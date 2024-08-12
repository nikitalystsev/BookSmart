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

func (h *Handler) signUp(c *gin.Context) {
	var inp dto.ReaderSignUpDTO
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid input body")
		return
	}

	reader := models.ReaderModel{
		ID:          uuid.New(),
		Fio:         inp.Fio,
		PhoneNumber: inp.PhoneNumber,
		Age:         inp.Age,
		Password:    inp.Password,
		Role:        "Reader",
	}

	err := h.readerService.SignUp(c.Request.Context(), &reader)
	if err != nil && errors.Is(err, errs.ErrReaderObjectIsNil) {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	if err != nil && errors.Is(err, errs.ErrReaderAlreadyExist) {
		c.AbortWithStatusJSON(http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) signIn(c *gin.Context) {
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

	c.JSON(http.StatusOK, dto.ReaderTokensDTO{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiredAt:    time.Now().Add(h.accessTokenTTL).UnixMilli(),
	})
}

func (h *Handler) refresh(c *gin.Context) {
	var inp string
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid input body")
		return
	}

	res, err := h.readerService.RefreshTokens(c.Request.Context(), inp)
	if err != nil && errors.Is(err, errs.ErrReaderDoesNotExists) {
		c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.ReaderTokensDTO{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiredAt:    time.Now().Add(h.accessTokenTTL).UnixMilli(),
	})
}

func (h *Handler) addToFavorites(c *gin.Context) {
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

	var bookID uuid.UUID
	if err = c.BindJSON(&bookID); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = h.readerService.AddToFavorites(c.Request.Context(), readerID, bookID)
	if err != nil && errors.Is(err, errs.ErrReaderDoesNotExists) {
		c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
		return
	}
	if err != nil && errors.Is(err, errs.ErrBookDoesNotExists) {
		c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
		return
	}
	if err != nil && errors.Is(err, errs.ErrBookAlreadyIsFavorite) {
		c.AbortWithStatusJSON(http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) getReaderByPhoneNumber(c *gin.Context) {
	phoneNumber := c.Param("phone_number")

	reader, err := h.readerService.GetByPhoneNumber(c.Request.Context(), phoneNumber)
	if err != nil && errors.Is(err, errs.ErrReaderDoesNotExists) {
		c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, reader)
}
