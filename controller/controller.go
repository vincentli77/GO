package controller

import (
	"Desktop/Go/model"
	"Desktop/Go/services"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func ReservationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCorsHeaders(w)
		if r.Method == http.MethodPost {
			var data model.ReservationData
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Error decoding JSON data: " + err.Error()))
				return
			}

			// Vérifie si une réservation existe déjà pour la même date et plage horaire
			var count int
			count = services.CheckReservation(db, data.Reservation_date, data.Start_time, data.End_time)
			if count > 0 {
				http.Error(w, "Reservation already exists for this time range", http.StatusConflict)
			} else {
				err = services.AddReservation(db, data)
				if err != nil {
					http.Error(w, "Error adding reservation1", http.StatusConflict)

				}
			}
			if err != nil {
				http.Error(w, "Error adding reservation2", http.StatusConflict)
			}

		}
	}
}

func HandleGetReservations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	setCorsHeaders(w)

	switch r.Method {
	case "GET":
		services.GetReservationsForWeek(db, w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
	}
}

func HandleGetAvailability(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	setCorsHeaders(w)

	switch r.Method {
	case "GET":
		data, err := services.GetAvailability(db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error getting availability: %s", err.Error())
			return
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error encoding availability: %s", err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
	}
}

func AvailabilityHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCorsHeaders(w)
		if r.Method == http.MethodPost {
			var data []map[string]string
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				http.Error(w, "Error decoding JSON data: "+err.Error(), http.StatusBadRequest)
				return
			}
			err = services.AddAvailability(db, data)
			if err != nil {
				http.Error(w, "Error adding availaibility", http.StatusConflict)

			}

		}
	}
}

func setCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}

func AdminHandleGetReservations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	setCorsHeaders(w)

	switch r.Method {
	case "GET":
		services.AdminGetReservationsForWeek(db, w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
	}
}

func GetUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	setCorsHeaders(w)

	switch r.Method {
	case "GET":
		userID := r.URL.Query().Get("id")
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
		users, err := services.GetUser(db, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error retrieving user: %v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
	}
}
