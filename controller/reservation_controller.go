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

func ReservationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Définir les en-têtes CORS pour permettre les requêtes cross-domain
		setCorsHeaders(w, r)

		// Vérifier si la requête est de type POST
		if r.Method == http.MethodPost {

			// Déclarer une variable pour stocker les données de réservation
			var data schemas.ReservationData

			// Décoder le corps de la requête JSON dans la variable de réservation
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				// Renvoyer une erreur HTTP si la lecture des données de la requête JSON a échoué
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Error decoding JSON data: " + err.Error()))
				return
			}

			// Vérifier si une réservation existe déjà pour la même date et plage horaire
			var count int
			count = model.CheckReservation(db, data.Reservation_date, data.Start_time, data.End_time)
			if count > 0 {
				// Renvoyer une erreur HTTP si une réservation existe déjà pour cette plage horaire
				http.Error(w, "Reservation already exists for this time range", http.StatusConflict)
			} else {
				// Ajouter la nouvelle réservation à la base de données
				err = model.AddReservation(db, data)
				if err != nil {
					// Renvoyer une erreur HTTP si l'ajout de la réservation a échoué
					http.Error(w, "Error adding reservation", http.StatusConflict)
				}
			}
			if err != nil {
				// Renvoyer une erreur HTTP si l'ajout de la réservation a échoué
				http.Error(w, "Error adding reservation", http.StatusConflict)
			}
		}
	}
}

func GetReservationsHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Autorise les requêtes Cross-Origin Resource Sharing (CORS)
	setCorsHeaders(w, r)

	// Détermine la méthode HTTP utilisée
	switch r.Method {
	case "GET":
		// Si la méthode est GET, appelle la fonction GetReservationsForWeek pour récupérer les réservations
		model.GetReservationsForWeek(db, w, r)
	default:
		// Si la méthode n'est pas GET, renvoie une réponse d'erreur avec le code HTTP 405 (Method Not Allowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Méthode %s non autorisée", r.Method)
	}
}

func AdminGetReservationsHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	setCorsHeaders(w, r)

	switch r.Method {
	case "GET":
		model.AdminGetReservationsForWeek(db, w, r)
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
