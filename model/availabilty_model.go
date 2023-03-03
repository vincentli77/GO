package model

import (
	"database/sql"
)

func GetAvailability(db *sql.DB) ([]map[string]string, error) {
	rows, err := db.Query("SELECT day, start_time, end_time FROM disponibilite")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]map[string]string, 0)

	for rows.Next() {
		var day string
		var start_time string
		var end_time string
		err := rows.Scan(&day, &start_time, &end_time)
		if err != nil {
			return nil, err
		}

		availibility := make(map[string]string)
		availibility["day"] = day
		availibility["start_time"] = start_time
		availibility["end_time"] = end_time

		data = append(data, availibility)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

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
	for _, data := range data {
		_, err = stmt.Exec(data["day"], data["start_time"], data["end_time"])
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
