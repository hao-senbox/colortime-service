package colortime

import (
	"colortime-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, colorTimeHandler *ColorTimeHandler) {
	colorTime := r.Group("api/v1/colortime").Use(middleware.Secured())
	{
		colorTime.POST("", colorTimeHandler.CreateColorTime)
		colorTime.GET("", colorTimeHandler.GetColorTimes)
		colorTime.GET("/:id", colorTimeHandler.GetColorTime)
		colorTime.PUT("/:id", colorTimeHandler.UpdateColorTime)
		colorTime.DELETE("/:id", colorTimeHandler.DeleteColorTime)
	}
}
