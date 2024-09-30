package models

import (
	"github.com/nikitalystsev/BookSmart-services/core/models"
	tdbmodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/testDataBuilder/models"
	"time"
)

type ReservationModelObjectMother struct {
}

func NewReservationModelObjectMother() *ReaderModelObjectMother {
	return &ReaderModelObjectMother{}
}

func (rmom *ReaderModelObjectMother) DefaultReservation() *models.ReservationModel {
	return tdbmodels.NewReservationModelBuilder().Build()
}

func (rmom *ReaderModelObjectMother) ExpiredReservation() *models.ReservationModel {
	return tdbmodels.NewReservationModelBuilder().
		WithReturnDate(time.Now().AddDate(0, -1, 0)).
		Build()
}
