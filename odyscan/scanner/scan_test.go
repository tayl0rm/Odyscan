package scanner

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"odyscan/config"
)

// Mock ClamAV Server
func startMockClamAVServer(t *testing.T) *httptest.Server {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate ClamAV response
		if strings.Contains(r.URL.Path, "stream") {
			// Simulate a clean scan result
			fmt.Fprintln(w, "stream: OK")
		} else {
			http.Error(w, "Invalid request", http.StatusBadRequest)
		}
	}))
	t.Cleanup(mockServer.Close)
	return mockServer
}

// createTestFiles creates dummy files to be scanned
func createTestFiles(scanDir string) error {
	err := os.MkdirAll(scanDir, 0755)
	if err != nil {
		return err
	}
	testFile := filepath.Join(scanDir, "testfile.txt")
	content := []byte("This is a harmless test file.")
	return os.WriteFile(testFile, content, 0644)
}

// TestScanWithClamAV tests scanning with a mocked ClamAV server
func TestScanWithClamAV(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test_scan")

	// Create test files
	err := createTestFiles(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create test files: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Start mock ClamAV server
	mockServer := startMockClamAVServer(t)
	mockHost, mockPort, _ := net.SplitHostPort(strings.TrimPrefix(mockServer.URL, "http://"))

	cfg := &config.Config{
		ExtractDir:     tmpDir,
		ClamdNamespace: mockHost,
		ClamdPort:      mockPort,
	}

	// Run scan with mock server
	err = ScanWithClamAV(cfg)
	if err != nil {
		t.Fatalf("ScanWithClamAV failed: %v", err)
	}
}
