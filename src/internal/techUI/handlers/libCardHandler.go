package handlers

import (
	"BookSmart-services/core/models"
	"BookSmart-services/errs"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) createLibCard(c *gin.Context) {
	readerIDStr, _, err := getReaderData(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	readerID, err := uuid.Parse(readerIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = h.libCardService.Create(c.Request.Context(), readerID)
	if err != nil && errors.Is(err, errs.ErrLibCardAlreadyExist) {
		c.AbortWithStatusJSON(http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) updateLibCard(c *gin.Context) {
	readerIDStr, _, err := getReaderData(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	readerID, err := uuid.Parse(readerIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	var libCard *models.LibCardModel
	libCard, err = h.libCardService.GetByReaderID(c.Request.Context(), readerID)
	if err != nil && !errors.Is(err, errs.ErrLibCardDoesNotExists) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	if err != nil && errors.Is(err, errs.ErrLibCardDoesNotExists) {
		c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
		return
	}

	err = h.libCardService.Update(c.Request.Context(), libCard)
	if err != nil && errors.Is(err, errs.ErrLibCardDoesNotExists) {
		c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
		return
	}
	if err != nil && errors.Is(err, errs.ErrLibCardIsValid) {
		c.AbortWithStatusJSON(http.StatusConflict, err.Error())
		return
	}
	if err != nil && errors.Is(err, errs.ErrLibCardObjectIsNil) {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) getLibCardByReaderID(c *gin.Context) {
	readerIDStr, _, err := getReaderData(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	readerID, err := uuid.Parse(readerIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	var libCard *models.LibCardModel
	libCard, err = h.libCardService.GetByReaderID(c.Request.Context(), readerID)
	if err != nil && !errors.Is(err, errs.ErrLibCardDoesNotExists) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	if err != nil && errors.Is(err, errs.ErrLibCardDoesNotExists) {
		c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, libCard)
}
