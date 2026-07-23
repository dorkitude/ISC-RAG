package cli

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const mediaSumRevision = "bcbcef7de3b4e702e526352d825a4ff06de2becb"

type corpusFile struct {
	name   string
	sha256 string
}

var mediaSumFiles = []corpusFile{
	{name: "train.json", sha256: "cd713e926637b34930c1c5595c9e97cec6d514f4d95921a73e3da766b5ba75af"},
	{name: "validation.json", sha256: "313dfd682b25a6151a760632c06ba2bbb6fbb585a7df5301b524bdb07cb499a0"},
	{name: "test.json", sha256: "136c241cdb4fba03395df259d1153d60a61fe274215e78ee1fc96211feff6cde"},
}

func newBootstrapCommand() *cobra.Command {
	var projectDir string

	command := &cobra.Command{
		Use:   "bootstrap",
		Short: "Validate local credentials and acquire the pinned MediaSum resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			absoluteProjectDir, err := filepath.Abs(projectDir)
			if err != nil {
				return fmt.Errorf("resolve project directory: %w", err)
			}
			if err := loadLocalEnv(absoluteProjectDir); err != nil {
				return err
			}
			if err := requireCredentials(); err != nil {
				return err
			}

			dataDir := filepath.Join(absoluteProjectDir, "resources", "mediasum")
			for _, file := range mediaSumFiles {
				if err := acquireFile(cmd.Context(), dataDir, file); err != nil {
					return err
				}
			}
			cmd.Printf("MediaSum resource ready: %s\n", dataDir)
			return nil
		},
	}

	command.Flags().StringVar(&projectDir, "project-dir", ".", "ISC-RAG project directory")
	return command
}

func loadLocalEnv(projectDir string) error {
	envFile := filepath.Join(projectDir, ".env")
	if _, err := os.Stat(envFile); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return fmt.Errorf("inspect %s: %w", envFile, err)
	}
	if err := godotenv.Load(envFile); err != nil {
		return fmt.Errorf("load %s: %w", envFile, err)
	}
	return nil
}

func requireCredentials() error {
	missing := make([]string, 0, 2)
	for _, name := range []string{"FIREWORKS_API_KEY", "TURBOPUFFER_API_KEY"} {
		if strings.TrimSpace(os.Getenv(name)) == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variable(s): %s; set them in .env or the shell environment", strings.Join(missing, ", "))
	}
	return nil
}

func acquireFile(ctx context.Context, dataDir string, file corpusFile) error {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return fmt.Errorf("create MediaSum directory: %w", err)
	}

	target := filepath.Join(dataDir, file.name)
	valid, err := hasExpectedChecksum(target, file.sha256)
	if err != nil {
		return err
	}
	if valid {
		fmt.Printf("Verified existing %s\n", file.name)
		return nil
	}

	temporary, err := os.CreateTemp(dataDir, "."+file.name+".download-*")
	if err != nil {
		return fmt.Errorf("create temporary file for %s: %w", file.name, err)
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)

	url := fmt.Sprintf("https://huggingface.co/datasets/nbroad/mediasum/resolve/%s/%s?download=true", mediaSumRevision, file.name)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request for %s: %w", file.name, err)
	}

	fmt.Printf("Downloading %s\n", file.name)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("download %s: %w", file.name, err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("download %s: unexpected HTTP status %s", file.name, response.Status)
	}

	hash := sha256.New()
	if _, err := io.Copy(io.MultiWriter(temporary, hash), response.Body); err != nil {
		return fmt.Errorf("write %s: %w", file.name, err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close %s: %w", file.name, err)
	}
	actualChecksum := fmt.Sprintf("%x", hash.Sum(nil))
	if actualChecksum != file.sha256 {
		return fmt.Errorf("checksum mismatch for %s: got %s", file.name, actualChecksum)
	}
	if err := os.Rename(temporaryName, target); err != nil {
		return fmt.Errorf("place %s: %w", file.name, err)
	}
	fmt.Printf("Verified %s\n", file.name)
	return nil
}

func hasExpectedChecksum(path, expected string) (bool, error) {
	input, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("open %s: %w", path, err)
	}
	defer input.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, input); err != nil {
		return false, fmt.Errorf("hash %s: %w", path, err)
	}
	return fmt.Sprintf("%x", hash.Sum(nil)) == expected, nil
}
