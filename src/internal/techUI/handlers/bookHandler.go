package handlers

import (
	"BookSmart-services/core/dto"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) getBooks(c *gin.Context) {
	var params dto.BookParamsDTO
	if err := c.BindJSON(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input body"})
		return
	}

	fmt.Println(params)
	books, err := h.bookService.GetByParams(c.Request.Context(), &params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, books)
}

func (h *Handler) getBookByID(c *gin.Context) {
	bookID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	book, err := h.bookService.GetByID(c.Request.Context(), bookID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get books"})
		return
	}

	c.JSON(http.StatusOK, book)
}
