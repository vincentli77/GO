package controller

import (
	"Desktop/Go/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

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
