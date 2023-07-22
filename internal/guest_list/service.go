package guest_list

import (
	"fmt"
	"time"

	"github.com/getground/tech-tasks/backend/internal/entity"
	"github.com/getground/tech-tasks/backend/pkg/database"
)

type GuestListService interface {
	CreateTable(table *entity.Table) (*entity.CreateTableResponseBody, error)
	AddGuest(guest *entity.Guest) (*entity.AddGuestResponseBody, error)
	GetAllGuests() ([]entity.GetAllGuestsElement, error)
	GetAllCheckedInGuests() ([]entity.GetAllCheckedInGuestsElement, error)
	CheckInGuest(guest *entity.Guest) (*entity.CheckInGuestResponseBody, error)
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
	values := []interface{}{guest.Name, guest.AccompanyingGuests, guest.TableID}
	_, err = s.dbClient.Create("guest", columns, values...)
	if err != nil {
		return nil, err
	}

	// Update the number of reserved seats
	updatedReservedSeats := table.ReservedSeats + (guest.AccompanyingGuests + 1)
	columnsToUpdate := []string{"reserved_seats"}
	values = []interface{}{updatedReservedSeats}
	err = s.dbClient.Update("table", guest.TableID, columnsToUpdate, values...)
	if err != nil {
		return nil, err
	}

	newGuest := entity.AddGuestResponseBody{
		Name: guest.Name,
	}

	return &newGuest, nil
}

func (s *service) GetAllGuests() ([]entity.GetAllGuestsElement, error) {
	guests := []entity.GetAllGuestsElement{}

	err := s.dbClient.FindMany(&guests, "guest", nil, nil)
	if err != nil {
		return nil, err
	}

	return guests, nil
}

func (s *service) CheckInGuest(guest *entity.Guest) (*entity.CheckInGuestResponseBody, error) {
	// Retrieve the guest info from the DB
	var retrievedGuest entity.Guest
	err := s.dbClient.FindUnique(&retrievedGuest, "guest", "name", guest.Name)
	if err != nil {
		return nil, err
	}

	// Check if the guest is already checked in
	if retrievedGuest.TimeArrived != nil {
		err := fmt.Errorf("guest with name `%s` is already checked in", guest.Name)
		return nil, err
	}

	// Check in the guest if they have extras
	if guest.AccompanyingGuests > retrievedGuest.AccompanyingGuests {
		var table entity.Table
		err := s.dbClient.FindUnique(&table, "table", "id", retrievedGuest.TableID)
		if err != nil {
			return nil, err
		}

		extras := guest.AccompanyingGuests - retrievedGuest.AccompanyingGuests

		if (extras + table.ReservedSeats) > table.Capacity {
			err = fmt.Errorf("no available seats on table %d", retrievedGuest.TableID)
			return nil, err
		}

		columnsToUpdate := []string{"accompanying_guests"}
		values := []interface{}{guest.AccompanyingGuests}
		err = s.dbClient.Update("guest", retrievedGuest.ID, columnsToUpdate, values...)
		if err != nil {
			return nil, err
		}

		reservedSeats := table.ReservedSeats + extras
		columnsToUpdate = []string{"reserved_seats"}
		values = []interface{}{reservedSeats}
		err = s.dbClient.Update("table", table.ID, columnsToUpdate, values...)
		if err != nil {
			return nil, err
		}
	}

	// Check in hte guest
	timeArrived := time.Now().UTC().String()
	columnsToUpdate := []string{"time_arrived"}
	values := []interface{}{timeArrived}
	err = s.dbClient.Update("guest", retrievedGuest.ID, columnsToUpdate, values...)
	if err != nil {
		return nil, err
	}

	result := entity.CheckInGuestResponseBody{
		Name: guest.Name,
	}

	return &result, nil
}

func (s *service) GetAllCheckedInGuests() ([]entity.GetAllCheckedInGuestsElement, error) {
	guests := []entity.GetAllCheckedInGuestsElement{}
	condition := "time_arrived IS NOT NULL"

	err := s.dbClient.FindMany(&guests, "guest", &condition, nil)
	if err != nil {
		return nil, err
	}

	return guests, nil
}
