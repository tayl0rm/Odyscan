package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"odyscan/config"
	"odyscan/scanner"
	"path/filepath"
)

func main() {
	http.HandleFunc("/", serveIndex)     // Serve index.html
	http.HandleFunc("/scan", handleScan) // Handle image scan requests

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

// Handle scan request
func handleScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	imageName := r.FormValue("imageName") // Get image name from user input
	if imageName == "" {
		http.Error(w, "Image name is required", http.StatusBadRequest)
		return
	}

	// Load config
	cfg, err := config.LoadConfig("/app/config/config.yaml")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading config: %v", err), http.StatusInternalServerError)
		return
	}
	cfg.ImageName = imageName
	cfg.LocalTar = fmt.Sprintf("/tmp/%s.tar", imageName)
	cfg.ExtractDir = fmt.Sprintf("/tmp/%s_extracted", imageName)

	// Pull image
	err = scanner.PullImageFromArtifactRegistry(cfg)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error pulling image: %v", err), http.StatusInternalServerError)
		return
	}

	// Extract image
	err = scanner.ExtractImage(cfg)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error extracting image: %v", err), http.StatusInternalServerError)
		return
	}

	// Scan extracted files with ClamAV
	err = scanner.ScanWithClamAV(cfg)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error scanning image: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "✅ Image pulled, extracted, and scanned successfully!")
}
