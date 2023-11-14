package main

import (
	"context"
	"github.com/SimikicAleksandar/AccomodationService/data"
	"github.com/SimikicAleksandar/AccomodationService/handlers"
	protos "github.com/SimikicAleksandar/AccomodationService/protos/main"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	lis, err := net.Listen("tcp", ":9093")
	if err != nil {
		log.Fatal(err)
	}
	serverRegister := grpc.NewServer()

	logger := log.New(os.Stdout, "[accommodation-main] ", log.LstdFlags)

	//Initialize the handler and inject said logger
	accommodationLog := log.New(os.Stdout, "[accommodation-repo-log] ", log.LstdFlags)
	accommodationRepo, err := data.New(timeoutContext, accommodationLog)
	service := handlers.NewServer(logger, accommodationRepo)

	protos.RegisterAccommodationServer(serverRegister, service)
	err = serverRegister.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
