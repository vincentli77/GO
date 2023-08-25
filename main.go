package main

import (
	"Desktop/Go/controller"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// Établissement de la connexion à la base de données.
	db, err := sql.Open("mysql", "b50a7750fdd7c4:5bedff1e@tcp(eu-cdbr-west-03.cleardb.net:3306)/heroku_0429c505d3dfa57")

	// Si une erreur survient lors de l'établissement de la connexion, arrêter le programme et afficher l'erreur.
	if err != nil {
		panic(err)
	}
	// Fermer la connexion à la base de données une fois la fonction main terminée.
	defer db.Close()

	// Vérification que la connexion à la base de données est opérationnelle.
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/reservations", controller.AddReservationHandler(db))
	http.HandleFunc("/get_reservations", func(w http.ResponseWriter, r *http.Request) {
		controller.GetReservationsHandler(db, w, r)
	})

	http.HandleFunc("/admin_get_reservations", func(w http.ResponseWriter, r *http.Request) {
		controller.AdminGetReservationsHandler(db, w, r)
	})
	http.HandleFunc("/delete_reservation", func(w http.ResponseWriter, r *http.Request) {
		controller.DeleteReservationHandler(db, w, r)
	})

	http.HandleFunc("/add_availability", controller.AddAvailabilityHandler(db))
	http.HandleFunc("/get_availability", func(w http.ResponseWriter, r *http.Request) {
		controller.GetAvailabilityHandler(db, w, r)
	})
	http.HandleFunc("/get_user", func(w http.ResponseWriter, r *http.Request) {
		controller.GetUserHandler(db, w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Serveur web démarré sur le port " + port + "...")
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}

}
