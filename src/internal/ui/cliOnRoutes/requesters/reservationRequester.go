package requesters

import (
	"github.com/sirupsen/logrus"
)

type ReservationRequester struct {
	logger *logrus.Entry
}

func NewReservationRequester(logger *logrus.Entry) *ReservationRequester {
	return &ReservationRequester{logger: logger}
}
