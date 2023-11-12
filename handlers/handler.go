package handlers

import (
	"context"
	protos "github.com/SimikicAleksandar/AccomodationService/protos/main"
	"log"
)

type MyAccommodationServer struct {
	protos.UnimplementedAccommodationServer
	logger *log.Logger
	// NoSQL: injecting product repository
}

func NewServer(l *log.Logger) *MyAccommodationServer {
	return &MyAccommodationServer{*new(protos.UnimplementedAccommodationServer), l}
}

func GetAccommodation(context.Context, *protos.AccommodationRequest) (*protos.DummyList, error) {
	return nil, nil
}
func SetAccommodation(context.Context, *protos.AccommodationResponse) (*protos.Emptya, error) {
	return nil, nil
}
func UpdateAccommodation(context.Context, *protos.AccommodationResponse) (*protos.Emptya, error) {
	return nil, nil
}
