// main.go

package main

import (
	"fmt"
	"net/http"

	_ "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title SDGP-team5-BE
// @version 1.0
// @description API documentation
// @host localhost:5000
// @BasePath /

// @Summary check server status
// @Description checks if the server is running
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
	http.HandleFunc("/", handler)

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	serverAddress := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Server is running at http://%s\n", serverAddress)
	fmt.Printf("Swagger UI available at http://%s/swagger\n", serverAddress)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
