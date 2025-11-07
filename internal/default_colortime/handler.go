package default_colortime

import (
	"colortime-service/helper"
	"colortime-service/pkg/constants"
	"context"
	"errors"
	"fmt"
	"net/http"

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

func (h *DefaultColorTimeHandler) GetDefaultColorTimeWeek(c *gin.Context) {
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

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user ID not found"), nil)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.DefaultColorTimeService.GetDefaultColorTimeWeek(ctx, orgID, start, end, userID.(string))

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "default colortime week fetched successfully", data)
}

func (h *DefaultColorTimeHandler) GetAllDefaultColorTimeWeeks(c *gin.Context) {
	orgID := c.Query("org_id")
	if orgID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("organization id required"), nil)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.DefaultColorTimeService.GetAllDefaultColorTimeWeeks(ctx, orgID)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "all default colortime weeks fetched successfully", data)
}

func (h *DefaultColorTimeHandler) UpdateDefaultColorTimeWeek(c *gin.Context) {
	var req UpdateDefaultColorTimeWeekRequest

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

	err := h.DefaultColorTimeService.UpdateDefaultColorTimeWeek(ctx, id, &req)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "default colortime week updated successfully", nil)
}

func (h *DefaultColorTimeHandler) DeleteDefaultColorTimeWeek(c *gin.Context) {
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

	err := h.DefaultColorTimeService.DeleteDefaultColorTimeWeek(ctx, id)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "default colortime week deleted successfully", nil)
}

func (h *DefaultColorTimeHandler) CreateDefaultColorBlockAndSaveSlotHandler(c *gin.Context) {
	weekColorTimeID := c.Param("week_colortime_id")

	var req CreateDefaultColorBlockWithSlotRequest
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

	colorBlock, err := h.DefaultColorTimeService.CreateDefaultColorBlockAndSaveSlot(c, weekColorTimeID, &req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "default color block with single slot created successfully", colorBlock)
}

func (h *DefaultColorTimeHandler) CreateDefaultColorBlockForSessionHandler(c *gin.Context) {
	weekColorTimeID := c.Param("week_colortime_id")

	var req CreateDefaultColorBlockWithSlotRequest
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

	colorBlock, err := h.DefaultColorTimeService.CreateDefaultColorBlockForSession(c, weekColorTimeID, &req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "default color block for session created successfully", colorBlock)
}

func (h *DefaultColorTimeHandler) UpdateDefaultColorSlot(c *gin.Context) {
	weekID := c.Param("week_id")
	slotID := c.Param("slot_id")

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

	err := h.DefaultColorTimeService.UpdateDefaultColorSlot(ctx, weekID, slotID, &req)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "default color slot updated successfully", nil)
}

func (h *DefaultColorTimeHandler) DeleteDefaultColorSlot(c *gin.Context) {
	weekID := c.Param("week_id")
	slotID := c.Param("slot_id")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.DefaultColorTimeService.DeleteDefaultColorSlot(ctx, weekID, slotID)

	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "default color slot deleted successfully", nil)
}
