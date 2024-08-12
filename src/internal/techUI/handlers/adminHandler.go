package handlers

import (
	"BookSmart-services/core/dto"
	"BookSmart-services/core/models"
	"BookSmart-services/errs"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (h *Handler) signInAsAdmin(c *gin.Context) {
	var inp dto.ReaderSignInDTO
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.readerService.SignIn(c.Request.Context(), &inp)
	if err != nil && errors.Is(err, errs.ErrReaderDoesNotExists) {
		c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
		return
	}
	if err != nil && errors.Is(err, errors.New("wrong password")) {
		c.AbortWithStatusJSON(http.StatusConflict, err.Error())
		return
	}
	if err != nil && errors.Is(err, errs.ErrReaderObjectIsNil) {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	_, role, err := h.tokenManager.Parse(res.AccessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	if role == "Reader" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "you are not authorized to perform this action")
		return
	}

	c.JSON(http.StatusOK, dto.ReaderTokensDTO{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiredAt:    time.Now().Add(h.accessTokenTTL).UnixMilli(),
	})
}

func (h *Handler) deleteBook(c *gin.Context) {
	_, readerRole, err := getReaderData(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if readerRole == "Reader" {
		c.AbortWithStatusJSON(http.StatusForbidden, "reader not delete book")
		return
	}

	bookID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = h.bookService.Delete(c.Request.Context(), bookID)
	if err != nil && errors.Is(err, errs.ErrBookObjectIsNil) {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	if err != nil && errors.Is(err, errs.ErrBookDoesNotExists) {
		c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) addNewBook(c *gin.Context) {
	_, readerRole, err := getReaderData(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if readerRole == "Reader" {
		c.AbortWithStatusJSON(http.StatusForbidden, "reader not delete book")
		return
	}

	var newBook dto.BookParamsDTO
	if err = c.BindJSON(&newBook); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	book := &models.BookModel{
		ID:             uuid.New(),
		Title:          newBook.Title,
		Author:         newBook.Author,
		Publisher:      newBook.Publisher,
		CopiesNumber:   newBook.CopiesNumber,
		Rarity:         newBook.Rarity,
		Genre:          newBook.Genre,
		PublishingYear: newBook.PublishingYear,
		Language:       newBook.Language,
		AgeLimit:       newBook.AgeLimit,
	}

	err = h.bookService.Create(c.Request.Context(), book)
	if err != nil && errors.Is(err, errs.ErrBookObjectIsNil) {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) getReservationsByBookID(c *gin.Context) {
	var bookID uuid.UUID
	if err := c.BindJSON(&bookID); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	reservations, err := h.reservationService.GetByBookID(c.Request.Context(), bookID)
	if err != nil && errors.Is(err, errs.ErrReservationDoesNotExists) {
		c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, reservations)
}
