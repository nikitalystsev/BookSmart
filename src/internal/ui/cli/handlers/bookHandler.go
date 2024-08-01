package handlers

import (
	"BookSmart/internal/services/intfServices"
	"github.com/sirupsen/logrus"
)

type BookHandler struct {
	bookService intfServices.IBookService
	logger      *logrus.Entry
}

func NewBookHandler(bookService intfServices.IBookService, logger *logrus.Entry) *BookHandler {
	return &BookHandler{bookService: bookService, logger: logger}
}
