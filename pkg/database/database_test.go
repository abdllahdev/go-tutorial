package database

import (
	"testing"
)

func TestConnect(t *testing.T) {
	// Test when database is not initialized
	db, err := Connect()

	if err != nil {
		t.Errorf("Error while getting DB instance: %v", err)
	}

	if db == nil {
		t.Errorf("DB instance is nil")
	}
}
