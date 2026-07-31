package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"codedock.run/codedock/internal/utils"
	"codedock.run/codedock/pkg/types"
)

var (
	deployProjectID string
	deployAppName   string
)

var deployCmd = &cobra.Command{
	Use:   "deploy [path|service-id]",
	Short: "Deploy a local directory or trigger deployment for an existing service",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := "."
		if len(args) > 0 {
			target = args[0]
		}

		fi, err := os.Stat(target)
		if err == nil && fi.IsDir() {
			runLocalDirectoryDeploy(target)
			return
		}

		if len(args) == 1 {
			runServiceTriggerDeploy(args[0])
			return
		}

		fmt.Println("❌ Specified target is not a valid directory or service ID")
		os.Exit(1)
	},
}

func runLocalDirectoryDeploy(targetDir string) {
	absDir, err := filepath.Abs(targetDir)
	if err != nil {
		fmt.Printf("❌ Invalid target directory: %v\n", err)
		os.Exit(1)
	}

	client := getClient()

	projectID := deployProjectID
	if projectID == "" {
		projects, err := client.ListProjects()
		if err != nil {
			fmt.Printf("❌ Failed to list projects: %v\n", err)
			os.Exit(1)
		}
		if len(projects) > 0 {
			projectID = projects[0].ID
		} else {
			p, err := client.CreateProject(&types.CreateProjectRequest{
				Name:        "default-project",
				Description: "Default Project",
			})
			if err != nil {
				fmt.Printf("❌ Failed to create project: %v\n", err)
				os.Exit(1)
			}
			projectID = p.ID
		}
	}

	appName := deployAppName
	if appName == "" {
		appName = filepath.Base(absDir)
	}

	fmt.Printf("📦 Packaging local directory: %s...\n", absDir)
	tmpArchive := filepath.Join(os.TempDir(), fmt.Sprintf("codedock-cli-%s.tar.gz", uuid.New().String()[:8]))
	if err := utils.CreateTarGzArchive(absDir, tmpArchive); err != nil {
		fmt.Printf("❌ Failed to package archive: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(tmpArchive)

	fmt.Printf("🚀 Deploying %s to Codedock...\n", appName)
	res, err := client.DeployArchive(projectID, appName, tmpArchive)
	if err != nil {
		fmt.Printf("❌ Deployment failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Deployment complete!")
	fmt.Printf("   App Name: %s\n", res.AppName)
	fmt.Printf("   App ID:   %s\n", res.AppID)
	fmt.Printf("   Container: %s\n", res.ContainerID)
}

func runServiceTriggerDeploy(serviceID string) {
	client := getClient()

	fmt.Printf("🚀 Triggering deployment for service %s...\n", serviceID)
	deployment, err := client.TriggerDeployment(serviceID)
	if err != nil {
		fmt.Printf("❌ Failed to trigger deployment: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Deployment started!\n")
	fmt.Printf("ID: %s\n", deployment.ID)
	fmt.Printf("Status: %s\n", deployment.Status)
	fmt.Printf("To check status, run: codedock status %s\n", serviceID)
}

func init() {
	deployCmd.Flags().StringVarP(&deployProjectID, "project", "p", "", "Target project ID")
	deployCmd.Flags().StringVarP(&deployAppName, "name", "n", "", "App name override")
	rootCmd.AddCommand(deployCmd)
}
