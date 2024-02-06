package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reservation-service/data"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KeyProduct struct{}

type ReservationController struct {
	logger *log.Logger
	repo   *data.ReservationRepo
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

func (r ReservationController) isRoomReservedForDates(roomID string, checkInDate time.Time, checkOutDate time.Time) bool {

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

type Accommodation struct {
	Id           primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Name         string                 `json:"name"`
	Location     string                 `json:"location"`
	Amenities    []string               `json:"amenities"`
	MinGuests    int                    `json:"min_guests"`
	MaxGuests    int                    `json:"max_guests"`
	Images       []string               `json:"images"`
	Availability []AvailabilityInterval `json:"availability"`
	PriceType    string                 `json:"price_type"`
	HostID       primitive.ObjectID     `bson:"host_id,omitempty" json:"host_id"`
}

type AvailabilityInterval struct {
	Start           string  `json:"start"`
	End             string  `json:"end"`
	PricePerNight   float64 `json:"price_per_night"`
	PriceOnWeekends float64 `json:"price_on_weekends"`
}

func (r ReservationController) isRoomAvailableForDates(roomID string, checkInDate time.Time, checkOutDate time.Time) (bool, error) {
	reservationURL := fmt.Sprintf("http://accommodation-service:8001/accommodations/%s", roomID)

	accoResponse, err := http.Get(reservationURL)
	if err != nil {
		r.logger.Printf("Error making GET request for accommodation: %v", err)
		return false, fmt.Errorf("error making GET request for accommodation")
	}
	defer accoResponse.Body.Close()

	if accoResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(accoResponse.Body)
		r.logger.Println("error: ", string(body))
		return false, fmt.Errorf("accommodation service returned error: %s", string(body))
	}
	var accommodation Accommodation
	if err := json.NewDecoder(accoResponse.Body).Decode(&accommodation); err != nil {
		r.logger.Printf("error parsing accommodation response body: %v\n", err)
		return false, fmt.Errorf("error parsing accommodation response body")
	}

	for _, interval := range accommodation.Availability {
		start, err := time.Parse("2006-01-02", interval.Start)
		if err != nil {
			r.logger.Printf("error parsing start date: %v", err)
			continue
		}
		end, err := time.Parse("2006-01-02", interval.End)
		if err != nil {
			r.logger.Printf("error parsing end date: %v", err)
			continue
		}
		if (checkInDate.Equal(start) || checkInDate.After(start)) && (checkOutDate.Before(end) || checkOutDate.Equal(end)) {
			return true, nil
		}
	}

	return false, nil
}

func countWeekendDays(startDate, endDate time.Time) (int, error) {
	if startDate.After(endDate) {
		return 0, fmt.Errorf("start date must be before end date")
	}

	var weekendDays int
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		switch d.Weekday() {
		case time.Saturday, time.Sunday:
			weekendDays++
		}
	}

	return weekendDays, nil
}

