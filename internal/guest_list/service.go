package guest_list

import (
	"github.com/getground/tech-tasks/backend/internal/entity"
	"github.com/getground/tech-tasks/backend/pkg/database"
)

type GuestListService interface {
	CreateTable(table *entity.Table) (*entity.Table, error)
}

type service struct {
	dbClient database.Client
}

func NewGuestListService(dbClient database.Client) GuestListService {
	return &service{dbClient}
}

func (s *service) CreateTable(table *entity.Table) (*entity.Table, error) {
	columns := []string{"capacity"}
	id, err := s.dbClient.Create("table", columns, table.Capacity)
	if err != nil {
		return nil, err
	}

	newTable := entity.Table{
		ID:       id,
		Capacity: table.Capacity,
	}

	return &newTable, nil
}
