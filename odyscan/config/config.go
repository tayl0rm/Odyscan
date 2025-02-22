package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	ProjectID  string
	RepoName   string
	ImageName  string
	Tag        string
	LocalTar   string
	ExtractDir string
	ClamdHost  string
	ClamdPort  string
}

// LoadConfig loads configuration from a file (local) or environment variables (Kubernetes)
func LoadConfig(configPath string) (*Config, error) {
	// Check if running inside Kubernetes (checks for a service account file)
	_, inCluster := os.LookupEnv("KUBERNETES_SERVICE_HOST")

	if configPath != "" && !inCluster {
		fmt.Println("📄 Loading configuration from YAML file...")
		viper.SetConfigFile(configPath)

		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
	} else {
		fmt.Println("🌐 Running inside Kubernetes, using environment variables...")
	}

	// Read from environment variables
	viper.AutomaticEnv()

	// Load values into struct
	config := &Config{
		ProjectID:  viper.GetString("GCP_PROJECT"),
		RepoName:   viper.GetString("GCP_ARTIFACT_REPO"),
		ImageName:  viper.GetString("GCP_IMAGE_NAME"),
		Tag:        viper.GetString("GCP_IMAGE_TAG"),
		LocalTar:   viper.GetString("LOCAL_TAR"),
		ExtractDir: viper.GetString("EXTRACT_DIR"),
		ClamdHost:  viper.GetString("CLAMD_HOST"),
		ClamdPort:  viper.GetString("CLAMD_PORT"),
	}

	return config, nil
}
