package default_colortime

import (
	"colortime-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, defaultColorTimeHandler *DefaultColorTimeHandler) {
	defaultColorTime := r.Group("api/v1/default-colortime").Use(middleware.Secured())
	{
		defaultColorTime.POST("/day", defaultColorTimeHandler.CreateDefaultDayColorTime)
		defaultColorTime.GET("/day", defaultColorTimeHandler.GetDefaultDayColorTime)
		defaultColorTime.GET("/days", defaultColorTimeHandler.GetDefaultDayColorTimesInRange)
		defaultColorTime.GET("/all-days", defaultColorTimeHandler.GetAllDefaultDayColorTimes)
		defaultColorTime.DELETE("/day/:id", defaultColorTimeHandler.DeleteDefaultDayColorTime)

		defaultColorTime.GET("/day/:id/block/:slot_id", defaultColorTimeHandler.GetBlockBySlotID)
		defaultColorTime.PUT("/day/:id/slot/edit/:slot_id", defaultColorTimeHandler.UpdateDefaultColorSlot)
	}
}
