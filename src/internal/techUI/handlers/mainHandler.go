package handlers

import (
	"BookSmart-services/intf"
	"BookSmart-services/pkg/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Handler struct {
	bookService        intf.IBookService
	libCardService     intf.ILibCardService
	readerService      intf.IReaderService
	reservationService intf.IReservationService
	tokenManager       auth.ITokenManager
	logger             *logrus.Entry
}

func NewHandler(
	bookService intf.IBookService,
	libCardService intf.ILibCardService,
	readerService intf.IReaderService,
	reservationService intf.IReservationService,
	tokenManager auth.ITokenManager,
	logger *logrus.Entry,
) *Handler {
	return &Handler{
		bookService:        bookService,
		libCardService:     libCardService,
		readerService:      readerService,
		reservationService: reservationService,
		tokenManager:       tokenManager,
		logger:             logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	router := gin.Default()

	router.Use(h.corsSettings())

	authenticate := router.Group("/auth")
	{
		authenticate.POST("/sign-up", h.signUp)
		authenticate.POST("/sign-in", h.signIn)
		authenticate.POST("/admin/sign-in", h.signInAsAdmin)
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

		admin := api.Group("/admin")
		{
			admin.POST("/books/:id", h.deleteBook)
			admin.POST("/books", h.addNewBook)
		}
	}

	return router
}

func (h *Handler) corsSettings() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowMethods: []string{
			http.MethodPost,
		},
		AllowOrigins: []string{
			"*",
		},
		AllowCredentials: true,
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
		},
		ExposeHeaders: []string{
			"Content-Type",
		},
	})
}
