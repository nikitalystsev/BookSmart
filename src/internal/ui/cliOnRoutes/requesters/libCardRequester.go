package requesters

import (
	"github.com/sirupsen/logrus"
)

type LibCardRequester struct {
	logger *logrus.Entry
}

func NewLibCardRequester(logger *logrus.Entry) *LibCardRequester {
	return &LibCardRequester{logger: logger}
}
