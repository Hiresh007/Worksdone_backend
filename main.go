package main

import (
	"log"
	"net/http"
	"os"
	"server/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(CORSMiddleware())
	routes.AuthRoutes(router)
	routes.DataRoutes(router)
	router.GET("/working", working)
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

type simple struct {
	Hello   string `json:"hello"`
	Message string `json:"message"`
}

func working(c *gin.Context) {
	message := simple{
		Hello:   "World",
		Message: "Api Working",
	}

	c.JSON(http.StatusOK, message)
}
