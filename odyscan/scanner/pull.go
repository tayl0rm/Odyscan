package scanner

import (
	"context"
	"fmt"

	"odyscan/config" // Import config package

	artifact "cloud.google.com/go/artifactregistry/apiv1"
	artifactpb "cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"
	"google.golang.org/api/iterator"
)

// PullImageFromArtifactRegistry retrieves image metadata using values from config
func PullImageFromArtifactRegistry(cfg *config.Config) error {
	ctx := context.Background()

	// Create Artifact Registry client
	client, err := artifact.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create artifact registry client: %v", err)
	}
	defer client.Close()

	// Construct request
	req := &artifactpb.ListDockerImagesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/europe-west1/repositories/%s", cfg.ProjectID, cfg.RepoName),
	}

	// Iterate through images
	it := client.ListDockerImages(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("error listing images: %v", err)
		}

		// Match requested image and tag
		expectedURI := fmt.Sprintf("europe-west1-docker.pkg.dev/%s/%s/%s:%s", cfg.ProjectID, cfg.RepoName, cfg.ImageName, cfg.Tag)
		if resp.Uri == expectedURI {
			fmt.Printf("✅ Found image: %s (Size: %d bytes)\n", resp.Uri, resp.ImageSizeBytes)
			return nil
		}
	}

	return fmt.Errorf("❌ image %s:%s not found in Artifact Registry", cfg.ImageName, cfg.Tag)
}
