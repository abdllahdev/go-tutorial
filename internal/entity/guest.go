package entity

type Guest struct {
	ID                 int     `json:"id"                  db:"id"`
	Name               string  `json:"name"                db:"name"`
	AccompanyingGuests int     `json:"accompanying_guests" db:"accompanying_guests"`
	TableID            int     `json:"table_id"            db:"table_id"`
	TimeArrived        *string `json:"time_arrived"        db:"time_arrived"`
}

type AddGuestRequestBody struct {
	Table              int `json:"table"`
	AccompanyingGuests int `json:"accompanying_guests"`
}

type AddGuestResponseBody struct {
	Name string `json:"name"`
}

type CheckInGuestRequestBody struct {
	AccompanyingGuests int `json:"accompanying_guests"`
}

type CheckInGuestResponseBody struct {
	Name string `json:"name"`
}

type GetAllGuestsElement struct {
	Name               string `json:"name"                db:"name"`
	AccompanyingGuests int    `json:"accompanying_guests" db:"accompanying_guests"`
	TableID            int    `json:"table_id"            db:"table_id"`
}

type GetAllGuestsResponseBody struct {
	Guests []GetAllGuestsElement `json:"guests"`
}

type GetAllCheckedInGuestsElement struct {
	Name               string `json:"name"                db:"name"`
	AccompanyingGuests int    `json:"accompanying_guests" db:"accompanying_guests"`
	TimeArrived        string `json:"time_arrived"        db:"time_arrived"`
}

type GetAllCheckedInGuestsResponseBody struct {
	Guests []GetAllCheckedInGuestsElement `json:"guests"`
}
