package default_colortime

import (
	"colortime-service/helper"
	"colortime-service/pkg/constants"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DefaultColorTimeHandler struct {
	DefaultColorTimeService DefaultColorTimeService
}

func NewDefaultColorTimeHandler(defaultColorTimeService DefaultColorTimeService) *DefaultColorTimeHandler {
	return &DefaultColorTimeHandler{
		DefaultColorTimeService: defaultColorTimeService,
	}
}

func (h *DefaultColorTimeHandler) CreateDefaultDayColorTime(c *gin.Context) {
	var req CreateDefaultDayColorTimeRequest
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

	dayColorTime, err := h.DefaultColorTimeService.CreateDefaultDayColorTime(ctx, &req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "default day color time created successfully", dayColorTime)
}

func (h *DefaultColorTimeHandler) GetDefaultDayColorTime(c *gin.Context) {
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

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, errors.New("user ID not found in context"), nil)
		return
	}

	languageIDStr := c.Query("language_id")
	var languageID *int
	if languageIDStr != "" {
		if parsedLanguageID, err := strconv.Atoi(languageIDStr); err == nil {
			languageID = &parsedLanguageID
		}
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	dayColorTime, err := h.DefaultColorTimeService.GetDefaultDayColorTime(ctx, orgID, date, userID.(string), languageID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "default day color time retrieved successfully", dayColorTime)
}

func (h *DefaultColorTimeHandler) GetDefaultDayColorTimesInRange(c *gin.Context) {
	orgID := c.Query("org_id")
	if orgID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("org_id is required"), nil)
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if startDate == "" || endDate == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("start_date and end_date are required"), nil)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, errors.New("user ID not found in context"), nil)
		return
	}

	languageIDStr := c.Query("language_id")
	var languageID *int
	if languageIDStr != "" {
		if parsedLanguageID, err := strconv.Atoi(languageIDStr); err == nil {
			languageID = &parsedLanguageID
		}
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	dayColorTimes, err := h.DefaultColorTimeService.GetDefaultDayColorTimesInRange(ctx, orgID, startDate, endDate, userID.(string), languageID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "default day color times retrieved successfully", dayColorTimes)
}

func (h *DefaultColorTimeHandler) GetAllDefaultDayColorTimes(c *gin.Context) {
	orgID := c.Query("org_id")
	if orgID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("org_id is required"), nil)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	dayColorTimes, err := h.DefaultColorTimeService.GetAllDefaultDayColorTimes(ctx, orgID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "all default day color times retrieved successfully", dayColorTimes)
}

func (h *DefaultColorTimeHandler) DeleteDefaultDayColorTime(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), nil)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.DefaultColorTimeService.DeleteDefaultDayColorTime(ctx, id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "default day color time deleted successfully", nil)
}

func (h *DefaultColorTimeHandler) GetBlockBySlotID(c *gin.Context) {
	dayID := c.Param("id")
	if dayID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("day id is required"), nil)
		return
	}

	slotID := c.Param("slot_id")
	if slotID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("slot_id is required"), nil)
		return
	}

	block, err := h.DefaultColorTimeService.GetBlockBySlotID(c, dayID, slotID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "block retrieved successfully", block)
}

func (h *DefaultColorTimeHandler) UpdateDefaultColorSlot(c *gin.Context) {
	dayID := c.Param("id")
	if dayID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("day id is required"), nil)
		return
	}

	slotID := c.Param("slot_id")
	if slotID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("slot_id is required"), nil)
		return
	}

	var req UpdateDefaultColorSlotRequest
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

	err := h.DefaultColorTimeService.UpdateDefaultColorSlot(ctx, dayID, slotID, &req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "slot updated successfully", nil)
}
