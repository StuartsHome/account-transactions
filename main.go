package main

import (
	"account-transactions/server"
	"account-transactions/store"
	"fmt"
	"log"
	"net/http"
)

var port = ":8080"

func main() {
	fmt.Printf("listening on port %s\n", port)

	db := store.New()
	r := server.NewRouter(db)

	log.Fatal(http.ListenAndServe(port, r))
}
