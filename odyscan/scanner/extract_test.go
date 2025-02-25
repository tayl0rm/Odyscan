package scanner

import (
	"archive/tar"
	"os"
	"path/filepath"
	"testing"

	"odyscan/config"
)

// createTestTar creates a temporary tar archive for testing extraction
func createTestTar(tarPath string) error {
	file, err := os.Create(tarPath)
	if err != nil {
		return err
	}
	defer file.Close()

	tw := tar.NewWriter(file)
	defer tw.Close()

	// Create a test file inside the tar
	files := []struct {
		Name, Body string
	}{
		{"testfile1.txt", "This is a test file."},
		{"testfile2.txt", "Another test file."},
	}

	for _, file := range files {
		header := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			return err
		}
	}
	return nil
}

func TestExtractImage(t *testing.T) {
	tmpDir := os.TempDir()
	tmpTar := filepath.Join(tmpDir, "test.tar")
	extractDir := filepath.Join(tmpDir, "test_extract")

	// Create test tar file
	err := createTestTar(tmpTar)
	if err != nil {
		t.Fatalf("Failed to create test tar file: %v", err)
	}

	cfg := &config.Config{
		LocalTar:   tmpTar,
		ExtractDir: extractDir,
	}

	// Run extraction
	err = ExtractImage(cfg)
	if err != nil {
		t.Fatalf("ExtractImage failed: %v", err)
	}

	// Check if extracted files exist
	expectedFiles := []string{"testfile1.txt", "testfile2.txt"}
	for _, file := range expectedFiles {
		if _, err := os.Stat(filepath.Join(extractDir, file)); os.IsNotExist(err) {
			t.Errorf("Expected file %s not found after extraction", file)
		}
	}

	// Cleanup
	os.RemoveAll(tmpTar)
	os.RemoveAll(extractDir)
}
