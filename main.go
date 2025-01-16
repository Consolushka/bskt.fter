package main

import (
	"IMP/app/api"
	"IMP/app/database"
	"log"
	"net/http"
)

func main() {
	database.OpenDbConnection()
	mux := api.Serve()

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
