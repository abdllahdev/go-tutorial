package entity

type Table struct {
	ID            int `json:"id" db:"id"`
	Capacity      int `json:"capacity" db:"capacity"`
	ReservedSeats int `json:"reserved_seats" db:"reserved_seats"`
}

type CreateTableResponseBody struct {
	ID       int `json:"id"`
	Capacity int `json:"capacity"`
}
