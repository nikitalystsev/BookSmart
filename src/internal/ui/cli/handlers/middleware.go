package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	readerCtx           = "readerID"
)

func (h *Handler) readerIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
	}

	c.Set(readerCtx, id)
}

func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.tokenManager.Parse(headerParts[1])
}

func (h *Handler) test(c *gin.Context) {
	var id uuid.UUID
	if err := c.BindJSON(&id); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid input body")
		return
	}

	c.Set("id", id.String())
}

func getReaderID(c *gin.Context) (string, error) {
	id, ok := c.Get(readerCtx)
	if !ok {
		return "", errors.New("user id not found")
	}

	idStr, ok := id.(string)
	if !ok {
		return "", errors.New("user id is of invalid type")
	}

	return idStr, nil
}
