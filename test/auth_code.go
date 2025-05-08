// test/auth_code.go
// issue temporary auth code
// go run test/auth_code.go

package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "os/exec"
    "runtime"

    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    RunAuthCodeTest()
}


func RunAuthCodeTest() {
    clientID := os.Getenv("GOOGLE_AUTH_CLIENT_ID")
    if clientID == "" {
        log.Fatal("GOOGLE_AUTH_CLIENT_ID env variable not set")
    }
	redirectURI := os.Getenv("GOOGLE_AUTH_REDIRECT_URL")
    scope := "openid email profile"

    authURL := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&access_type=offline&prompt=consent",
        clientID,
        redirectURI,
        scope,
    )

    fmt.Println("✅️ Opening browser for Google login...")
    openBrowser(authURL)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        code := r.URL.Query().Get("code")
        if code == "" {
            http.Error(w, "No code found", http.StatusBadRequest)
            return
        }
        fmt.Printf("\n✅ Authorization Code: %s\n", code)
        w.Write([]byte("<h1>You can close this tab.</h1>"))
        os.Exit(0)
    })
    fmt.Println("✅️ Waiting for redirect with auth code...")
    log.Fatal(http.ListenAndServe(":5100", nil))
}


func openBrowser(url string) {
    var err error
    switch runtime.GOOS {
    case "linux":
        err = exec.Command("xdg-open", url).Start()
    case "windows":
        err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
    case "darwin":
        err = exec.Command("open", url).Start()
    default:
        err = fmt.Errorf("unsupported platform")
    }
    if err != nil {
        log.Fatalf("failed to open browser: %v", err)
    }
}