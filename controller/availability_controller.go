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
	setCorsHeaders(w, r)

	switch r.Method {
	case "GET":
		data, err := model.GetAvailability(db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Erreur lors de la récupération de la disponibilité : %s", err.Error())
			return
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Erreur lors de l'encodage en JSON de la disponibilité : %s", err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	default:
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
