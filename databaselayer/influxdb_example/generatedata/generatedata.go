package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

var animaltags = []string{"Tyrannosaurus rex;Rex", "Velociraptor;Rapto", "Velociraptor;Velo", "Carnotaurus;Carno"}

const myDB = "micro"

func main() {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://172.17.0.11:8086",
		Username: "",
		Password: "",
	})
	if err != nil {
		log.Fatal(err)
	}
	queryDB(c, "", "Create DATABASE "+myDB)
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  myDB,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}
	wg := sync.WaitGroup{}
	detectSignal := checkStopOsSignals(&wg)
	rand.Seed(time.Now().UnixNano())
	for !(*detectSignal) {
		animaltag := animaltags[rand.Intn(len(animaltags))]
		split := strings.Split(animaltag, ";")
		tags := map[string]string{
			"animal_type": split[0],
			"nickname":    split[1],
		}
		fields := map[string]interface{}{
			"weight": rand.Intn(300) + 1,
		}
		fmt.Println(animaltag, fields["weight"])
		pt, err := client.NewPoint("weightmeasures", tags, fields, time.Now())
		if err != nil {
			log.Println(err)
			continue
		}
		bp.AddPoint(pt)
		time.Sleep(1 * time.Second)
	}
	log.Println("Exit signal triggered, writing data...")
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}
	wg.Wait()
	log.Println("Exiting program...")

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

func checkStopOsSignals(wg *sync.WaitGroup) *bool {
	Signal := false
	wg.Add(1)
	go func(s *bool) {
		ch := make(chan os.Signal, 1000)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		log.Println("Exit signals received")
		*s = true
		wg.Done()
	}(&Signal)
	return &Signal
}
