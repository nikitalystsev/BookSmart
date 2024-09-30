package models

import (
	"github.com/nikitalystsev/BookSmart-services/core/models"
	tdbmodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/testDataBuilder/models"
	"time"
)

type LibCardModelObjectMother struct {
}

func NewLibCardModelObjectMother() *LibCardModelObjectMother {
	return &LibCardModelObjectMother{}
}

func (lcmom *LibCardModelObjectMother) DefaultLibCard() *models.LibCardModel {
	return tdbmodels.NewLibCardModelBuilder().Build()
}

func (lcmom *LibCardModelObjectMother) ExpiredLibCard() *models.LibCardModel {
	return tdbmodels.NewLibCardModelBuilder().
		WithIssueDate(time.Now().AddDate(0, 0, -370)).
		WithActionStatus(false).
		Build()
}
