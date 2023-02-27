package main

import (
	"Desktop/Go/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// Établissez la connexion à la base de données.
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/reservation")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Vérifiez si la connexion est opérationnelle.
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/reservations", reservationHandler(db))
	http.HandleFunc("/get_reservations", func(w http.ResponseWriter, r *http.Request) {
		handleGetReservations(db, w, r)
	})
	http.HandleFunc("/addAvailability", availabilityHandler(db))

	fmt.Println("Serveur web démarré sur le port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}

func reservationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			count = model.CheckReservation(db, data.ReservationDate, data.StartTime, data.EndTime)
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

func handleGetReservations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		model.GetReservationsForWeek(db, w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
	}
}

func availabilityHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var data model.Availability
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Error decoding JSON data: " + err.Error()))
				return
			}
			err = model.AddAvailability(db, data)
			fmt.Println(data)

			if err != nil {
				http.Error(w, "Error adding availaibility", http.StatusConflict)

			}

		}
	}
}
