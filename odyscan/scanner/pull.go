package scanner

import (
	"context"
	"fmt"
	"strings"

	"odyscan/config" // Import config package

	artifact "cloud.google.com/go/artifactregistry/apiv1"
	artifactpb "cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"
	"github.com/google/go-containerregistry/pkg/crane"
)

// PullImageFromArtifactRegistry fetches an image from Artifact Registry based on user input.
func PullImageFromArtifactRegistry(cfg *config.Config) error {
	ctx := context.Background()
	client, err := artifact.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create artifact registry client: %v", err)
	}
	defer client.Close()

	req := &artifactpb.ListDockerImagesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/europe-west1/repositories/%s", cfg.ProjectID, cfg.RepoName),
	}

	it := client.ListDockerImages(ctx, req)
	for {
		resp, err := it.Next()
		if err != nil {
			break
		}
		expectedURI := fmt.Sprintf("europe-west1-docker.pkg.dev/%s/%s/%s", cfg.ProjectID, cfg.RepoName, cfg.ImageName)
		if strings.Contains(resp.Uri, expectedURI) {
			fmt.Printf("✅ Found image: %s\n", resp.Uri)
			localTarPath := fmt.Sprintf("/tmp/%s.tar", cfg.ImageName)
			cfg.LocalTar = localTarPath
			img, err := crane.Pull(expectedURI)
			if err != nil {
				return fmt.Errorf("failed to pull image: %v", err)
			}
			if err := crane.Save(img, expectedURI, localTarPath); err != nil {
				return fmt.Errorf("failed to pull and save image: %v", err)
			}
			return nil
		}
	}
	return fmt.Errorf("❌ image %s not found in Artifact Registry", cfg.ImageName)
}
