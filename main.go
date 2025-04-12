// main.go

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"

	_ "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/language"
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
func handler(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("Server is running at http://%s:%s", host, port)
	fmt.Fprintf(w, "%s\n", message)
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

	// load Gemini API Key
	if os.Getenv("GEMINI_API_KEY") == "" {
		log.Println("ERROR: 'GEMINI_API_KEY' is NOT set.")
	}

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	// routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("/translate", language.TranslateHandler)

	handlerWithCORS := c.Handler(http.DefaultServeMux)

	serverAddress := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Server is running at http://%s\n", serverAddress)
	fmt.Printf("Swagger UI available at http://%s/swagger\n", serverAddress)

	if err := http.ListenAndServe(":"+port, handlerWithCORS); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
