package handlers

import (
	"BookSmart-services/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) createLibCard(c *gin.Context) {
	readerIDStr, err := getReaderID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	readerID, err := uuid.Parse(readerIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = h.libCardService.Create(c.Request.Context(), readerID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) updateLibCard(c *gin.Context) {
	readerIDStr, err := getReaderID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	readerID, err := uuid.Parse(readerIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	var libCard *models.LibCardModel
	libCard, err = h.libCardService.GetByReaderID(c.Request.Context(), readerID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = h.libCardService.Update(c.Request.Context(), libCard)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) getLibCardByReaderID(c *gin.Context) {
	readerIDStr, err := getReaderID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	readerID, err := uuid.Parse(readerIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	var libCard *models.LibCardModel
	libCard, err = h.libCardService.GetByReaderID(c.Request.Context(), readerID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, libCard)
}
