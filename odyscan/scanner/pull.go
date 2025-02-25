package scanner

import (
	"context"
	"fmt"
	"odyscan/config"
	"path/filepath"

	artifactregistry "cloud.google.com/go/artifactregistry/apiv1"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/v1/google"
)

// PullImageFromArtifactRegistry pulls an image using crane with authentication
func PullImageFromArtifactRegistry(cfg *config.Config) error {
	ctx := context.Background()
	client, err := artifactregistry.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create artifact registry client: %v", err)
	}
	defer client.Close()

	imageURI := fmt.Sprintf("europe-west1-docker.pkg.dev/%s/%s/%s", cfg.ProjectID, cfg.RepoName, cfg.ImageName)
	localTarPath := filepath.Join("/tmp", fmt.Sprintf("%s.tar", cfg.ImageName))
	cfg.LocalTar = localTarPath

	// Ensure authentication
	keychain := google.Keychain
	fmt.Println("🔄 Authenticating with Artifact Registry...")
	craneOpts := crane.WithAuthFromKeychain(keychain)

	// Pull the image
	fmt.Printf("🔄 Pulling image: %s\n", imageURI)
	img, err := crane.Pull(imageURI, craneOpts)
	if err != nil {
		return fmt.Errorf("failed to pull image: %v", err)
	}

	// Save the pulled image as a tar file
	if err := crane.Save(img, imageURI, localTarPath); err != nil {
		return fmt.Errorf("failed to save image: %v", err)
	}

	fmt.Printf("✅ Image pulled and saved to %s\n", localTarPath)
	return nil
}
