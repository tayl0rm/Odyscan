package main

import (
	"fmt"
	"log"
	"net/http"

	"odyscan/config"
	"odyscan/scanner"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "odyscan/templates/index.html")
	})

	http.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		r.ParseForm()
		imageName := r.FormValue("image")
		if imageName == "" {
			http.Error(w, "Image name is required", http.StatusBadRequest)
			return
		}

		cfg.ImageName = imageName // Update config with user input
		err := scanner.PullImageFromArtifactRegistry(cfg)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to pull image: %v", err), http.StatusInternalServerError)
			return
		}

		err = scanner.ScanWithClamAV(cfg)
		if err != nil {
			http.Error(w, fmt.Sprintf("Scan failed: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "✅ Scan completed successfully!")
	})

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
