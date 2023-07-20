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
	db, err := database.Connect()
	if err != nil {
		log.Printf("Fuck shit happened, %s", err)
	}
	defer db.Close()

	// Start API
	r := mux.NewRouter()
	guestListService := guest_list.NewGuestListService(db)
	guest_list.RegisterHandlers(r, guestListService)
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
