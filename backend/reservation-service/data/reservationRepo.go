package data

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocql/gocql"
)

type ReservationRepo struct {
	session *gocql.Session
	logger  *log.Logger
}

func New(logger *log.Logger) (*ReservationRepo, error) {
	db := os.Getenv("CASS_DB")

	// Connect to default keyspace
	cluster := gocql.NewCluster(db)
	cluster.Keyspace = "system"
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	// Create 'reservation' keyspace
	err = session.Query(
		fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
					WITH replication = {
						'class' : 'SimpleStrategy',
						'replication_factor' : %d
					}`, "reservations", 1)).Exec()
	if err != nil {
		logger.Println(err)
	}
	session.Close()

	cluster.Keyspace = "reservations"
	cluster.Consistency = gocql.One
	session, err = cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	// Return repository with logger and DB session
	return &ReservationRepo{
		session: session,
		logger:  logger,
	}, nil
}

// Disconnect from database
func (rr *ReservationRepo) CloseSession() {
	rr.session.Close()
}

func (rr *ReservationRepo) ConvertStringToUUID(id string) gocql.UUID {
	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		return gocql.UUID{}
	}
	return uuid
}
func (rr *ReservationRepo) CreateTables() {
	rr.logger.Printf("Creating table reservations_by_room")
	err := rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			room_id text,
			reservation_id UUID,
			guest_id text,
			guest_username text,
			reservation_date date,
			checkin_date date,
			checkout_date date,
			PRIMARY KEY (room_id, reservation_id)
		);
		`,
			"reservations_by_room")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	rr.logger.Printf("Creating table reservations_by_guest")
	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			guest_id text,
			guest_username text,
			reservation_id UUID,
			room_id text,
			reservation_date date,
			checkin_date date,
			checkout_date date,
			PRIMARY KEY (guest_id,reservation_id)
		);`,
			"reservations_by_guest")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}
}

// func (rr *ReservationRepo) DeletePastReservations() {
// 	err := rr.session.Query(
// 		`DELETE FROM reservations_by_room WHERE checkout_date < todate(now());`,
// 	).Exec()
// 	if err != nil {
// 		rr.logger.Printf("error deleting inactive reservations reservations_by_room: ", err)
// 	}
// 	err1 := rr.session.Query(
// 		`DELETE FROM reservations_by_guest WHERE checkout_date < todate(now());`,
// 	).Exec()
// 	if err1 != nil {
// 		rr.logger.Println("error deleting inactive reservations in reservations_by_guest: ", err)
// 		rr.logger.Println("Past reservations have been deleted")
// 	}
// }

func (rr *ReservationRepo) GetReservationsByRoom(id string) (ReservationsByRoom, error) {
	scanner := rr.session.Query(`SELECT 
	room_id, reservation_id, 
	guest_id,guest_username, reservation_date,
	checkin_date,checkout_date 
	FROM reservations_by_room
	 WHERE room_id = ? and checkout_date >= todate(now())`,
		id).Iter().Scanner()

	var reservations ReservationsByRoom
	for scanner.Next() {
		var reservation ReservationByRoom
		err := scanner.Scan(
			&reservation.RoomId, &reservation.ReservationId,
			&reservation.GuestID, &reservation.GuestUsername, &reservation.ReservationDate,
			&reservation.CheckInDate, &reservation.CheckOutDate)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		reservations = append(reservations, &reservation)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return reservations, nil
}

func (rr *ReservationRepo) GetReservationByIdAndRoom(resId string, roomId string) (ReservationByRoom, error) {
	resUuid := rr.ConvertStringToUUID(resId)
	scanner := rr.session.Query(`SELECT 
	room_id, reservation_id, 
	guest_id,guest_username, reservation_date,
	checkin_date,checkout_date 
	FROM reservations_by_room
	 WHERE reservation_id = ? and room_id = ? `,
		resUuid, roomId).Iter().Scanner()
	var reservation ReservationByRoom
	for scanner.Next() {
		err := scanner.Scan(
			&reservation.RoomId, &reservation.ReservationId,
			&reservation.GuestID, &reservation.GuestUsername, &reservation.ReservationDate,
			&reservation.CheckInDate, &reservation.CheckOutDate)
		if err != nil {
			rr.logger.Println(err)
			return reservation, err
		}
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return reservation, err
	}
	return reservation, nil
}

func (rr *ReservationRepo) GetReservationsByGuest(id string) (ReservationsByGuest, error) {
	scanner := rr.session.Query(`SELECT 
	guest_id,guest_username,reservation_id,
	room_id, 
	reservation_date,checkin_date,
	checkout_date 
	FROM reservations_by_guest
	 WHERE guest_id = ?`,
		id).Iter().Scanner()
	var reservations ReservationsByGuest
	for scanner.Next() {
		var reservation ReservationByGuest
		err := scanner.Scan(
			&reservation.GuestID, &reservation.GuestUsername, &reservation.ReservationId,
			&reservation.RoomId,
			&reservation.ReservationDate,
			&reservation.CheckInDate, &reservation.CheckOutDate)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		reservations = append(reservations, &reservation)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return reservations, nil
}
func (rr *ReservationRepo) InsertReservationByGuest(resRoom *ReservationByGuest) (string, error) {
	res_id := gocql.TimeUUID()
	today := time.Now().Format("2006-01-02")
	resRoom.ReservationDate = today
	err := rr.session.Query(
		`INSERT INTO reservations_by_guest (
			guest_id,
			guest_username, 
			reservation_id, 
			room_id, 
			reservation_date, 
			checkin_date, 
			checkout_date) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		resRoom.GuestID, resRoom.GuestUsername, res_id, resRoom.RoomId,
		resRoom.ReservationDate, resRoom.CheckInDate, resRoom.CheckOutDate).Exec()
	if err != nil {
		rr.logger.Println(err)
		return "", err
	}
	err = rr.session.Query(
		`INSERT INTO reservations_by_room (
			guest_id, 
			guest_username, 
			reservation_id,
			room_id, 
			reservation_date, 
			checkin_date, 
			checkout_date) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		resRoom.GuestID, resRoom.GuestUsername, res_id, resRoom.RoomId,
		resRoom.ReservationDate, resRoom.CheckInDate, resRoom.CheckOutDate).Exec()
	if err != nil {
		rr.logger.Println(err)
		return "", err
	}
	rr.logger.Println("reservation created , id: ", res_id)
	return res_id.String(), nil
}

func (rr *ReservationRepo) CancelReservation(resId string, roomId string) error {

	resUUID := rr.ConvertStringToUUID(resId)
	var guestID string
	scanner := rr.session.Query(
		`SELECT guest_id FROM reservations_by_room WHERE reservation_id = ? AND room_id = ?`,
		resUUID, roomId).Iter().Scanner()

	for scanner.Next() {
		if err := scanner.Scan(&guestID); err != nil {
			rr.logger.Println(err)
			return err
		}
	}
	err := rr.session.Query(
		`DELETE FROM reservations_by_room where reservation_id = ? and room_id= ?`,
		resUUID, roomId).Exec()
	if err != nil {
		rr.logger.Println(err)
		return err
	}

	err1 := rr.session.Query(
		`DELETE FROM reservations_by_guest where reservation_id = ? and guest_id = ?`,
		resUUID, guestID).Exec()
	if err1 != nil {
		rr.logger.Println(err)
		return err
	}

	return nil
}

func (rr *ReservationRepo) GetDistinctIds(idColumnName string, tableName string) ([]string, error) {
	scanner := rr.session.Query(
		fmt.Sprintf(`SELECT DISTINCT %s FROM %s`, idColumnName, tableName)).
		Iter().Scanner()
	var ids []string
	for scanner.Next() {
		var id string
		err := scanner.Scan(&id)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return ids, nil
}

func (rr *ReservationRepo) GetAllFromByGuest() (ReservationsByGuest, error) {
	scanner := rr.session.Query(
		(`SELECT * FROM reservations_by_guest`)).
		Iter().Scanner()
	var reservations ReservationsByGuest
	for scanner.Next() {
		var reservation ReservationByGuest
		err := scanner.Scan(&reservation.RoomId, &reservation.ReservationId, &reservation.GuestID, &reservation.ReservationDate, &reservation.CheckInDate, &reservation.CheckOutDate)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		reservations = append(reservations, &reservation)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return reservations, nil
}
func (rr *ReservationRepo) GetAccommodationIDByReservation(reservationID string) (string, error) {
	resUUID := rr.ConvertStringToUUID(reservationID)

	query := rr.session.Query(
		`SELECT room_id FROM reservations_by_room
		WHERE reservation_id = ?`,
		resUUID)

	var roomID string
	if err := query.Scan(&roomID); err != nil {
		rr.logger.Println(err)
		return "", err
	}

	// Replace the following line with the actual logic to get the accommodation ID
	// For now, let's assume roomID itself is the accommodation ID
	accommodationID := roomID

	return accommodationID, nil
}

/*func (rr *ReservationRepo) HasStayedInAccommodation(ctx context.Context, userID, accommodationID string) (bool, error) {
	// Convert string user and accommodation IDs to UUIDs
	userUUID := rr.ConvertStringToUUID(userID)
	accommodationUUID := rr.ConvertStringToUUID(accommodationID)

	// Query reservations_by_guest to check if the user has an active reservation
	query := rr.session.Query(
		`SELECT COUNT(*) FROM reservations_by_guest
		WHERE guest_id = ?
		AND reservation_id IN
			(SELECT reservation_id FROM reservations_by_room
			WHERE room_id = ? AND checkout_date >= todate(now()))`,
		userUUID, accommodationUUID)

	var count int
	if err := query.Scan(&count); err != nil {
		rr.logger.Println(err)
		return false, err
	}

	// If count is greater than 0, the user has an active reservation for the accommodation
	return count > 0, nil
}
*/
