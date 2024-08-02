package handlers

import (
	"BookSmart/internal/models"
	"BookSmart/internal/services/errsService"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initReaderRoutes(api *gin.RouterGroup) {
	readers := api.Group("/readers")
	{
		readers.POST("/sign-up", h.readerSignUp)
		readers.POST("/sign-in")
		readers.POST("/auth/refresh")
	}
}

func (h *Handler) readerSignUp(c *gin.Context) {
	fmt.Println("call readerSignUp")

	var inp models.ReaderModel
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid input body")
		return
	}

	err := h.readerService.SignUp(c.Request.Context(), &inp)
	if err != nil && !errors.Is(err, errsService.ErrReaderAlreadyExist) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	if err != nil && errors.Is(err, errsService.ErrReaderAlreadyExist) {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}
