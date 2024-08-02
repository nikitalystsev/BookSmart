package handlers

import (
	"BookSmart/internal/services/intfServices"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	bookService        intfServices.IBookService
	libCardService     intfServices.ILibCardService
	readerService      intfServices.IReaderService
	reservationService intfServices.IReservationService
	logger             *logrus.Entry
}

func NewHandler(
	bookService intfServices.IBookService,
	libCardService intfServices.ILibCardService,
	readerService intfServices.IReaderService,
	reservationService intfServices.IReservationService,
	logger *logrus.Entry) *Handler {
	return &Handler{
		bookService:        bookService,
		libCardService:     libCardService,
		readerService:      readerService,
		reservationService: reservationService,
		logger:             logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	//gin.DefaultWriter = ioutil.Discard

	router := gin.Default()

	api := router.Group("/api")
	{
		h.initReaderRoutes(api)
		h.initAdminRoutes(api)
	}

	return router
}
