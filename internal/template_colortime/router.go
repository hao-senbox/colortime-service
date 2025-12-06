package templatecolortime

import (
	"colortime-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, templateColorTimeHandler *TemplateColorTimeHandler) {
	templateColorTime := r.Group("api/v1/template-colortime").Use(middleware.Secured())
	{
		templateColorTime.GET("", templateColorTimeHandler.GetTemplateColorTime)
		templateColorTime.POST("", templateColorTimeHandler.CreateTemplateColorTime)
		templateColorTime.POST("/duplicate", templateColorTimeHandler.DuplicateTemplateColorTime)
		templateColorTime.POST("/apply-template", templateColorTimeHandler.ApplyTemplateColorTime)
		templateColorTime.PUT("/copy-slot/:block_id", templateColorTimeHandler.CopySlotToTemplateColorTime)
		templateColorTime.PUT("/:id/update-slot/:slot_id", templateColorTimeHandler.UpdateTemplateColorTimeSlot)
		templateColorTime.DELETE("/:id/delete-block/:block_id", templateColorTimeHandler.DeleteTemplateColorTimeBlock)
		templateColorTime.DELETE("/:id/delete-slot/:slot_id", templateColorTimeHandler.DeleteTemplateColorTimeSlot)
	}
}
