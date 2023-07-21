package database

import (
	"testing"
)

func TestInsertIntoTable(t *testing.T) {
	dbClient, err := NewClient("username:password@/getground")

	if err != nil {
		t.Errorf("Error creating DB client: %s", err)
	}

	defer dbClient.Close()
}
