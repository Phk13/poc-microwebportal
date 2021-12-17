package main

import (
	"log"

	"github.com/influxdata/influxdb/client/v2"
)

func main() {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://172.17.0.11:8086",
		Username: "",
		Password: "",
	})
	if err != nil {
		log.Fatal(err)
	}
	res, err := queryDB(c, "micro", "select * from weightmeasures where animal_type = 'Velociraptor'")
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range res {
		log.Println("messages:", v.Messages)
		for _, s := range v.Series {
			log.Println("series name:", s.Name)
			log.Println("series columns:", s.Columns)
			for _, r := range s.Values {
				log.Println("row values:", r)
			}

		}
	}
}

func queryDB(c client.Client, database, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: database,
	}
	response, err := c.Query(q)
	if err != nil {
		return res, err
	}
	if response.Error() != nil {
		return res, response.Error()
	}
	res = response.Results

	return res, nil
}
