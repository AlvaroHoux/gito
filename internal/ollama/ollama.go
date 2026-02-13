package ollama

import "net/http"

func isOllamaRunning() bool {
	resp, err := http.Get("http://localhost:11434")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}