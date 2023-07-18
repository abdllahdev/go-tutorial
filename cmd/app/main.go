package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// init mysql.
	db, err := sql.Open("mysql", "user:password@/getground")

	if err != nil {
		log.Fatal(err)
	}
	
	defer db.Close()

	// ping
	http.HandleFunc("/ping", handlerPing)
	http.ListenAndServe(":3000", nil)
}

func handlerPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong\n")
}
