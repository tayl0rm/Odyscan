package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"odyscan/config"
	"odyscan/scanner"
	"path/filepath"
	"time"
)

func main() {
	http.HandleFunc("/", serveIndex)     // Serve index.html
	http.HandleFunc("/scan", handleScan) // Handle image scan requests with SSE

	fmt.Println("🚀 Server started on port 8080")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("❌ Server failed to start: %v\n", err)
	}
}

// Serve index.html
func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	tmplPath := filepath.Join("templates", "index.html") // Ensure correct path
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Handle scan request with SSE progress updates
func handleScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	imageName := r.FormValue("imageName")
	if imageName == "" {
		http.Error(w, "Image name is required", http.StatusBadRequest)
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Load config
	cfg, err := config.LoadConfig("/app/config/config.yaml")
	if err != nil {
		sendSSEMessage(w, "error", fmt.Sprintf("Error loading config: %v", err))
		return
	}
	cfg.ImageName = imageName
	cfg.LocalTar = fmt.Sprintf("/tmp/%s.tar", imageName)
	cfg.ExtractDir = fmt.Sprintf("/tmp/%s_extracted", imageName)

	// Steps with SSE updates
	steps := []struct {
		Message string
		Action  func() error
	}{
		{"Pulling image from Artifact Registry...", func() error { return scanner.PullImageFromArtifactRegistry(cfg) }},
		{"Extracting image...", func() error { return scanner.ExtractImage(cfg) }},
		{"Scanning with ClamAV...", func() error { return scanner.ScanWithClamAV(cfg) }},
	}

	for i, step := range steps {
		sendSSEMessage(w, "progress", step.Message)
		time.Sleep(2 * time.Second) // Simulate processing delay
		if err := step.Action(); err != nil {
			sendSSEMessage(w, "error", fmt.Sprintf("Error: %v", err))
			return
		}
		sendSSEMessage(w, "progress", fmt.Sprintf("%d%% complete", (i+1)*33))
	}

	sendSSEMessage(w, "complete", "✅ Image pulled, extracted, and scanned successfully!")
}

// sendSSEMessage sends an SSE event
func sendSSEMessage(w http.ResponseWriter, eventType, message string) {
	msg := map[string]string{"type": eventType, "message": message}
	jsonData, _ := json.Marshal(msg)
	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	w.(http.Flusher).Flush()
}
