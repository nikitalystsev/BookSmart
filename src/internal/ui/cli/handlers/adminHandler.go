package handlers

import (
	"BookSmart-services/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) deleteBook(c *gin.Context) {
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
	var newBook models.BookModel
	if err := c.BindJSON(&newBook); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input body"})
		return
	}

	if newBook.ID == uuid.Nil {
		newBook.ID = uuid.New()
	}

	err := h.bookService.Create(c.Request.Context(), &newBook)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}
