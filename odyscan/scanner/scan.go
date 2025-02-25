package scanner

import (
	"fmt"
	"os"
	"strings"

	"odyscan/config"

	"github.com/dutchcoders/go-clamd"
)

// ScanWithClamAV scans extracted files using ClamAV running in another namespace
func ScanWithClamAV(cfg *config.Config) error {
	// Construct ClamAV service address using namespace
	clamavService := fmt.Sprintf("tcp://clamd.%s.svc.cluster.local:%s", cfg.ClamdNamespace, cfg.ClamdPort)
	clam := clamd.NewClamd(clamavService)

	fmt.Printf("🔍 Connecting to ClamAV at %s\n", clamavService)

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
			if strings.Contains(scanRes.Status, "FOUND") {
				fmt.Printf("🚨 Malware found in %s: %s\n", filePath, scanRes.Description)
			} else {
				fmt.Printf("✅ File %s is clean.\n", filePath)
			}
		}
	}
	return nil
}
