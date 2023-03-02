package controller

import (
	"Desktop/Go/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

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
