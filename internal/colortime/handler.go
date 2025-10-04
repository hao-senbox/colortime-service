package colortime

import (
	"colortime-service/helper"
	"colortime-service/pkg/constants"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ColorTimeHandler struct {
	ColorTimeService ColorTimeService
}

func NewColorTimeHandler(colorTimeService ColorTimeService) *ColorTimeHandler {
	return &ColorTimeHandler{
		ColorTimeService: colorTimeService,
	}
}

func (h *ColorTimeHandler) CreateColorTime(c *gin.Context) {

	var req CreateColorTimeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request body"), nil)
		return
	}

	tokenString, exist := c.Get(constants.Token)
	if !exist {
		helper.SendError(c, http.StatusUnauthorized, errors.New("unauthorized"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, tokenString)

	err := h.ColorTimeService.CreateColorTime(ctx, req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, errors.New("failed to create color time"), nil)
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "color time created successfully", nil)

}

func (h *ColorTimeHandler) GetColorTimes(c *gin.Context) {

	userID := c.Query("user_id")
	baseDate := c.Query("base_date")
	timeRange := c.Query("time_range")
	organizationID := c.Query("organization_id")

	if userID == "" || baseDate == "" || timeRange == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request parameters"), nil)
		return
	}

	tokenString, exist := c.Get(constants.Token)
	if !exist {
		helper.SendError(c, http.StatusUnauthorized, errors.New("unauthorized"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, tokenString)

	colorTimes, err := h.ColorTimeService.GetColorTimes(ctx, userID, organizationID, baseDate, timeRange)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, errors.New("failed to get color times"), nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "color times retrieved successfully", colorTimes)

}

func (h *ColorTimeHandler) GetColorTime(c *gin.Context) {
	
	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request parameters"), nil)
		return
	}

	tokenString, exist := c.Get(constants.Token)
	if !exist {
		helper.SendError(c, http.StatusUnauthorized, errors.New("unauthorized"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, tokenString)

	colorTime, err := h.ColorTimeService.GetColorTime(ctx, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, errors.New("failed to get color time"), nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "color time retrieved successfully", colorTime)

}

func (h *ColorTimeHandler) UpdateColorTime(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request parameters"), nil)
		return
	}

	var req UpdateColorTimeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request body"), nil)
		return
	}

	tokenString, exist := c.Get(constants.Token)
	if !exist {
		helper.SendError(c, http.StatusUnauthorized, errors.New("unauthorized"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, tokenString)

	err := h.ColorTimeService.UpdateColorTime(ctx, req, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, errors.New("failed to update color time"), nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "color time updated successfully", nil)

}

func (h *ColorTimeHandler) DeleteColorTime(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request parameters"), nil)
		return
	}

	tokenString, exist := c.Get(constants.Token)
	if !exist {
		helper.SendError(c, http.StatusUnauthorized, errors.New("unauthorized"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, tokenString)

	err := h.ColorTimeService.DeleteColorTime(ctx, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, errors.New("failed to delete color time"), nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "color time deleted successfully", nil)
	
}