func (r ReservationController) calculatePrice(roomID string, checkInDate, checkOutDate time.Time, number_of_guests int) (float64, error) {
	reservationURL := fmt.Sprintf("http://accommodation-service:8001/accommodations/%s", roomID)

	accoResponse, err := http.Get(reservationURL)
	if err != nil {
		r.logger.Printf("Error making GET request for accommodation: %v", err)
		return 0, fmt.Errorf("error making GET request for accommodation")
	}
	defer accoResponse.Body.Close()

	if accoResponse.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(accoResponse.Body)
		r.logger.Println("error: ", string(body))
		return 0, fmt.Errorf("accommodation service returned error: %s", string(body))
	}
	var accommodation Accommodation
	if err := json.NewDecoder(accoResponse.Body).Decode(&accommodation); err != nil {
		r.logger.Printf("error parsing accommodation response body: %v\n", err)
		return 0, fmt.Errorf("error parsing accommodation response body")
	}

	for _, interval := range accommodation.Availability {
		start, err := time.Parse("2006-01-02", interval.Start)
		if err != nil {
			r.logger.Printf("error parsing start date: %v", err)
			continue
		}
		end, err := time.Parse("2006-01-02", interval.End)
		if err != nil {
			r.logger.Printf("error parsing end date: %v", err)
			continue
		}
		if (checkInDate.Equal(start) || checkInDate.After(start)) && (checkOutDate.Before(end) || checkOutDate.Equal(end)) {
			price := 0.0
			pt := accommodation.PriceType
			gn := float64(number_of_guests)
			ppn := interval.PricePerNight
			pow := interval.PriceOnWeekends
			weekendDays, err := countWeekendDays(checkInDate, checkOutDate)
			d := (checkOutDate.Sub(checkInDate).Hours() / 24) - float64(weekendDays)

			if err != nil {
				return price, err
			}
			if pt == "Per_Guest" {
				price = gn*ppn*d + (pow * gn * float64(weekendDays))
			} else {
				price = ppn*d + (pow * float64(weekendDays))
			}
			return price, nil
		}
	}

	return 0, nil
}
func (r ReservationController) InsertReservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		type FieldsForInsertion struct {
			RoomId         string `json:"room_id"`
			NumberOfGuests int    `json:"number_of_guests"`
			CheckInDate    string `json:"checkin_date"`
			CheckOutDate   string `json:"checkout_date"`
		}

		var reservationBG data.ReservationByGuest
		var fields FieldsForInsertion

		claims, err := r.getUserInfoFromToken(c.Request.Header["Authorization"])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		userType, ok := claims["User_type"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user role not found in token"})
			return
		}
		userID, ok := claims["Uid"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user id not found in token"})
			return
		}
		username, ok := claims["First_name"].(string) //should be username
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username not found in token"})
			return
		}

		if userType != "Guest" {
			c.JSON(http.StatusForbidden, gin.H{"error": "only users that are guests can access this page"})
			return
		}

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
		if checkOutDate.Before(time.Now()) || checkInDate.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot make a reservation in the past"})
			return
		}

		if !r.isRoomReservedForDates(fields.RoomId, checkInDate, checkOutDate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "There is a reservation already for the selected dates"})
			return
		}
		isRoomAvailableForDates, err := r.isRoomAvailableForDates(fields.RoomId, checkInDate, checkOutDate)
		if !isRoomAvailableForDates {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The accommodation cannot be reserved for the selected dates"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		reservationBG.GuestID = userID
		reservationBG.RoomId = fields.RoomId
		reservationBG.GuestUsername = username
		reservationBG.CheckInDate = fields.CheckInDate
		reservationBG.CheckOutDate = fields.CheckOutDate
		reservationBG.NumberOfGuests = fields.NumberOfGuests
		price, err := r.calculatePrice(fields.RoomId, checkInDate, checkOutDate, fields.NumberOfGuests)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		reservationBG.Price = price

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

		claims, err := r.getUserInfoFromToken(c.Request.Header["Authorization"])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		userType, ok := claims["User_type"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user role not found in token"})
			return
		}
		if userType != "Guest" {
			c.JSON(http.StatusForbidden, gin.H{"error": "only users that are guests can access this page"})
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
		if resCheckInDate.Year() == now.Year() && resCheckInDate.Month() == now.Month() && resCheckInDate.Day() == now.Day() {
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
		guestIds, err := s.repo.GetDistinctIds("guest_id", "reservations_by_guest")
		if err != nil {
			s.logger.Print("Database exception: ", err)
		}

		if guestIds == nil {
			return
		}

		s.logger.Println(guestIds)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			s.logger.Fatal("Unable to convert to json :", err)
			return
		}
		c.JSON(http.StatusOK, guestIds)
	}
}
func (s *ReservationController) GetAllReservationIDs() gin.HandlerFunc {
	return func(c *gin.Context) {
		resIds, err := s.repo.GetDistinctIds("reservation_id", "reservations_by_guest")
		if err != nil {
			s.logger.Print("Database exception: ", err)
		}

		if resIds == nil {
			return
		}

		s.logger.Println(resIds)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			s.logger.Fatal("Unable to convert to json :", err)
			return
		}
		c.JSON(http.StatusOK, resIds)
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
func (r ReservationController) getUserInfoFromToken(authHeader []string) (jwt.MapClaims, error) {

	if len(authHeader) == 0 {
		return nil, fmt.Errorf("no header")
	}
	authString := strings.Join(authHeader, "")
	tokenString := strings.Split(authString, "Bearer ")[1]

	if len(tokenString) == 0 {
		return nil, fmt.Errorf("token empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
