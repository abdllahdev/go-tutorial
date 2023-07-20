package guest_list

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/getground/tech-tasks/backend/internal/entity"
)

type GuestListService interface {
	CreateTable(table *entity.Table) (*entity.Table, error)
}

type service struct {
	db *sql.DB
}

func NewGuestListService(db *sql.DB) GuestListService {
	return &service{db}
}

func (s *service) CreateTable(table *entity.Table) (*entity.Table, error) {
	query := "INSERT INTO `table` (`capacity`) VALUES (?)"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, table.Capacity)
	if err != nil {
		log.Printf("Error %s when inserting row into table", err)
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error %s while getting created row ID", err)
		return nil, err
	}
	log.Printf("Created a new table")

	newTable := entity.Table{
		ID:       int(id),
		Capacity: table.Capacity,
	}

	return &newTable, nil
}
