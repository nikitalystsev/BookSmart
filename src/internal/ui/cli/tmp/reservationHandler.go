package tmp

import (
	"BookSmart/internal/services/intfServices"
	"github.com/sirupsen/logrus"
)

type ReservationHandler struct {
	reservationService intfServices.IReservationService
	logger             *logrus.Entry
}

func NewReservationHandler(
	reservationService intfServices.IReservationService,
	logger *logrus.Entry,
) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService, logger: logger}
}
