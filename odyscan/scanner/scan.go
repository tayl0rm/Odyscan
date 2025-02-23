package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"odyscan/config"

	"github.com/dutchcoders/go-clamd"
)

// ScanWithClamAV scans the extracted files using ClamAV running in the same K3s cluster
func ScanWithClamAV(cfg *config.Config) error {
	clamdSocket := fmt.Sprintf("%s:%s", cfg.ClamdHost, cfg.ClamdPort)
	clam := clamd.NewClamd(clamdSocket)

	// Get list of extracted files
	files, err := getFilesInDir(cfg.ExtractDir)
	if err != nil {
		return fmt.Errorf("failed to read extracted files: %v", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files found in extracted directory: %s", cfg.ExtractDir)
	}

	fmt.Printf("🔍 Scanning %d files with ClamAV...\n", len(files))

	for _, filePath := range files {
		err := scanFileWithClamAV(clam, filePath)
		if err != nil {
			fmt.Printf("⚠️ Error scanning %s: %v\n", filePath, err)
		}
	}

	fmt.Println("✅ ClamAV scanning completed!")
	return nil
}

// scanFileWithClamAV sends a file to ClamAV for scanning
func scanFileWithClamAV(clam *clamd.Clamd, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Perform scan
	res, err := clam.ScanStream(file, make(chan bool))
	if err != nil {
		return fmt.Errorf("ClamAV scan error: %v", err)
	}

	// Read results
	for scanRes := range res {
		if strings.Contains(scanRes.Status, "FOUND") {
			fmt.Printf("🚨 Malware found in %s: %s\n", filePath, scanRes.Description)
		} else {
			fmt.Printf("✅ File %s is clean.\n", filePath)
		}
	}

	return nil
}

// getFilesInDir returns a list of all file paths in a directory
func getFilesInDir(dir string) ([]string, error) {
	var filePaths []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filePaths = append(filePaths, path)
		}
		return nil
	})

	return filePaths, err
}
