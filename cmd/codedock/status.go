package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status [service-id]",
	Short: "Check the status of an application service and its recent deployments",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceID := args[0]
		client := getClient()

		service, err := client.GetService(serviceID)
		if err != nil {
			fmt.Printf("Error retrieving service status: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Service:       %s (%s)\n", service.Name, service.ID)
		fmt.Printf("Status:        %s\n", service.Status)
		if service.Domain != "" {
			fmt.Printf("Domain:        %s\n", service.Domain)
		}
		if service.RepositoryURL != "" {
			fmt.Printf("Repository:    %s (branch: %s)\n", service.RepositoryURL, service.Branch)
		}
		fmt.Printf("Created At:    %s\n", service.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated At:    %s\n", service.UpdatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()

		deployments, err := client.ListDeployments(serviceID)
		if err == nil && len(deployments) > 0 {
			fmt.Println("Recent Deployments:")
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
			fmt.Fprintln(w, "ID\tSTATUS\tBRANCH\tCOMMIT\tCREATED")
			limit := min(len(deployments), 5)
			for i := range limit {
				d := deployments[i]
				createdAt := d.CreatedAt.Format("2006-01-02 15:04:05")
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", d.ID, d.Status, d.Branch, d.CommitHash, createdAt)
			}
			w.Flush()
		} else {
			fmt.Println("No deployments found for this service.")
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	appsCmd.AddCommand(statusCmd)
}
