package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("users/signin", controllers.Signin())
	router.POST("users/signup", controllers.Signup())
}
