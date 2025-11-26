package colortime

import (
	"colortime-service/helper"
	"colortime-service/pkg/constants"
	"context"
	"errors"
	"fmt"
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

func (h *ColorTimeHandler) AddTopicToColorTimeWeek(c *gin.Context) {

	var req AddTopicToColorTimeWeekRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request parameters"), nil)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), nil)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.ColorTimeService.AddTopicToColorTimeWeek(ctx, id, &req, userID.(string))

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "topic added to color time week successfully", nil)

}

func (h *ColorTimeHandler) GetToColorTimeWeek(c *gin.Context) {

	userID := c.Query("user_id")
	if userID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request parameters"), nil)
		return
	}

	role := c.Query("role")
	if role == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request parameters"), nil)
		return
	}

	orgID := c.Query("org_id")
	if orgID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request parameters"), nil)
		return
	}

	start := c.Query("start")
	if start == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request parameters"), nil)
		return
	}

	end := c.Query("end")
	if end == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request parameters"), nil)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.ColorTimeService.GetColorTimeWeek(ctx, userID, role, orgID, start, end)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "topic to color time week fetched successfully", data)

}

func (h *ColorTimeHandler) DeleteTopicToColorTimeWeek(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request parameters"), nil)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.ColorTimeService.DeleteTopicToColorTimeWeek(ctx, id)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "topic deleted from color time week successfully", nil)

}

func (h *ColorTimeHandler) AddTopicToColorTimeDay(c *gin.Context) {

	var req AddTopicToColorTimeDayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request parameters"), nil)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.ColorTimeService.AddTopicToColorTimeDay(ctx, id, &req)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "topic added to color time day successfully", nil)

}

func (h *ColorTimeHandler) DeleteTopicToColorTimeDay(c *gin.Context) {

	var req DeleteTopicToColorTimeDayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request parameters"), nil)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.ColorTimeService.DeleteTopicToColorTimeDay(ctx, id, &req)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "topic deleted from color time day successfully", nil)

}

func (h *ColorTimeHandler) UpdateColorSlotHandler(c *gin.Context) {
	weekColorTimeID := c.Param("week_colortime_id")
	slotID := c.Param("slot_id")

	var req UpdateColorSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}
	
	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, errors.New("user ID not found in context"), nil)
		return
	}
	
	if userID == "" {
		helper.SendError(c, http.StatusUnauthorized, errors.New("user ID not found in context"), nil)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.ColorTimeService.UpdateColorSlot(ctx, weekColorTimeID, slotID, &req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "color slot updated successfully", nil)
	
}

func (h *ColorTimeHandler) GetColorTimeDay(c *gin.Context) {
	orgID := c.Query("org_id")
	if orgID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("org_id is required"), nil)
		return
	}

	date := c.Query("date")
	if date == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("date is required"), nil)
		return
	}

	userID := c.Query("user_id")
	if userID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("user_id is required"), nil)
		return
	}

	role := c.Query("role")
	if role == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("role is required"), nil)
		return
	}
	
	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.ColorTimeService.GetColorTimeDay(ctx, orgID, date, userID, role)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "color time day retrieved successfully", data)
}

func (h *ColorTimeHandler) GetTopicByTerm(c *gin.Context) {
	orgID := c.Query("org_id")
	if orgID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("org_id is required"), nil)
		return
	}

	userID := c.Query("user_id")
	if userID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("user_id is required"), nil)
		return
	}

	role := c.Query("role")
	if role == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("role is required"), nil)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.ColorTimeService.GetTopicByTerm(ctx, orgID, userID, role)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "topic by term retrieved successfully", data)
}