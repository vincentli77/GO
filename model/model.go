package model

type ReservationData struct {
	First_name       string `json:"first_name"`
	Last_name        string `json:"last_name"`
	Reservation_date string `json:"reservation_date"`
	Start_time       string `json:"start_time"`
	End_time         string `json:"end_time"`
}

type Reservation struct {
	Id               int    `json:"id"`
	User_id          int    `json:"user_id"`
	Reservation_date string `json:"reservation_date"`
	Start_time       string `json:"start_time"`
	End_time         string `json:"end_time"`
	Status           string `json:"status"`
	Created_at       string `json:"created_at"`
	Updated_at       string `json:"updated_at"`
}

type Availability struct {
	Day        string `json:"day"`
	Start_time string `json:"start_time"`
	End_time   string `json:"end_time"`
}
