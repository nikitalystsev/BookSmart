package models

import (
	"github.com/nikitalystsev/BookSmart-services/core/models"
	tdbmodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/testDataBuilder/models"
)

type RatingModelObjectMother struct {
}

func NewRatingModelObjectMother() *RatingModelObjectMother {
	return &RatingModelObjectMother{}
}

func (rmom *RatingModelObjectMother) DefaultRating() *models.RatingModel {
	return tdbmodels.NewRatingModelBuilder().Build()
}
