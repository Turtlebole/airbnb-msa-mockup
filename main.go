package main

import (
	"github.com/SimikicAleksandar/AccomodationService/handlers"
	protos "github.com/SimikicAleksandar/AccomodationService/protos/main"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {

	lis, err := net.Listen("tcp", ":9093")
	if err != nil {
		log.Fatal(err)
	}
	serverRegister := grpc.NewServer()

	logger := log.New(os.Stdout, "[accommodation-glavno] ", log.LstdFlags)

	//Initialize the handler and inject said logger
	service := handlers.NewServer(logger)

	protos.RegisterAccommodationServer(serverRegister, service)
	err = serverRegister.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
