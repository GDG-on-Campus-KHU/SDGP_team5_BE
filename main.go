// main.go

package main

import (
	"fmt"
	"net/http"
)

var (
	host = "localhost"
	port = "5000"
)

func handler(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("Server is running at http://%s:%s", host, port)
	fmt.Fprintf(w, "%s\n", message)
}

func main() {
	http.HandleFunc("/", handler)

	serverAddress := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Server is running at http://%s\n", serverAddress)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
