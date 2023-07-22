package entity

type Table struct {
	ID            int `db:"id"`
	Capacity      int `db:"capacity"`
	ReservedSeats int `db:"reserved_seats"`
}

type CreateTableResponseBody struct {
	ID       int `json:"id"`
	Capacity int `json:"capacity"`
}
