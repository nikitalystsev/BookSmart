package handlers

import (
	"BookSmart-services/dto"
	"BookSmart-services/errs"
	"BookSmart-services/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) signInAsAdmin(c *gin.Context) {
	var inp dto.ReaderSignInDTO
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.readerService.SignIn(c.Request.Context(), &inp)
	if err != nil && !errors.Is(err, errs.ErrReaderDoesNotExists) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	if err != nil && errors.Is(err, errs.ErrReaderDoesNotExists) {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
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

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

func (h *Handler) deleteBook(c *gin.Context) {
	_, readerRole, err := getReaderData(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	if readerRole == "Reader" {
		c.AbortWithStatusJSON(http.StatusForbidden, "reader not delete book")
		return
	}

	bookID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = h.bookService.Delete(c.Request.Context(), bookID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) addNewBook(c *gin.Context) {
	_, readerRole, err := getReaderData(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	if readerRole == "Reader" {
		c.AbortWithStatusJSON(http.StatusForbidden, "reader not delete book")
		return
	}

	var newBook models.BookModel
	if err = c.BindJSON(&newBook); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input body"})
		return
	}

	if newBook.ID == uuid.Nil {
		newBook.ID = uuid.New()
	}

	err = h.bookService.Create(c.Request.Context(), &newBook)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}
