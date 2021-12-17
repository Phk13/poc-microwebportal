package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/phk13/poc-micro/communicationlayer/microproto2"
)

func main() {
	op := flag.String("op", "s", "s for server, c for client")
	flag.Parse()
	switch strings.ToLower(*op) {
	case "s":
		RunProto2Server()
	case "c":
		RunProto2Client()
	}
}

func RunProto2Client() {
	a := &microproto2.Animal{
		Id:         proto.Int(1),
		AnimalType: proto.String("Raptor"),
		Nickname:   proto.String("rapto"),
		Zone:       proto.Int(3),
		Age:        proto.Int(20),
	}
	data, err := proto.Marshal(a)
	if err != nil {
		log.Fatal(err)
	}
	SendData(data)
}

func SendData(data []byte) {
	c, err := net.Dial("tcp", "127.0.0.1:8282")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	c.Write(data)
}

func RunProto2Server() {
	l, err := net.Listen("tcp", ":8282")
	if err != nil {
		log.Fatal(err)
	}
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
			a := &microproto2.Animal{}
			proto.Unmarshal(data, a)
			fmt.Println(a)
		}(c)
	}
}
