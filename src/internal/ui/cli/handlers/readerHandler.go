package handlers

import (
	"BookSmart/internal/services/intfServices"
	"github.com/sirupsen/logrus"
)

type ReaderHandler struct {
	readerService intfServices.IReaderService
	logger        *logrus.Entry
}

func NewReaderHandler(readerService intfServices.IReaderService, logger *logrus.Entry) *ReaderHandler {
	return &ReaderHandler{readerService: readerService, logger: logger}
}
