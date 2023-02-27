package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ReservationData struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	ReservationDate string `json:"reservation_date"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
}

type Reservation struct {
	ID              int
	UserID          int
	ReservationDate string
	StartTime       string
	EndTime         string
	Status          string
	CreatedAt       string
	UpdatedAt       string
}

type Availability struct {
	Day       string `json:"day"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func AddReservation(db *sql.DB, data ReservationData) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// insert user if it doesn't exist
	var userID int64
	err = tx.QueryRow("SELECT id FROM users WHERE first_name = ? AND last_name = ?", data.FirstName, data.LastName).Scan(&userID)
	if err == sql.ErrNoRows {
		res, err := tx.Exec("INSERT INTO users (first_name, last_name) VALUES (?, ?)", data.FirstName, data.LastName)
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
	_, err = tx.Exec("INSERT INTO reservations (user_id, reservation_date, start_time, end_time) VALUES (?, ?, ?, ?)", userID, data.ReservationDate, data.StartTime, data.EndTime)
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

	var reservations []Reservation
	for rows.Next() {
		var r Reservation
		if err := rows.Scan(&r.ID, &r.UserID, &r.ReservationDate, &r.StartTime, &r.EndTime, &r.Status, &r.CreatedAt, &r.UpdatedAt); err != nil {
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

func CheckReservation(db *sql.DB, reservationDate string, startTime string, endTime string) int {
	query := "SELECT COUNT(*) FROM reservations WHERE reservation_date = ? AND start_time >= ? AND end_time <= ?"
	var count int
	err := db.QueryRow(query, reservationDate, startTime, endTime).Scan(&count)
	fmt.Println(count)
	fmt.Println(err)

	return count
}

func AddAvailability(db *sql.DB, data Availability) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	var data2 = []map[string]string{
		{"day": "monday", "start_time": "09:00", "end_time": "09:30"},
		{"day": "monday", "start_time": "09:30", "end_time": "10:00"},
		{"day": "monday", "start_time": "10:00", "end_time": "10:30"},
	}

	stmt, err := tx.Prepare("INSERT INTO disponibilite (day, start_time, end_time) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	for _, d := range data2 {
		_, err = stmt.Exec(d["day"], d["start_time"], d["end_time"])
		if err != nil {
			return err
		}
	}
	// insert reservation
	// _, err = tx.Exec("INSERT INTO disponibilite (day, start_time, end_time) VALUES (?, ?, ?)", data.Day, data.StartTime, data.EndTime)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
