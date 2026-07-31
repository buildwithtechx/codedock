package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CreateTarContext(sourceDir string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("failed to determine relative path for %s: %w", path, err)
		}
		if relPath == "." {
			return nil
		}
		if shouldIgnoreFile(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		return addFileToTar(tw, path, info, relPath)
	})
	if err != nil {
		return nil, err
	}
	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close tar writer: %w", err)
	}
	return &buf, nil
}

func CreateTarGzArchive(sourceDir, targetTarGzPath string) error {
	outFile, err := os.OpenFile(targetTarGzPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create target archive: %w", err)
	}
	defer outFile.Close()

	gw := gzip.NewWriter(outFile)
	tw := tar.NewWriter(gw)

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("failed to determine relative path for %s: %w", path, err)
		}
		if relPath == "." {
			return nil
		}
		if shouldIgnoreFile(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		return addFileToTar(tw, path, info, relPath)
	})
	if err != nil {
		return err
	}
	if err := tw.Close(); err != nil {
		return fmt.Errorf("failed to close tar writer: %w", err)
	}
	if err := gw.Close(); err != nil {
		return fmt.Errorf("failed to close gzip writer: %w", err)
	}
	return nil
}

func shouldIgnoreFile(relPath string) bool {
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	for _, p := range parts {
		switch p {
		case ".git", "node_modules", ".next", "dist", "build", ".cache", ".DS_Store":
			return true
		}
		if p == ".env" || strings.HasPrefix(p, ".env.") {
			return true
		}
	}
	return false
}

func addFileToTar(tw *tar.Writer, path string, info os.FileInfo, relPath string) error {
	if err := writeTarHeader(tw, info, relPath); err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	return writeTarFileContent(tw, path)
}

func writeTarHeader(tw *tar.Writer, info os.FileInfo, relPath string) error {
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return fmt.Errorf("failed to create tar header for %s: %w", relPath, err)
	}
	header.Name = relPath
	if err := tw.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write tar header for %s: %w", relPath, err)
	}
	return nil
}

func writeTarFileContent(tw *tar.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()
	if _, err := io.Copy(tw, file); err != nil {
		return fmt.Errorf("failed to copy file %s into tar: %w", path, err)
	}
	return nil
}
