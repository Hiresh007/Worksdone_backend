package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

func DataRoutes(router *gin.Engine) {
	router.POST("resume/scan", controllers.ResumeScan())
	router.POST("resume/summarize", controllers.Summarize())
	router.POST("resume/score", controllers.ResumeScore())
}
