package controllers

import (
	"log"
	"net/http"
	"reservation-service/data"
	"time"

	"github.com/gin-gonic/gin"
)

type KeyProduct struct{}

type ReservationController struct {
	logger *log.Logger
	// NoSQL: injecting student repository
	repo *data.ReservationRepo
}

func NewReservationController(l *log.Logger, r *data.ReservationRepo) *ReservationController {
	return &ReservationController{l, r}
}

func (r ReservationController) GetReservationByRoom() gin.HandlerFunc {
	return func(c *gin.Context) {

		roomId := c.Param("room_id")

		resBR, err := r.repo.GetReservationsByRoom(roomId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Database exception": err.Error()})
			return
		}

		if resBR == nil {
			c.JSON(http.StatusNotFound, "No reservations found for this room")
			return
		}

		c.JSON(http.StatusOK, resBR)
	}
}

func (r ReservationController) GetReservationByGuest() gin.HandlerFunc {
	return func(c *gin.Context) {
		guestID := c.Param("guest_id")
		r.logger.Println("get reservation by guest_id: ", guestID)

		resBG, err := r.repo.GetReservationsByGuest(guestID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Database exception": err.Error()})
			return
		}

		if resBG == nil {
			c.JSON(http.StatusNotFound, gin.H{"Reservations not found for id:": guestID})
			return
		}

		c.JSON(http.StatusOK, resBG)
	}
}

func (r ReservationController) IsRoomAvailableForDates(roomID string, checkInDate time.Time, checkOutDate time.Time) bool {

	reservations, err := r.repo.GetReservationsByRoom(roomID)
	if err != nil {
		return false
	}

	for _, reservation := range reservations {
		resCheckInDate, err := time.Parse("2006-01-02", reservation.CheckInDate)
		if err != nil {
			return false
		}

		resCheckOutDate, err := time.Parse("2006-01-02", reservation.CheckOutDate)
		if err != nil {
			return false
		}

		if (checkInDate.Before(resCheckOutDate) || checkInDate.Equal(resCheckOutDate)) &&
			(checkOutDate.After(resCheckInDate) || checkOutDate.Equal(resCheckInDate)) {

			return false
		}
	}
	return true
}

func (r ReservationController) InsertReservationByGuest() gin.HandlerFunc {
	return func(c *gin.Context) {
		type FieldsForInsertion struct {
			RoomId         string `json:"room_id"`
			GuestID        string `json:"guest_id"`
			Guest_Username string `json:"guest_username"`
			RoomName       string `json:"room_name"`
			CheckInDate    string `json:"checkin_date"`
			CheckOutDate   string `json:"checkout_date"`
		}
		var reservationBG data.ReservationByGuest
		var fields FieldsForInsertion

		if err := c.BindJSON(&fields); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Data you provided causes an error": err.Error()})
			r.logger.Println(err)
			return
		}
		checkInDate, err := time.Parse("2006-01-02", fields.CheckInDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error parsing checkin_date": err.Error()})
			return
		}

		checkOutDate, err := time.Parse("2006-01-02", fields.CheckOutDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error parsing checkout_date": err.Error()})
			return
		}

		if checkOutDate.Before(checkInDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Checkout date must be after check-in date"})
			return
		}
		if !r.IsRoomAvailableForDates(fields.RoomId, checkInDate, checkOutDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "there is a reservation already for the selected dates"})
			return
		}

		reservationBG.GuestID = fields.GuestID
		reservationBG.RoomId = fields.RoomId
		reservationBG.RoomName = fields.RoomName
		reservationBG.GuestUsername = fields.Guest_Username
		reservationBG.CheckInDate = fields.CheckInDate
		reservationBG.CheckOutDate = fields.CheckOutDate

		res_id, err := r.repo.InsertReservationByGuest(&reservationBG)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Database exception": err.Error()})
			r.logger.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"Reservation created": res_id})

	}
}

func (r ReservationController) CancelReservation() gin.HandlerFunc {
	return func(c *gin.Context) {

		type CancelInfo struct {
			RoomId        string `json:"room_id"`
			ReservationId string `json:"reservation_id"`
		}

		var cancelInfo CancelInfo

		if err := c.BindJSON(&cancelInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Data you provided causes an error": err.Error()})
			r.logger.Println(err)
			return
		}

		reservation, err := r.repo.GetReservationByIdAndRoom(cancelInfo.ReservationId, cancelInfo.RoomId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		now := time.Now()
		resCheckInDate, err := time.Parse("2006-01-02", reservation.CheckInDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if resCheckInDate.Year() == now.Year() && resCheckInDate.Month() == now.Month() && resCheckInDate.Day() == now.Day() || resCheckInDate.After(now) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You can only cancel before your reservation date"})
			return
		}
		err = r.repo.CancelReservation(cancelInfo.ReservationId, cancelInfo.RoomId)
		if err != nil {
			r.logger.Print("Database exception: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"Database exception": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "Reservation Canceled")
	}
}
func (s *ReservationController) GetAllGuestIDs() gin.HandlerFunc {
	return func(c *gin.Context) {
		studentIds, err := s.repo.GetDistinctIds("guest_id", "reservations_by_guest")
		if err != nil {
			s.logger.Print("Database exception: ", err)
		}

		if studentIds == nil {
			return
		}

		s.logger.Println(studentIds)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			s.logger.Fatal("Unable to convert to json :", err)
			return
		}
		c.JSON(http.StatusOK, studentIds)
	}
}
func (s *ReservationController) GetAllReservationIDs() gin.HandlerFunc {
	return func(c *gin.Context) {
		studentIds, err := s.repo.GetDistinctIds("reservation_id", "reservations_by_guest")
		if err != nil {
			s.logger.Print("Database exception: ", err)
		}

		if studentIds == nil {
			return
		}

		s.logger.Println(studentIds)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			s.logger.Fatal("Unable to convert to json :", err)
			return
		}
		c.JSON(http.StatusOK, studentIds)
	}
}
func (r ReservationController) GetAllFromByGuest() gin.HandlerFunc {
	return func(c *gin.Context) {

		reservationsByGuest, err := r.repo.GetAllFromByGuest()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Database exception": err.Error()})
			return
		}

		if reservationsByGuest == nil {
			return
		}

		err = c.BindJSON(&reservationsByGuest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable To convert to JSON"})
			return
		}
		c.JSON(http.StatusOK, reservationsByGuest)
	}
}
