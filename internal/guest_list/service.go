package guest_list

import (
	"fmt"

	"github.com/getground/tech-tasks/backend/internal/entity"
	"github.com/getground/tech-tasks/backend/pkg/database"
)

type GuestListService interface {
	CreateTable(table *entity.Table) (*entity.CreateTableResponseBody, error)
	AddGuest(guest *entity.Guest) (*entity.AddGuestResponseBody, error)
	GetAllGuests() ([]entity.GuestData, error)
}

type service struct {
	dbClient database.Client
}

func NewGuestListService(dbClient database.Client) GuestListService {
	return &service{dbClient}
}

func (s *service) CreateTable(table *entity.Table) (*entity.CreateTableResponseBody, error) {
	columns := []string{"capacity"}
	id, err := s.dbClient.Create("table", columns, table.Capacity)
	if err != nil {
		return nil, err
	}

	newTable := entity.CreateTableResponseBody{
		ID:       id,
		Capacity: table.Capacity,
	}

	return &newTable, nil
}

func (s *service) AddGuest(guest *entity.Guest) (*entity.AddGuestResponseBody, error) {
	// Check if a guest with the same already exists in the DB
	guestExists, err := s.dbClient.Exists("guest", "name", guest.Name)
	if err != nil {
		return nil, err
	}
	if guestExists {
		err = fmt.Errorf("guest with name %s already exists", guest.Name)
		return nil, err
	}

	// Check if table already exists in the DB
	var table entity.Table
	err = s.dbClient.FindUnique(&table, "table", "id", guest.TableID)
	if err != nil {
		return nil, err
	}

	// Check if there are enough seats
	if table.ReservedSeats+(guest.AccompanyingGuests+1) > table.Capacity {
		err = fmt.Errorf("no available seats on table %d", guest.TableID)
		return nil, err
	}

	// Add a new guest
	columns := []string{"name", "accompanying_guests", "table_id"}
	_, err = s.dbClient.Create("guest", columns, guest.Name, guest.AccompanyingGuests, guest.TableID)
	if err != nil {
		return nil, err
	}

	// Update the number of reserved seats
	columnsToUpdate := []string{"reserved_seats"}
	updatedReservedSeats := table.ReservedSeats + (guest.AccompanyingGuests + 1)
	s.dbClient.Update("table", guest.TableID, columnsToUpdate, updatedReservedSeats)

	newGuest := entity.AddGuestResponseBody{
		Name: guest.Name,
	}

	return &newGuest, nil
}

func (s *service) GetAllGuests() ([]entity.GuestData, error) {
	guests := []entity.GuestData{}

	err := s.dbClient.FindMany(&guests, "guest", nil, nil)
	if err != nil {
		return nil, err
	}

	return guests, nil
}
