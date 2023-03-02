package model

import (
	"database/sql"
)

func GetAvailability(db *sql.DB) ([]map[string]string, error) {
	// Récupère toutes les disponibilités enregistrées dans la table 'disponibilite'
	rows, err := db.Query("SELECT day, start_time, end_time FROM disponibilite")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Crée une slice de maps pour stocker les données récupérées
	data := make([]map[string]string, 0)

	// Parcourt chaque ligne renvoyée par la requête SQL
	for rows.Next() {
		var day string
		var start_time string
		var end_time string
		err := rows.Scan(&day, &start_time, &end_time)
		if err != nil {
			return nil, err
		}

		// Crée une map pour stocker les données de chaque ligne
		d := make(map[string]string)
		d["day"] = day
		d["start_time"] = start_time
		d["end_time"] = end_time

		// Ajoute la map créée à la slice de données
		data = append(data, d)
	}
	// Vérifie s'il y a eu une erreur lors du parcours des lignes renvoyées
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// Retourne les données récupérées sous forme de slice de maps
	return data, nil
}

func AddAvailability(db *sql.DB, data []map[string]string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM disponibilite")
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO disponibilite (day, start_time, end_time) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	for _, d := range data {
		_, err = stmt.Exec(d["day"], d["start_time"], d["end_time"])
		if err != nil {
			return err
		}
	}
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
