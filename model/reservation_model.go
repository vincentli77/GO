package model

import (
	"Desktop/Go/schemas"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func AddReservation(db *sql.DB, data schemas.ReservationData) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var userID int64

	// Check if user exists, insert if it doesn't
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

	// Insert reservation
	_, err = tx.Exec("INSERT INTO reservations (user_id, reservation_date, start_time, end_time) VALUES (?, ?, ?, ?)", userID, data.Reservation_date, data.Start_time, data.End_time)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GetReservations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Définit la requête SQL pour récupérer les réservations d'une semaine spécifiée.
	query := "SELECT reservation_id, reservation_date, start_time, end_time, created_at, updated_at FROM reservations WHERE reservation_date >= ? AND reservation_date <= ?"

	// Définit une fonction différée qui récupère une éventuelle erreur panique.
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Error getting reservations: %v", r)
		}
	}()

	// Obtient la plage de dates de la semaine actuelle en utilisant une fonction.
	startDate, endDate := getWeekRange(r)

	// Prépare la requête SQL pour récupérer les réservations pour la plage de dates spécifiée.
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

	// Parcourt chaque ligne de résultats et ajoute les réservations à une liste de réservations.
	var reservations []schemas.Reservation
	for rows.Next() {
		var r schemas.Reservation
		//Extrait les  valeurs de de la query et les stock dans la variable à l'aide de ses pointeurs.
		if err := rows.Scan(&r.Id, &r.Reservation_date, &r.Start_time, &r.End_time, &r.Created_at, &r.Updated_at); err != nil {
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

	// Si aucune réservation n'a été trouvée, renvoie une réponse avec code 404 et un message d'erreur.
	if len(reservations) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "No reservations found for the specified date range")
		return
	}

	// Encode les réservations en JSON et les renvoie dans la réponse HTTP.
	json.NewEncoder(w).Encode(reservations)
}

func CheckReservation(db *sql.DB, reservation_date string, start_time string, end_time string) int {
	query := "SELECT COUNT(*) FROM reservations WHERE reservation_date = ? AND start_time >= ? AND end_time <= ?"
	var count int
	err := db.QueryRow(query, reservation_date, start_time, end_time).Scan(&count)
	fmt.Println(err)

	return count
}

func AdminGetReservations(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	query := "SELECT reservation_id,user_id,reservation_date,start_time,end_time,created_at,updated_at FROM reservations WHERE reservation_date >= ? AND reservation_date <= ?"

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

	var reservations []schemas.ReservationAdmin
	for rows.Next() {
		var row schemas.ReservationAdmin
		if err := rows.Scan(&row.Id, &row.User_id, &row.Reservation_date, &row.Start_time, &row.End_time, &row.Created_at, &row.Updated_at); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error scanning row: %v", err)
			return
		}
		reservations = append(reservations, row)
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

func DeleteReservation(db *sql.DB, userID int) error {
	query := "DELETE FROM reservations WHERE reservation_id = ?"
	_, err := db.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

func getWeekRange(r *http.Request) (startDate string, endDate string) {

	startDate = r.URL.Query().Get("start_date")
	endDate = r.URL.Query().Get("end_date")

	return startDate, endDate
}
