package default_colortime

import (
	"colortime-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, defaultColorTimeHandler *DefaultColorTimeHandler) {
	defaultColorTime := r.Group("api/v1/default-colortime").Use(middleware.Secured())
	{
		// CRUD for default colortime week
		// defaultColorTime.POST("", defaultColorTimeHandler.CreateDefaultColorTimeWeek)
		defaultColorTime.GET("/week", defaultColorTimeHandler.GetDefaultColorTimeWeek)
		defaultColorTime.GET("/weeks", defaultColorTimeHandler.GetAllDefaultColorTimeWeeks)
		defaultColorTime.PUT("/:id", defaultColorTimeHandler.UpdateDefaultColorTimeWeek)
		defaultColorTime.DELETE("/:id", defaultColorTimeHandler.DeleteDefaultColorTimeWeek)

		// Color block and slot management
		defaultColorTime.POST("/week/:week_colortime_id/block/save", defaultColorTimeHandler.CreateDefaultColorBlockAndSaveSlotHandler)
		defaultColorTime.POST("/week/:week_colortime_id/block/add-session", defaultColorTimeHandler.CreateDefaultColorBlockForSessionHandler)

		// Slot update and delete
		defaultColorTime.PUT("/week/:week_id/slot/:slot_id", defaultColorTimeHandler.UpdateDefaultColorSlot)
		defaultColorTime.DELETE("/week/:week_id/slot/:slot_id", defaultColorTimeHandler.DeleteDefaultColorSlot)
	}
}
