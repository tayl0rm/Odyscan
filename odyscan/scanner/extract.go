package scanner

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"odyscan/config"
)

// ExtractImage extracts the filesystem from the saved Docker image tar
func ExtractImage(cfg *config.Config) error {
	fmt.Printf("📂 Extracting image from: %s\n", cfg.LocalTar)

	// Open the tar file
	file, err := os.Open(cfg.LocalTar)
	if err != nil {
		return fmt.Errorf("failed to open tar file: %v", err)
	}
	defer file.Close()

	tr := tar.NewReader(file)

	// Ensure extraction directory exists
	if err := os.MkdirAll(cfg.ExtractDir, 0755); err != nil {
		return fmt.Errorf("failed to create extract directory: %v", err)
	}

	// Extract each file in the tar archive
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar file: %v", err)
		}

		targetPath := filepath.Join(cfg.ExtractDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directories
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
		case tar.TypeReg:
			// Create and write files
			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create file: %v", err)
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, tr)
			if err != nil {
				return fmt.Errorf("failed to extract file: %v", err)
			}
		default:
			fmt.Printf("⚠️ Skipping unknown file type: %s\n", header.Name)
		}
	}

	fmt.Println("✅ Image extracted successfully!")
	return nil
}
