package scanner

import (
	"archive/tar"
	"fmt"
	"io"
	"odyscan/config"
	"os"
)

// Refactor to mirror changes made to pull.go & use the ArfitactRegistry client to pull the image
// ExtractImage extracts the filesystem of the saved Docker image
func ExtractImage(cfg *config.Config) error {
	file, err := os.Open(cfg.LocalTar)
	if err != nil {
		return err
	}
	defer file.Close()

	tr := tar.NewReader(file)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := fmt.Sprintf("%s/%s", cfg.ExtractDir, header.Name)
		if header.Typeflag == tar.TypeDir {
			os.MkdirAll(target, 0755)
		} else {
			f, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}
	return nil
}
