package routes

import (
	"reservation-service/controllers"

	"github.com/gin-gonic/gin"
)

func MainRoutes(routes *gin.Engine, reservationController controllers.ReservationController) {

	routes.GET("/reservations/by_room/:room_id", reservationController.GetReservationByRoom())
	routes.GET("/reservations/by_guest/:guest_id", reservationController.GetReservationByGuest())
	routes.POST("/reservations/by_guest/insert", reservationController.InsertReservation())
	routes.DELETE("/reservations/cancel", reservationController.CancelReservation())
	routes.GET("/reservations/all-guest-ids", reservationController.GetAllGuestIDs())
	routes.GET("/reservations/all-reservation-ids", reservationController.GetAllReservationIDs())
	routes.GET("/reservations/all-res-guest", reservationController.GetAllFromByGuest())
}
