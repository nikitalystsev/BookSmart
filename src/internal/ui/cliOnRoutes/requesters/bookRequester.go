package requesters

import (
	"github.com/sirupsen/logrus"
)

type BookRequester struct {
	logger *logrus.Entry
}

func NewBookRequester(logger *logrus.Entry) *BookRequester {
	return &BookRequester{logger: logger}
}
