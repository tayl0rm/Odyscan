package scanner

import (
	"fmt"
	"os"
	"strings"

	"odyscan/config"

	"github.com/dutchcoders/go-clamd"
)

// ScanWithClamAV scans the extracted files using ClamAV
func ScanWithClamAV(cfg *config.Config) error {
	clam := clamd.NewClamd(fmt.Sprintf("%s:%s", cfg.ClamdHost, cfg.ClamdPort))
	files, err := os.ReadDir(cfg.ExtractDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", cfg.ExtractDir, file.Name())
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file: %v", err)
		}
		defer file.Close()

		res, err := clam.ScanStream(file, make(chan bool))
		if err != nil {
			return fmt.Errorf("ClamAV scan error: %v", err)
		}

		for scanRes := range res {
			fmt.Printf("Scan result for %s: %s\n", filePath, scanRes.Status)
		}

		for scanRes := range res {
			if strings.Contains(scanRes.Status, "FOUND") {
				fmt.Printf("🚨 Malware found in %s: %s\n", filePath, scanRes.Description)
			} else {
				fmt.Printf("✅ File %s is clean.\n", filePath)
			}
		}
	}
	return nil
}
