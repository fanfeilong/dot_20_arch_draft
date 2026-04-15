package updater

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	defaultRepo = "fanfeilong/dot_20_arch_draft"
)

// SelfUpdate downloads the latest d2a release asset and replaces the current executable.
func SelfUpdate() (string, error) {
	if runtime.GOOS == "windows" {
		return "", fmt.Errorf("self-update via -U is not supported on windows; use install.ps1")
	}

	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("resolve executable path: %w", err)
	}
	if resolved, err := filepath.EvalSymlinks(exePath); err == nil {
		exePath = resolved
	}

	asset, err := releaseAssetName()
	if err != nil {
		return "", err
	}

	repo := getenv("D2A_REPO", defaultRepo)
	version := getenv("D2A_VERSION", "latest")
	baseURL := strings.TrimSpace(os.Getenv("D2A_BASE_URL"))
	if baseURL == "" {
		if version == "latest" {
			baseURL = fmt.Sprintf("https://github.com/%s/releases/latest/download", repo)
		} else {
			baseURL = fmt.Sprintf("https://github.com/%s/releases/download/%s", repo, version)
		}
	}

	tmpDir, err := os.MkdirTemp("", "d2a-update-*")
	if err != nil {
		return "", fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, asset)
	if err := downloadFile(fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), asset), archivePath); err != nil {
		return "", err
	}

	newBinaryPath := filepath.Join(tmpDir, "d2a.new")
	if err := extractBinary(archivePath, newBinaryPath); err != nil {
		return "", err
	}
	if err := os.Chmod(newBinaryPath, 0o755); err != nil {
		return "", fmt.Errorf("chmod updated binary: %w", err)
	}

	replacePath := exePath + ".tmp-update"
	if err := os.Rename(newBinaryPath, replacePath); err != nil {
		return "", fmt.Errorf("prepare replacement binary: %w", err)
	}
	if err := os.Rename(replacePath, exePath); err != nil {
		return "", fmt.Errorf("replace executable: %w", err)
	}
	return exePath, nil
}

func releaseAssetName() (string, error) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	var arch string
	switch goarch {
	case "amd64":
		arch = "amd64"
	case "arm64":
		arch = "arm64"
	default:
		return "", fmt.Errorf("unsupported architecture: %s", goarch)
	}

	switch goos {
	case "darwin", "linux":
		return fmt.Sprintf("d2a_%s_%s.tar.gz", goos, arch), nil
	case "windows":
		return fmt.Sprintf("d2a_windows_%s.zip", arch), nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", goos)
	}
}

func downloadFile(url, dst string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download release asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download release asset: unexpected status %s", resp.Status)
	}

	file, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create archive file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("write archive file: %w", err)
	}
	return nil
}

func extractBinary(archivePath, dst string) error {
	if strings.HasSuffix(archivePath, ".tar.gz") {
		return extractTarGzBinary(archivePath, dst)
	}
	if strings.HasSuffix(archivePath, ".zip") {
		return extractZipBinary(archivePath, dst)
	}
	return fmt.Errorf("unsupported archive format: %s", archivePath)
}

func extractTarGzBinary(archivePath, dst string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("open gzip stream: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read tar entry: %w", err)
		}

		if hdr.FileInfo().IsDir() {
			continue
		}
		if filepath.Base(hdr.Name) != "d2a" {
			continue
		}

		out, err := os.Create(dst)
		if err != nil {
			return fmt.Errorf("create extracted binary: %w", err)
		}
		if _, err := io.Copy(out, tr); err != nil {
			out.Close()
			return fmt.Errorf("extract binary: %w", err)
		}
		if err := out.Close(); err != nil {
			return fmt.Errorf("close extracted binary: %w", err)
		}
		return nil
	}
	return fmt.Errorf("binary d2a not found in archive")
}

func extractZipBinary(archivePath, dst string) error {
	zr, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("open zip archive: %w", err)
	}
	defer zr.Close()

	for _, file := range zr.File {
		if filepath.Base(file.Name) != "d2a.exe" {
			continue
		}

		src, err := file.Open()
		if err != nil {
			return fmt.Errorf("open zip entry: %w", err)
		}

		out, err := os.Create(dst)
		if err != nil {
			src.Close()
			return fmt.Errorf("create extracted binary: %w", err)
		}
		if _, err := io.Copy(out, src); err != nil {
			out.Close()
			src.Close()
			return fmt.Errorf("extract binary: %w", err)
		}
		if err := out.Close(); err != nil {
			src.Close()
			return fmt.Errorf("close extracted binary: %w", err)
		}
		if err := src.Close(); err != nil {
			return fmt.Errorf("close zip entry: %w", err)
		}
		return nil
	}
	return fmt.Errorf("binary d2a.exe not found in archive")
}

func getenv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
