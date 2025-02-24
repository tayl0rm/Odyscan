package scanner

import (
	"context"
	"fmt"
	"strings"

	"odyscan/config"

	artifact "cloud.google.com/go/artifactregistry/apiv1"
	artifactpb "cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"
	"github.com/google/go-containerregistry/pkg/crane"
)

// PullImageFromArtifactRegistry fetches an image from Artifact Registry based on user input.
func PullImageFromArtifactRegistry(cfg *config.Config) error {
	ctx := context.Background()
	client, err := artifact.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("❌ failed to create Artifact Registry client: %v", err)
	}
	defer client.Close()

	// Construct the expected image URI correctly based on whether it's a digest or a tag
	var expectedURI string
	if strings.Contains(cfg.ImageName, "@sha256:") {
		// Image is specified by digest
		expectedURI = fmt.Sprintf("europe-west1-docker.pkg.dev/%s/%s/%s", cfg.ProjectID, cfg.RepoName, cfg.ImageName)
	} else {
		// Assume it's a tag
		expectedURI = fmt.Sprintf("europe-west1-docker.pkg.dev/%s/%s/%s:%s", cfg.ProjectID, cfg.RepoName, cfg.ImageName, cfg.Tag)
	}

	fmt.Println("🔍 Checking for image in Artifact Registry:", expectedURI)

	// Verify image existence
	req := &artifactpb.ListDockerImagesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/europe-west1/repositories/%s", cfg.ProjectID, cfg.RepoName),
	}

	it := client.ListDockerImages(ctx, req)
	imageFound := false

	for {
		resp, err := it.Next()
		if err != nil {
			break // No more images in the list
		}
		if strings.Contains(resp.Uri, expectedURI) {
			imageFound = true
			break
		}
	}

	if !imageFound {
		return fmt.Errorf("❌ image %s not found in Artifact Registry", expectedURI)
	}

	fmt.Printf("✅ Found image: %s\n", expectedURI)

	// Define local tar path
	localTarPath := fmt.Sprintf("/tmp/%s.tar", cfg.ImageName)
	cfg.LocalTar = localTarPath

	fmt.Println("🔄 Attempting to pull image using go-containerregistry (crane)...")

	// Pull and save the image
	img, err := crane.Pull(expectedURI)
	if err != nil {
		return fmt.Errorf("⚠️ failed to pull image with crane: %v", err)
	}

	err = crane.Save(img, expectedURI, localTarPath)
	if err != nil {
		return fmt.Errorf("❌ failed to save image using crane: %v", err)
	}

	fmt.Printf("✅ Image pulled and saved to %s\n", localTarPath)
	return nil
}
