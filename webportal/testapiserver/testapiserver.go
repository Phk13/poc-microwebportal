package main

import (
	"log"

	"github.com/phk13/poc-micro/databaselayer"
	"github.com/phk13/poc-micro/webportal/api"
)

func main() {
	db, err := databaselayer.GetDatabaseHandler(databaselayer.MONGODB, "mongodb://172.17.0.8")
	if err != nil {
		log.Fatal(err)
	}
	api.RunApi("localhost:8080", db)
}
