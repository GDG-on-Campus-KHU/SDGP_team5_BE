// language/gemini.go

package language

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

func TranslateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TranslationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		err := fmt.Errorf("ERROR: 'GEMINI_API_KEY' is NOT set.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	translated, err := callGeminiAPI(apiKey, req.Text, req.TargetLang)
	if err != nil {
		http.Error(w, "Failed to translate: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := TranslationResponse{TranslatedText: translated}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func removePronunciation(text string) string {
	re := regexp.MustCompile(`\s?\(.*\)`)
	return re.ReplaceAllString(text, "")
}

func callGeminiAPI(apiKey, text, targetLang string) (string, error) {
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + apiKey

	// prompt := fmt.Sprintf("Translate the following text to %s:\n\n%s", targetLang, text)
	prompt := fmt.Sprintf("Please translate the following medical information to %s, preserving its technical accuracy:\n\n%s", targetLang, text)

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Gemini API error: %s", string(bodyBytes))
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no translation returned")
	}

	// return result.Candidates[0].Content.Parts[0].Text, nil
	translatedText := result.Candidates[0].Content.Parts[0].Text
	return removePronunciation(translatedText), nil
}
