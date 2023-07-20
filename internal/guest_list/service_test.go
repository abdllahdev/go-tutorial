package guest_list

import (
	"testing"

	"github.com/getground/tech-tasks/backend/internal/entity"

	"github.com/getground/tech-tasks/backend/pkg/database"
)

func TestCreate(t *testing.T) {
	db, err := database.Connect()
	if err != nil {
		t.Errorf("Error while getting DB instance: %v", err)
	}
	defer db.Close()

	tableService := NewGuestListService(db)

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
