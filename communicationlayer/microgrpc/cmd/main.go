package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/phk13/poc-micro/communicationlayer/microgrpc"
	"github.com/phk13/poc-micro/databaselayer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	op := flag.String("op", "s", "s for server, c for client")
	flag.Parse()
	switch strings.ToLower(*op) {
	case "s":
		runGRPCServer()
	case "c":
		runGRPCClient()
	}
}

func runGRPCServer() {
	grpclog.Infoln("Starting GRPC Server")
	lis, err := net.Listen("tcp", ":8282")
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	microServer, err := microgrpc.NewMicroGrpcServer(databaselayer.MONGODB, "mongodb://172.17.0.8")
	if err != nil {
		log.Fatal(err)
	}
	microgrpc.RegisterMicroServiceServer(grpcServer, microServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func runGRPCClient() {
	conn, err := grpc.Dial("127.0.0.1:8282", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := microgrpc.NewMicroServiceClient(conn)
	input := ""
	fmt.Println("All animals? (y/n): ")
	fmt.Scanln(&input)
	if strings.EqualFold(input, "y") {
		animals, err := client.GetAllAnimals(context.Background(), &microgrpc.Request{})

		if err != nil {
			log.Fatal(err)
		}
		for {
			animal, err := animals.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				grpclog.Fatal(err)
			}
			fmt.Println(animal)
			grpclog.Infoln(animal)
		}
		return
	}
	fmt.Println("Nickname? ")
	fmt.Scanln(&input)
	a, err := client.GetAnimal(context.Background(), &microgrpc.Request{Nickname: input})
	if err != nil {
		log.Fatal(err)
	}
	grpclog.Infoln(*a)
}
