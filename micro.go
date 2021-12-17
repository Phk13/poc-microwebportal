package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/phk13/poc-micro/webportal"
)

type configuration struct {
	ServerAddress      string `json:"webserver"`
	DatabaseType       uint8  `json:"databasetype"`
	DatabaseConnection string `json:"dbconnection"`
	FrontEnd           string `json:"frontend"`
	Secret             string `json:"secret"`
}

func main() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	config := new(configuration)
	json.NewDecoder(file).Decode(config)
	log.Println("Starting web server on address", config.ServerAddress)
	webportal.RunWebPortal(config.DatabaseType, config.ServerAddress, config.DatabaseConnection, config.FrontEnd, config.Secret)
}
