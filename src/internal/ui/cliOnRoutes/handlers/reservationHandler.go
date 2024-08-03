package handlers

import (
	"BookSmart/internal/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) reserveBook(c *gin.Context) {
	var inp dto.ReserveBookDTO
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid input body")
		return
	}

	err := h.reservationService.Create(c.Request.Context(), inp.ReaderID, inp.BookID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}
