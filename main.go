package main

import (
	"log"
	"net/http"

	"github.com/chankei613/context-bundle-builder/internal/api"
	"github.com/chankei613/context-bundle-builder/internal/db"
)

func main() {
	conn, err := db.Init("context-bundle-builder.db")
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	router := api.NewRouter(conn)
	log.Println("context-bundle-builder backend listening on :8422")
	if err := http.ListenAndServe(":8422", router); err != nil {
		log.Fatal(err)
	}
}
