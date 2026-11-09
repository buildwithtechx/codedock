package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"codedock/internal/platform/storage"
)

type deployArgs struct {
	gitURL      string
	imageRef    string
	archivePath string
	dirName     string
	projectID   string
	branch      string
	rootDir     string
	port        int
}

func parseDeployArgs(args []string) deployArgs {
	d := deployArgs{
		branch: "main",
		port:   3000,
	}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--git", "-g":
			if i+1 < len(args) {
				d.gitURL = args[i+1]
				i++
			}
		case "--image", "-i":
			if i+1 < len(args) {
				d.imageRef = args[i+1]
				i++
			}
		case "--archive", "-a":
			if i+1 < len(args) {
				d.archivePath = args[i+1]
				i++
			}
		case "--project", "-p":
			if i+1 < len(args) {
				d.projectID = args[i+1]
				i++
			}
		case "--branch", "-b":
			if i+1 < len(args) {
				d.branch = args[i+1]
				i++
			}
		case "--dir", "-d":
			if i+1 < len(args) {
				d.rootDir = args[i+1]
				i++
			}
		case "--template", "-t":
			if i+1 < len(args) {
				templateName := args[i+1]
				d.gitURL = "https://github.com/buildwithtechx/codedock-examples.git"
				d.branch = "main"
				d.rootDir = templateName
				i++
			}
		case "--port":
			if i+1 < len(args) {
				if p, err := parseUint(args[i+1]); err == nil {
					d.port = p
				}
				i++
			}
		default:
			if fi, err := os.Stat(args[i]); err == nil && fi.IsDir() {
				absDir, err := filepath.Abs(args[i])
				if err != nil {
					exitError("Failed to get absolute directory: %v", err)
				}
				tmpArchive := filepath.Join(os.TempDir(), fmt.Sprintf("codedock-deploy-%s.tar.gz", uuid.New().String()[:8]))
				if err := storage.CreateTarGzArchive(absDir, tmpArchive); err != nil {
					os.Remove(tmpArchive)
					exitError("Failed to create archive: %v", err)
				}
				d.archivePath = tmpArchive
				d.dirName = filepath.Base(absDir)
			} else if strings.HasPrefix(args[i], "http") || strings.HasPrefix(args[i], "git@") {
				d.gitURL = args[i]
			} else if strings.Contains(args[i], ":") || strings.Contains(args[i], "/") {
				d.imageRef = args[i]
			}
		}
	}
	return d
}

func resolveAppName(gitURL, imageRef, dirName string) string {
	if dirName != "" {
		return dirName
	}
	appName := extractRepoName(gitURL)
	if appName == "app" && imageRef != "" {
		appName = imageRef
		if idx := strings.LastIndex(appName, "/"); idx >= 0 {
			appName = appName[idx+1:]
		}
		appName = strings.Split(appName, ":")[0]
	}
	return appName
}
