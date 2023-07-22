package guest_list

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/getground/tech-tasks/backend/internal/entity"

	"github.com/getground/tech-tasks/backend/pkg/database"
)

const (
	dsn = "username:password@/getground"
)

func setup() database.Client {
	dbClient, err := database.NewClient(dsn)
	if err != nil {
		log.Fatalf("Error while connecting to the DB, %v", err)
	}
	return dbClient
}

func cleanupTable(dbClient database.Client, tableName string) {
	err := dbClient.DeleteAll(tableName)
	if err != nil {
		log.Fatalf("Error while cleaning table %s, %v", tableName, err)
	}
}

func TestCreateTable(t *testing.T) {
	// Setup database
	dbClient := setup()
	defer dbClient.Close()

	// Cleanup tables
	cleanupTable(dbClient, "table")

	// Create guest list service
	guestListService := NewGuestListService(dbClient)

	// Test creating a new table
	var table entity.Table
	table.Capacity = 10
	newTable, err := guestListService.CreateTable(&table)
	assert.Nil(t, err, "Error while creating a new table, %v", err)
	assert.NotNil(t, newTable, "Expected table to have value but found nil")
}

func TestAddGuest(t *testing.T) {
	// Setup database
	dbClient := setup()
	defer dbClient.Close()

	// Cleanup tables
	cleanupTable(dbClient, "table")
	cleanupTable(dbClient, "guest")

	// Create guest list service
	guestListService := NewGuestListService(dbClient)

	// Create a new table
	var table entity.Table
	table.Capacity = 5
	newTable, err := guestListService.CreateTable(&table)
	assert.Nil(t, err, "Error while creating a new table, %v", err)
	assert.NotNil(t, newTable, "Expected table to have value but found nil")

	// Test adding a new guest with a table id that exists
	var guest entity.Guest
	guest.Name = "john"
	guest.AccompanyingGuests = 3
	guest.TableID = newTable.ID
	newGuest, err := guestListService.AddGuest(&guest)
	assert.Nil(t, err, "Error while creating a new guest, %v", err)
	assert.NotNil(t, newGuest, "Expected guest to have value but found nil")

	// Test adding a new guest with the same name
	_, err = guestListService.AddGuest(&guest)
	expectedErrorMsg := fmt.Sprintf("guest with name %s already exists", guest.Name)
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be %v but found %v", err, expectedErrorMsg)

	// Test adding a new guest in a table with no available seats
	guest.Name = "rob"
	guest.AccompanyingGuests = 1
	_, err = guestListService.AddGuest(&guest)
	expectedErrorMsg = fmt.Sprintf("no available seats on table %d", guest.TableID)
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be %v but found %v", err, expectedErrorMsg)

	// Test adding a new guest with a table id that does not exist
	guest.TableID = newTable.ID + 1
	_, err = guestListService.AddGuest(&guest)
	expectedErrorMsg = "sql: no rows in result set"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be %v but found %v", err, expectedErrorMsg)
}
