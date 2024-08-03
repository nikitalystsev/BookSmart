package handlers

import (
	"BookSmart/internal/services/intfServices"
	"BookSmart/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

type Handler struct {
	bookService        intfServices.IBookService
	libCardService     intfServices.ILibCardService
	readerService      intfServices.IReaderService
	reservationService intfServices.IReservationService
	logger             *logrus.Entry
	tokenManager       auth.ITokenManager
}

func NewHandler(
	bookService intfServices.IBookService,
	libCardService intfServices.ILibCardService,
	readerService intfServices.IReaderService,
	reservationService intfServices.IReservationService,
	logger *logrus.Entry,
	tokenManager auth.ITokenManager,
) *Handler {
	return &Handler{
		bookService:        bookService,
		libCardService:     libCardService,
		readerService:      readerService,
		reservationService: reservationService,
		logger:             logger,
		tokenManager:       tokenManager,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	router := gin.Default()

	authenticate := router.Group("/auth")
	{
		authenticate.POST("/sign-up", h.signUp)
		authenticate.POST("/sign-in", h.signIn)
		authenticate.POST("/refresh", h.refresh)
	}

	general := router.Group("/general")
	{
		general.GET("/books", h.getBooks)
		general.GET("/books/:id", h.getBookByID)

	}

	api := router.Group("/api", h.readerIdentity)
	{
		api.POST("/favorites", h.addToFavorites)

		libCards := api.Group("/lib-cards")
		{
			libCards.POST("/", h.createLibCard)
			libCards.PUT("/", h.updateLibCard)
			libCards.GET("/", h.getLibCardByReaderID)
		}

		reservations := api.Group("/reservations")
		{
			reservations.POST("/", h.reserveBook)
			reservations.GET("/", h.getReservationsByReaderID)
			reservations.PUT("/:id", h.updateReservation)
		}
	}

	return router
}
