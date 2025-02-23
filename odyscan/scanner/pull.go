package scanner

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"odyscan/config" // Import config package

	artifact "cloud.google.com/go/artifactregistry/apiv1"
	artifactpb "cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"
	"github.com/google/go-containerregistry/pkg/crane"
	"google.golang.org/api/iterator"
)

// PullImageFromArtifactRegistry pulls an image and extracts it
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
		// expectedURI := fmt.Sprintf("europe-west1-docker.pkg.dev/%s/%s/%s:%s", cfg.ProjectID, cfg.RepoName, cfg.ImageName, cfg.Tag)
		expectedURI := fmt.Sprintf("europe-west1-docker.pkg.dev/%s/%s/%s", cfg.ProjectID, cfg.RepoName, cfg.ImageName)
		if resp.Uri == expectedURI {
			fmt.Printf("✅ Found image: %s (Size: %d bytes)\n", resp.Uri, resp.ImageSizeBytes)

			// Define output file for tar
			localTarPath := filepath.Join("/tmp", fmt.Sprintf("%s_%s.tar", cfg.ImageName, cfg.Tag))
			cfg.LocalTar = localTarPath // Store tar path in config

			// Pull and save image as tar
			if err := pullAndSaveImage(cfg, expectedURI); err != nil {
				return fmt.Errorf("failed to pull and save image: %v", err)
			}

			// Extract the image after pulling
			fmt.Println("📦 Extracting image layers...")
			if err := ExtractImage(cfg); err != nil {
				return fmt.Errorf("failed to extract image: %v", err)
			}

			fmt.Println("✅ Image pulled and extracted successfully!")
			return nil
		}
	}

	return fmt.Errorf("❌ image %s:%s not found in Artifact Registry", cfg.ImageName, cfg.Tag)
}

func pullAndSaveImage(cfg *config.Config, imageURI string) error {
	fmt.Printf("🔄 Pulling image: %s\n", imageURI)

	// Define output tar file path
	tarPath := filepath.Join("/tmp", fmt.Sprintf("%s.tar", cfg.ImageName))
	cfg.LocalTar = tarPath // Store tar path in config

	// Attempt to pull using `crane` (Containerd/K3s)
	fmt.Println("🔄 Attempting to pull image using go-containerregistry (crane)...")
	img, err := crane.Pull(imageURI)
	if err != nil {
		fmt.Printf("⚠️  Crane pull failed: %v. Falling back to Docker...\n", err)
		return pullAndSaveImageWithDocker(imageURI, tarPath)
	}

	// Save the pulled image as a tar file
	err = crane.Save(img, imageURI, tarPath)
	if err != nil {
		return fmt.Errorf("failed to save image using crane: %v", err)
	}

	fmt.Printf("✅ Image pulled and saved to %s using containerd\n", tarPath)
	return nil
}

// pullAndSaveImageWithDocker pulls an image using Docker (for local testing)
func pullAndSaveImageWithDocker(imageURI, tarPath string) error {
	fmt.Println("🔄 Pulling image using Docker...")

	// Pull the image with Docker
	err := exec.Command("docker", "pull", imageURI).Run()
	if err != nil {
		return fmt.Errorf("failed to pull image with Docker: %v", err)
	}

	// Save the image as a tar file
	err = exec.Command("docker", "save", "-o", tarPath, imageURI).Run()
	if err != nil {
		return fmt.Errorf("failed to save image as tar: %v", err)
	}

	fmt.Printf("✅ Image saved to %s using Docker\n", tarPath)
	return nil
}
