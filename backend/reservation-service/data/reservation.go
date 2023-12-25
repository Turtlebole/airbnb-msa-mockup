package data

import (
	"encoding/json"
	"io"

	"github.com/gocql/gocql"
)

type ReservationByRoom struct {
	RoomId          string     `json:"room_id"`
	ReservationId   gocql.UUID `json:"reservation_id"`
	GuestID         string     `json:"guest_id"`
	GuestUsername   string     `json:"guest_username"`
	ReservationDate string     `json:"reservation_date"`
	CheckInDate     string     `json:"checkin_date"`
	CheckOutDate    string     `json:"checkout_date"`
}

type ReservationByGuest struct {
	RoomId          string     `json:"room_id"`
	ReservationId   gocql.UUID `json:"reservation_id"`
	GuestID         string     `json:"guest_id"`
	GuestUsername   string     `json:"guest_username"`
	ReservationDate string     `json:"reservation_date"`
	CheckInDate     string     `json:"checkin_date"`
	CheckOutDate    string     `json:"checkout_date"`
}

// Columns:
// reservation_id (UUID): Unique identifier for each reservation
// guest_name (text): Name of the guest making the reservation
// room_id (int): ID of the room reserved
// check_in_date (timestamp): Date and time of check-in
// check_out_date (timestamp): Date and time of check-out
// reservation_date (timestamp): Date and time when the reservation was made
// other_details (map<text, text>): Additional details as a map (e.g., special requests, preferences)

// Primary Key:
// reservation_	id

type ReservationsByRoom []*ReservationByRoom
type ReservationsByGuest []*ReservationByGuest

func (o *ReservationsByRoom) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}
func (o *ReservationsByGuest) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}
