package guest_list

import (
	"log"
	"net/http"
	"testing"

	"github.com/getground/tech-tasks/backend/internal/entity"
	"github.com/getground/tech-tasks/backend/internal/test"
	"github.com/getground/tech-tasks/backend/pkg/database"
	"github.com/gorilla/mux"
)

func TestAPI(t *testing.T) {
	dbClient, err := database.NewClient(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Close()

	// Cleanup tables
	cleanupTable(dbClient, "table")
	cleanupTable(dbClient, "guest")

	// Register routes
	r := mux.NewRouter()
	guestListService := NewGuestListService(dbClient)
	RegisterHandlers(r, guestListService)
	http.Handle("/", r)

	// Create new table
	var table entity.Table
	table.Capacity = 5
	tableResponse, err := guestListService.CreateTable(&table)
	if err != nil {
		log.Fatal(err)
	}
	table.ID = tableResponse.ID

	tests := []test.APITestCase{
		{
			Name:   "Create table",
			Method: "POST",
			URL:    "/tables",
			Body: entity.CreateTableRequestBody{
				Capacity: 5,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedResponse: map[string]interface{}{
				"capacity": 5,
			},
		},
		{
			Name:   "Add a new guest",
			Method: "POST",
			URL:    "/guest_list/john",
			Body: entity.AddGuestRequestBody{
				Table:              table.ID,
				AccompanyingGuests: 0,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedResponse: map[string]interface{}{
				"name": "john",
			},
		},
		{
			Name:           "Get all guests",
			Method:         "GET",
			URL:            "/guest_list",
			Body:           nil,
			ExpectedStatus: http.StatusOK,
			ExpectedResponse: map[string]interface{}{
				"guests": []interface{}{
					map[string]interface{}{
						"accompanying_guests": 0,
						"name":                "john",
						"table_id":            table.ID,
					},
				},
			},
		},
		{
			Name:   "Check in guest",
			Method: "PUT",
			URL:    "/guests/john",
			Body: map[string]interface{}{
				"accompanying_guests": 3,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedResponse: map[string]interface{}{
				"name": "john",
			},
		},
		{
			Name:           "Get all checked in guests",
			Method:         "GET",
			URL:            "/guests",
			Body:           nil,
			ExpectedStatus: http.StatusOK,
			ExpectedResponse: map[string]interface{}{
				"guests": []interface{}{
					map[string]interface{}{
						"name": "john",
					},
				},
			},
		},
		{
			Name:           "Count empty seats",
			Method:         "GET",
			URL:            "/seats_empty",
			Body:           nil,
			ExpectedStatus: http.StatusOK,
			ExpectedResponse: map[string]interface{}{
				"seats_empty": 6,
			},
		},
		{
			Name:             "Check out guest",
			Method:           "DELETE",
			URL:              "/guests/john",
			Body:             nil,
			ExpectedStatus:   http.StatusNoContent,
			ExpectedResponse: nil,
		},
	}

	for _, tc := range tests {
		test.Endpoint(t, r, tc)
	}
}
