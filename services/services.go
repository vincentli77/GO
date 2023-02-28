package services

import (
	"Desktop/Go/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func AddReservation(db *sql.DB, data model.ReservationData) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// insert user if it doesn't exist
	var userID int64
	err = tx.QueryRow("SELECT id FROM users WHERE first_name = ? AND last_name = ?", data.First_name, data.Last_name).Scan(&userID)
	if err == sql.ErrNoRows {
		res, err := tx.Exec("INSERT INTO users (first_name, last_name) VALUES (?, ?)", data.First_name, data.Last_name)
		if err != nil {
			tx.Rollback()
			return err
		}
		userID, err = res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return err
		}
	} else if err != nil {
		tx.Rollback()
		return err
	}

	// insert reservation
	_, err = tx.Exec("INSERT INTO reservations (userID, reservation_date, start_time, end_time) VALUES (?, ?, ?, ?)", userID, data.Reservation_date, data.Start_time, data.End_time)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GetReservationsForWeek(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	query := "SELECT * FROM reservations WHERE reservation_date >= ? AND reservation_date <= ?"

	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Error getting reservations: %v", r)
		}
	}()

	startDate, endDate := getWeekRange(r)
	stmt, err := db.Prepare(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error preparing statement: %v", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error executing query: %v", err)
		return
	}
	defer rows.Close()

	var reservations []model.Reservation
	for rows.Next() {
		var r model.Reservation
		if err := rows.Scan(&r.Id, &r.User_id, &r.Reservation_date, &r.Start_time, &r.End_time, &r.Status, &r.Created_at, &r.Updated_at); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error scanning row: %v", err)
			return
		}
		reservations = append(reservations, r)
	}

	if err := rows.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error after scanning rows: %v", err)
		return
	}

	if len(reservations) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "No reservations found for the specified date range")
		return
	}

	json.NewEncoder(w).Encode(reservations)
}

func getWeekRange(r *http.Request) (startDate string, endDate string) {

	startDate = r.URL.Query().Get("start_date")
	endDate = r.URL.Query().Get("end_date")

	return startDate, endDate
}

func CheckReservation(db *sql.DB, reservation_date string, start_time string, end_time string) int {
	query := "SELECT COUNT(*) FROM reservations WHERE reservation_date = ? AND start_time >= ? AND end_time <= ?"
	var count int
	err := db.QueryRow(query, reservation_date, start_time, end_time).Scan(&count)
	fmt.Println(err)

	return count
}

func GetAvailability(db *sql.DB) ([]map[string]string, error) {
	rows, err := db.Query("SELECT day, start_time, end_time FROM disponibilite")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]map[string]string, 0)
	for rows.Next() {
		var day string
		var start_time string
		var end_time string
		err := rows.Scan(&day, &start_time, &end_time)
		if err != nil {
			return nil, err
		}

		d := make(map[string]string)
		d["day"] = day
		d["start_time"] = start_time
		d["end_time"] = end_time

		data = append(data, d)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func AddAvailability(db *sql.DB, data []map[string]string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM disponibilite")
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO disponibilite (day, start_time, end_time) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	for _, d := range data {
		_, err = stmt.Exec(d["day"], d["start_time"], d["end_time"])
		if err != nil {
			return err
		}
	}
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
