package guest_list

import (
	"log"
	"testing"

	"github.com/getground/tech-tasks/backend/internal/entity"

	"github.com/getground/tech-tasks/backend/pkg/database"
)

func TestCreateTable(t *testing.T) {
	dbClient, err := database.NewClient("username:password@/getground")
	if err != nil {
		log.Printf("Error %s while creating DB client", err)
	}
	defer dbClient.Close()

	tableService := NewGuestListService(dbClient)

	var table entity.Table
	table.Capacity = 10

	newTable, err := tableService.CreateTable(&table)

	if err != nil || newTable == nil {
		t.Errorf("Error while creating a table, %v", err)
	}

	if newTable.Capacity != 10 {
		t.Errorf("Expected table capacity to be 10 but found %d", table.Capacity)
	}
}
