package main

import (
	"fmt"
	"odyscan/config"
	"odyscan/scanner"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		fmt.Printf("❌ Error loading config: %v\n", err)
		return
	}

	fmt.Println("🔹 Pulling image from GCP Artifact Registry...")
	if err := scanner.PullImageFromArtifactRegistry(cfg); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Println("🔹 Extracting image contents...")
	if err := scanner.ExtractImage(cfg); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Println("🔹 Scanning extracted files with ClamAV...")
	if err := scanner.ScanWithClamAV(cfg); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Println("✅ Scan complete!")
}
