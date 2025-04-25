// main.go

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/docs"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/auth"
	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/language"
	situationHandler "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/situation/handler"
)

// @title SDGP-team5-ResQ-BE
// @version 1.0
// @description API documentation
// @host localhost:5100
// @BasePath /

// @Summary check server status
// @Description check if the server is running
// @Tags Status
// @Accept json
// @Produce json
// @Success 200 {string} string "Server is running"
// @Router / [get]
func StatusHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Server is running!",
	})
}

var (
	host = "localhost"
	port = "5100"
)

func main() {

	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("ERROR: .env file NOT FOUND.")
	}

	fmt.Println("GOOGLE_AUTH_CLIENT_ID:", os.Getenv("GOOGLE_AUTH_CLIENT_ID"))
	fmt.Println("GOOGLE_AUTH_CLIENT_SECRET:", os.Getenv("GOOGLE_AUTH_CLIENT_SECRET"))

	// initialize Firebase authentication
	// firebase.InitFirebase()

	// initialize Google OAuth2 configuration
	auth.InitGoogleOAuthConfig()

	// initialize database & GCS
	dbConfig.InitMongo()
	dbConfig.InitGCS()

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "*")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	// health check
	r.GET("/", StatusHandler)

	// Swagger
	r.GET("/swagger/*any", gin.WrapH(httpSwagger.Handler()))

	// routes
	auth.RegisterAuthRoutes(r)
	situationHandler.RegisterSituationRoutes(r)
	r.POST("/translate", gin.WrapF(language.TranslateHandler))

	// start the server
	serverAddress := host + ":" + port
	fmt.Printf("Server is running at http://%s\n", serverAddress)
	fmt.Printf("Swagger UI available at http://%s/swagger\n", serverAddress)

	if err := r.Run(":5100"); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
