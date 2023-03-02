package controller

import (
	"Desktop/Go/entities"
	"Desktop/Go/model"
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
			var data entities.ReservationData

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

func GetAvailabilityHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Autorise les requêtes Cross-Origin Resource Sharing (CORS)
	setCorsHeaders(w, r)

	// Détermine la méthode HTTP utilisée
	switch r.Method {
	case "GET":
		// Appelle la fonction GetAvailability pour récupérer les données de disponibilité
		data, err := model.GetAvailability(db)
		if err != nil {
			// Si une erreur se produit, renvoie une réponse d'erreur avec le code HTTP 500 (Internal Server Error)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Erreur lors de la récupération de la disponibilité : %s", err.Error())
			return
		}
		// Encode les données en JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			// Si une erreur se produit lors de l'encodage en JSON, renvoie une réponse d'erreur avec le code HTTP 500 (Internal Server Error)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Erreur lors de l'encodage en JSON de la disponibilité : %s", err.Error())
			return
		}
		// Configure le type de contenu de la réponse
		w.Header().Set("Content-Type", "application/json")
		// Écrit la réponse JSON
		w.Write(jsonData)
	default:
		// Si la méthode n'est pas GET, renvoie une réponse d'erreur avec le code HTTP 405 (Method Not Allowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Méthode %s non autorisée", r.Method)
	}
}

func AvailabilityHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCorsHeaders(w, r)
		if r.Method == http.MethodPost {
			var data []map[string]string
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				http.Error(w, "Error decoding JSON data: "+err.Error(), http.StatusBadRequest)
				return
			}
			err = model.AddAvailability(db, data)
			if err != nil {
				http.Error(w, "Error adding availaibility", http.StatusConflict)

			}

		}
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

func GetUserHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	setCorsHeaders(w, r)

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
		users, err := model.GetUser(db, id)
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
