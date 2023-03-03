package controller

import (
	"Desktop/Go/model"
	"Desktop/Go/schemas"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func AddReservationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		setCorsHeaders(w, r)

		if r.Method == http.MethodPost {

			var data schemas.ReservationData

			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Error decoding JSON data: " + err.Error()))
				return
			}

			// Vérifier si une réservation existe déjà pour la même date et plage horaire
			var count int
			count = model.CheckReservation(db, data.Reservation_date, data.Start_time, data.End_time)
			if count > 0 {
				http.Error(w, "Reservation already exists for this time range", http.StatusConflict)
			} else {
				err = model.AddReservation(db, data)
				if err != nil {
					http.Error(w, "Error adding reservation", http.StatusConflict)
				}
			}
			if err != nil {
				http.Error(w, "Error adding reservation", http.StatusConflict)
			}
		}
	}
}

func GetReservationsHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	setCorsHeaders(w, r)

	switch r.Method {
	case "GET":
		model.GetReservations(db, w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Méthode %s non autorisée", r.Method)
	}
}

func AdminGetReservationsHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	setCorsHeaders(w, r)

	switch r.Method {
	case "GET":
		model.AdminGetReservations(db, w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
	}
}

func DeleteReservationHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	setCorsHeaders(w, r)

	switch r.Method {
	case "DELETE":
		userID := r.URL.Query().Get("reservation_id")
		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Reservation ID missing")
			return
		}
		id, err := strconv.Atoi(userID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid reservation ID")
			return
		}
		err = model.DeleteReservation(db, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error deleting user: %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "User deleted successfully")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
	}
}

func setCorsHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}
