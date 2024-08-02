package tmp

import (
	"BookSmart/internal/services/intfServices"
	"github.com/sirupsen/logrus"
)

type LibCardHandler struct {
	libCardService intfServices.ILibCardService
	logger         *logrus.Entry
}

func NewLibCardHandler(libCardService intfServices.ILibCardService, logger *logrus.Entry) *LibCardHandler {
	return &LibCardHandler{libCardService: libCardService, logger: logger}
}
