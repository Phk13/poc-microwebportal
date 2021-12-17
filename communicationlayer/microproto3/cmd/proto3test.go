package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"

	"github.com/phk13/poc-micro/communicationlayer/microproto3"
	"github.com/phk13/poc-micro/databaselayer"
	"google.golang.org/protobuf/proto"
)

func main() {
	op := flag.String("op", "s", "s for server, c for client")
	flag.Parse()
	switch strings.ToLower(*op) {
	case "s":
		RunProto3Server()
	case "c":
		log.Println("Launching client...")
		RunProto3Client()
	}
}

func RunProto3Client() {
	a := &microproto3.Animal{
		Id:         1,
		AnimalType: "Raptor",
		Nickname:   "rapto",
		Zone:       3,
		Age:        20,
	}
	log.Println("Marshaling:", a)
	data, err := proto.Marshal(a)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Sending...")
	SendData(data)
	log.Println("Data sent:", data)
}

func SendData(data []byte) {
	c, err := net.Dial("tcp", "127.0.0.1:8282")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	c.Write(data)
}

func RunProto3Server() {
	l, err := net.Listen("tcp", ":8282")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on port 8282...")
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		defer l.Close()
		go func(conn net.Conn) {
			defer conn.Close()
			data, err := ioutil.ReadAll(conn)
			if err != nil {
				return
			}
			a := &microproto3.Animal{}
			err = proto.Unmarshal(data, a)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("Received data...")
			fmt.Println(a)
		}(c)
	}
}

func SendDBToServer() {
	handler, err := databaselayer.GetDatabaseHandler(databaselayer.MONGODB, "mongodb://172.17.0.10")
	if err != nil {
		log.Fatal(err)
	}
	animals, err := handler.GetAvailableDinos()
	for _, animal := range animals {
		a := &microproto3.Animal{
			Id:         int32(animal.ID),
			AnimalType: animal.AnimalType,
			Nickname:   animal.Nickname,
			Zone:       int32(animal.Zone),
			Age:        int32(animal.Age),
		}
		data, err := proto.Marshal(a)
		if err != nil {
			log.Fatal(err)
		}
		SendData(data)
	}
}
