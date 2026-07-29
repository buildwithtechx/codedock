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
	outFile, err := os.Create(targetTarGzPath)
	if err != nil {
		return fmt.Errorf("failed to create target archive: %w", err)
	}
	defer outFile.Close()

	gw := gzip.NewWriter(outFile)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
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
}

func shouldIgnoreFile(relPath string) bool {
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	for _, p := range parts {
		switch p {
		case ".git", "node_modules", ".next", "dist", "build", ".cache", ".env", ".DS_Store":
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
