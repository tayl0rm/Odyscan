package main

import (
	"fmt"
	"log"
	"odyscan/config"
	"odyscan/scanner"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Pull image and extract layers
	err = scanner.PullImageFromArtifactRegistry(cfg)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Scan extracted files with ClamAV
	err = scanner.ScanWithClamAV(cfg)
	if err != nil {
		log.Fatalf("ClamAV scan failed: %v", err)
	} else {
		fmt.Println("✅ Malware scan completed successfully!")
	}
}
