package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

const DefaultAddr = "127.0.0.1:4173"

func ReportHandler(repoRoot string) (http.Handler, error) {
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("resolve repo root: %w", err)
	}

	reportDir := filepath.Join(repoRoot, "report")
	indexPath := filepath.Join(reportDir, "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("report not prepared: %s", repoRoot)
		}
		return nil, fmt.Errorf("check report html %s: %w", indexPath, err)
	}

	return http.FileServer(http.Dir(reportDir)), nil
}

func Serve(repoRoot, addr string) error {
	if addr == "" {
		addr = DefaultAddr
	}

	handler, err := ReportHandler(repoRoot)
	if err != nil {
		return err
	}

	fmt.Printf("serving d2a report at http://%s\n", addr)
	return http.ListenAndServe(addr, handler)
}
