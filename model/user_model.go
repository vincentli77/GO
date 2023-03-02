package model

import (
	"Desktop/Go/schemas"
	"database/sql"
)

func GetUser(db *sql.DB, userID int) (*schemas.User, error) {
	query := "SELECT id, first_name,last_name,created_at,updated_at FROM users WHERE id = ?"
	var user schemas.User
	err := db.QueryRow(query, userID).Scan(&user.Id, &user.First_name, &user.Last_name, &user.Created_at, &user.Updated_at)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
