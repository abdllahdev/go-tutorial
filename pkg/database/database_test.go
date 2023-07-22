package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertIntoTable(t *testing.T) {
	dbClient, err := NewClient("username:password@/getground")
	assert.NotNil(t, err)
	defer dbClient.Close()
}
