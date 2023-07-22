package entity

type Guest struct {
	ID                 int    `db:"id"`
	Name               string `db:"name"`
	AccompanyingGuests int    `db:"accompanying_guests"`
	TableID            int    `db:"table_id"`
	TimeArrived        string `db:"time_arrived"`
}

type AddGuestRequestBody struct {
	Table              int `json:"table"`
	AccompanyingGuests int `json:"accompanying_guests"`
}

type AddGuestResponseBody struct {
	Name string `json:"name"`
}
