package colortime

import (
	"colortime-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, colorTimeHandler *ColorTimeHandler) {
	colorTime := r.Group("api/v1/colortime").Use(middleware.Secured())
	{
		// colorTime.POST("", colorTimeHandler.CreateColorTime)
		// colorTime.GET("", colorTimeHandler.GetColorTimes)
		// colorTime.GET("/:id", colorTimeHandler.GetColorTime)
		// colorTime.PUT("/:id", colorTimeHandler.UpdateColorTime)
		// colorTime.DELETE("/:id", colorTimeHandler.DeleteColorTime)
		colorTime.GET("/week", colorTimeHandler.GetToColorTimeWeek)
		colorTime.POST("/add-topic/week/:id", colorTimeHandler.AddTopicToColorTimeWeek)
		colorTime.POST("/delete-topic/week/:id", colorTimeHandler.DeleteTopicToColorTimeWeek)
		colorTime.POST("/add-topic/day/:id", colorTimeHandler.AddTopicToColorTimeDay)
		colorTime.POST("/delete-topic/day/:id", colorTimeHandler.DeleteTopicToColorTimeDay)

		colorTime.POST("/week/:week_colortime_id/block/save", colorTimeHandler.CreateColorBlockAndSaveSlotHandler)
		colorTime.POST("/week/:week_colortime_id/block/add-session", colorTimeHandler.CreateColorBlockForSessionHandler)
		
		// colorTime.POST("/template", colorTimeHandler.CreateTemplateColorTime)
		// colorTime.GET("/template", colorTimeHandler.GetTemplateColorTimes)
		// colorTime.GET("template/:id", colorTimeHandler.GetTemplateColorTime)
		// colorTime.PUT("template/:id", colorTimeHandler.UpdateTemplateColorTime)
		// colorTime.DELETE("template/:id", colorTimeHandler.DeleteTemplateColorTime)
		// colorTime.POST("template/:id/add-slots", colorTimeHandler.AddSlotsToTemplateColorTime)
		// colorTime.POST("template/:id/edit-slots/:slot_id", colorTimeHandler.EditSlotsToTemplateColorTime)
		// colorTime.POST("template/:id/apply-template", colorTimeHandler.ApplyTemplateColorTime)
	}
}
