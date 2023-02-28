package main

import (
	"Desktop/Go/controller"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// Établissez la connexion à la base de données.
	// db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/reservation")
	db, err := sql.Open("mysql", "b50a7750fdd7c4:5bedff1e@tcp(eu-cdbr-west-03.cleardb.net:3306)/heroku_0429c505d3dfa57")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Vérifiez si la connexion est opérationnelle.
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/reservations", controller.ReservationHandler(db))
	http.HandleFunc("/get_reservations", func(w http.ResponseWriter, r *http.Request) {
		controller.HandleGetReservations(db, w, r)
	})
	http.HandleFunc("/admin_get_reservations", func(w http.ResponseWriter, r *http.Request) {
		controller.AdminHandleGetReservations(db, w, r)
	})
	http.HandleFunc("/add_availability", controller.AvailabilityHandler(db))
	http.HandleFunc("/get_availability", func(w http.ResponseWriter, r *http.Request) {
		controller.HandleGetAvailability(db, w, r)
	})
	http.HandleFunc("/get_user", func(w http.ResponseWriter, r *http.Request) {
		controller.GetUser(db, w, r)
	})
	fmt.Println("Serveur web démarré sur le port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
