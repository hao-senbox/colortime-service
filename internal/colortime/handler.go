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

func (h *ColorTimeHandler) CreateColorTime(c *gin.Context) {

	var req CreateColorTimeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
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

	err := h.ColorTimeService.CreateColorTime(ctx, &req, userID.(string))

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "color time created successfully", nil)

}

func (h *ColorTimeHandler) CreateTemplateColorTime(c *gin.Context) {

	var req CreateTemplateColorTimeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
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

	err := h.ColorTimeService.CreateTemplateColorTime(ctx, &req, userID.(string))

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "template color time created successfully", nil)
}

func (h *ColorTimeHandler) GetColorTimes(c *gin.Context) {

	start := c.Query("start")
	end := c.Query("end")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	colorTimes, err := h.ColorTimeService.GetColorTimes(ctx, start, end)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "color times retrieved successfully", colorTimes)

}

func (h *ColorTimeHandler) GetColorTime(c *gin.Context) {

	id := c.Query("id")
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

	colorTime, err := h.ColorTimeService.GetColorTime(ctx, id)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "color time retrieved successfully", colorTime)

}

func (h *ColorTimeHandler) UpdateColorTime(c *gin.Context) {

	var req UpdateColorTimeRequest

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

	err := h.ColorTimeService.UpdateColorTime(ctx, &req, id)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
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

	var req DeleteColorTimeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.ColorTimeService.DeleteColorTime(ctx, &req, id)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "color time deleted successfully", nil)
	
}

func (h *ColorTimeHandler) GetTemplateColorTimes(c *gin.Context) {

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	templateColorTimes, err := h.ColorTimeService.GetTemplateColorTimes(ctx)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "template color times retrieved successfully", templateColorTimes)

}

func (h *ColorTimeHandler) GetTemplateColorTime(c *gin.Context) {

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

	templateColorTime, err := h.ColorTimeService.GetTemplateColorTime(ctx, id)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "template color time retrieved successfully", templateColorTime)

}

func (h *ColorTimeHandler) UpdateTemplateColorTime(c *gin.Context) {

	var req UpdateTemplateColorTimeRequest

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

	err := h.ColorTimeService.UpdateTemplateColorTime(ctx, &req, id)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "template color time updated successfully", nil)

}

func (h *ColorTimeHandler) DeleteTemplateColorTime(c *gin.Context) {

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

	err := h.ColorTimeService.DeleteTemplateColorTime(ctx, id)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "template color time deleted successfully", nil)

}

func (h *ColorTimeHandler) AddSlotsToTemplateColorTime(c *gin.Context) {

	var req AddSlotsToTemplateColorTimeRequest

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

	err := h.ColorTimeService.AddSlotsToTemplateColorTime(ctx, &req, id)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "slots added to template color time successfully", nil)

}

func (h *ColorTimeHandler) EditSlotsToTemplateColorTime(c *gin.Context) {

	var req EditSlotsToTemplateColorTimeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request parameters"), nil)
		return
	}

	slot_id := c.Param("slot_id")
	if slot_id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request parameters"), nil)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.ColorTimeService.EditSlotsToTemplateColorTime(ctx, &req, id, slot_id)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "slots edited to template color time successfully", nil)

}

func (h *ColorTimeHandler) ApplyTemplateColorTime(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("invalid request parameters"), nil)
		return
	}

	var req ApplyTemplateColorTimeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
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

	err := h.ColorTimeService.ApplyTemplateColorTime(ctx, &req, id, userID.(string))

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "template color time applied successfully", nil)

}

func (h *ColorTimeHandler) AddTopicToColorTimeWeek(c *gin.Context) {

	var req AddTopicToColorTimeWeekRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
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

	err := h.ColorTimeService.AddTopicToColorTimeWeek(ctx, &req, userID.(string))

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "topic added to color time week successfully", nil)

}

func (h *ColorTimeHandler) GetTopicToColorTimeWeek(c *gin.Context) {

	userID := c.Query("user_id")
	if userID == "" {
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

	data, err := h.ColorTimeService.GetTopicToColorTimeWeek(ctx, userID, orgID, start, end)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "topic to color time week fetched successfully", data)

}

func (h *ColorTimeHandler) DeleteTopicToColorTimeWeek(c *gin.Context) {

	var req DeleteTopicToColorTimeWeekRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
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

	err := h.ColorTimeService.DeleteTopicToColorTimeWeek(ctx, &req, userID.(string))

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "topic deleted from color time week successfully", nil)

}