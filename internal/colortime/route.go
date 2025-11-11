package colortime

import (
	"colortime-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, colorTimeHandler *ColorTimeHandler) {
	colorTime := r.Group("api/v1/colortime").Use(middleware.Secured())
	{

		colorTime.GET("/week", colorTimeHandler.GetToColorTimeWeek)
		colorTime.POST("/add-topic/week/:id", colorTimeHandler.AddTopicToColorTimeWeek)
		colorTime.DELETE("/delete-topic/week/:id", colorTimeHandler.DeleteTopicToColorTimeWeek)
		colorTime.POST("/add-topic/day/:id", colorTimeHandler.AddTopicToColorTimeDay)
		colorTime.DELETE("/delete-topic/day/:id", colorTimeHandler.DeleteTopicToColorTimeDay)

		colorTime.PUT("/week/:week_colortime_id/slot/:slot_id", colorTimeHandler.UpdateColorSlotHandler)
	}
}
