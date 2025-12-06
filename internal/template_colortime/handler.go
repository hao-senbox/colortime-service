package templatecolortime

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

type TemplateColorTimeHandler struct {
	TemplateColorTimeService TemplateColorTimeService
}

func NewTemplateColorTimeHandler(templateColorTimeService TemplateColorTimeService) *TemplateColorTimeHandler {
	return &TemplateColorTimeHandler{
		TemplateColorTimeService: templateColorTimeService,
	}
}

func (h *TemplateColorTimeHandler) CreateTemplateColorTime(c *gin.Context) {
	var request CreateTemplateColorTimeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	templateColorTime, err := h.TemplateColorTimeService.CreateTemplateColorTime(ctx, request, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "template color time created successfully", templateColorTime)
}

func (h *TemplateColorTimeHandler) GetTemplateColorTime(c *gin.Context) {
	orgID := c.Query("org_id")
	if orgID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("org_id is required"), nil)
		return
	}

	termID := c.Query("term_id")
	if termID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("term_id is required"), nil)
		return
	}

	date := c.Query("date")
	if date == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("date is required"), nil)
		return
	}

	languageIDStr := c.Query("language_id")
	var languageID *int
	if languageIDStr != "" {
		if parsedLanguageID, err := strconv.Atoi(languageIDStr); err == nil {
			languageID = &parsedLanguageID
		}
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

	templateColorTime, err := h.TemplateColorTimeService.GetTemplateColorTime(ctx, orgID, termID, date, languageID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "template color time retrieved successfully", templateColorTime)
}

func (h *TemplateColorTimeHandler) DuplicateTemplateColorTime(c *gin.Context) {
	var request DuplicateTemplateColorTimeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	err := h.TemplateColorTimeService.DuplicateTemplateColorTime(ctx, request, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "template color time duplicated successfully", nil)
}

func (h *TemplateColorTimeHandler) ApplyTemplateColorTime(c *gin.Context) {
	var request ApplyTemplateColorTimeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	err := h.TemplateColorTimeService.ApplyTemplateColorTime(ctx, request, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "template color time applied successfully", nil)
}

func (h *TemplateColorTimeHandler) UpdateTemplateColorTimeSlot(c *gin.Context) {
	templateColorTimeID := c.Param("id")
	if templateColorTimeID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("template color time id is required"), nil)
		return
	}

	slotID := c.Param("slot_id")
	if slotID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("slot id is required"), nil)
		return
	}

	var request UpdateTemplateColorTimeSlotRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	err := h.TemplateColorTimeService.UpdateTemplateColorTimeSlot(ctx, templateColorTimeID, slotID, &request, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "template color time slot updated successfully", nil)
}

func (h *TemplateColorTimeHandler) DeleteTemplateColorTimeBlock(c *gin.Context) {
	templateColorTimeID := c.Param("id")
	if templateColorTimeID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("template color time id is required"), nil)
		return
	}

	blockID := c.Param("block_id")
	if blockID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("block id is required"), nil)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, errors.New("user ID not found in context"), nil)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.TemplateColorTimeService.DeleteTemplateColorTimeBlock(ctx, templateColorTimeID, blockID, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "template color time block deleted successfully", nil)
}

func (h *TemplateColorTimeHandler) DeleteTemplateColorTimeSlot(c *gin.Context) {
	templateColorTimeID := c.Param("id")
	if templateColorTimeID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("template color time id is required"), nil)
		return
	}

	slotID := c.Param("slot_id")
	if slotID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("slot id is required"), nil)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, errors.New("user ID not found in context"), nil)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.TemplateColorTimeService.DeleteTemplateColorTimeSlot(ctx, templateColorTimeID, slotID, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "template color time slot deleted successfully", nil)
}

func (h *TemplateColorTimeHandler) CopySlotToTemplateColorTime(c *gin.Context) {
	blockID := c.Param("block_id")
	if blockID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("block id is required"), nil)
		return
	}

	var request CopySlotToTemplateColorTimeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, errors.New("user ID not found in context"), nil)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, errors.New("token not found in context"), nil)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.TemplateColorTimeService.CopySlotToTemplateColorTime(ctx, blockID, &request, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, nil)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "block copied to template color time successfully", nil)
}
