package main

import (
	"log"
	"net/http"

	"github.com/getground/tech-tasks/backend/internal/guest_list"
	"github.com/getground/tech-tasks/backend/pkg/database"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// Initiate DB
	dbClient, err := database.NewClient("username:password@tcp(mysql:3306)/getground")
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Close()

	// Start API
	r := mux.NewRouter()
	guestListService := guest_list.NewGuestListService(dbClient)
	guest_list.RegisterHandlers(r, guestListService)
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
