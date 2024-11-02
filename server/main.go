package main

import (
	"log"
	"net/http"

	"github.com/ArpitKRS/go-med-record/config"
	"github.com/ArpitKRS/go-med-record/routes"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Client().Disconnect(config.Ctx)

	router := routes.SetupRouter()
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
