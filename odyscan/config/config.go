package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config struct for app settings
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

// LoadConfig loads configuration from YAML or environment variables
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	return &Config{
		ProjectID:  viper.GetString("GCP_PROJECT"),
		RepoName:   viper.GetString("GCP_ARTIFACT_REPO"),
		ImageName:  viper.GetString("GCP_IMAGE_NAME"),
		Tag:        viper.GetString("GCP_IMAGE_TAG"),
		LocalTar:   viper.GetString("LOCAL_TAR"),
		ExtractDir: viper.GetString("EXTRACT_DIR"),
		ClamdHost:  viper.GetString("CLAMD_HOST"),
		ClamdPort:  viper.GetString("CLAMD_PORT"),
	}, nil
}
