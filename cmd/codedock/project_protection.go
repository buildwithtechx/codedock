package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"codedock/internal/models"
	"github.com/spf13/cobra"
)

var projectRoutingCmd = &cobra.Command{
	Use:   "routing",
	Short: "Manage project routing and edge protection",
}

var routeProtectionCmd = &cobra.Command{
	Use:   "protection",
	Short: "Manage edge password protection for project routes",
}

var routeProtectionListCmd = &cobra.Command{
	Use:   "list <project-id>",
	Short: "List route protection statuses for a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		statuses, err := client.ListRouteProtection(args[0])
		if err != nil {
			fmt.Printf("Error listing route protection: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "HOSTNAME\tENABLED\tSTATUS")
		for _, s := range statuses {
			fmt.Fprintf(w, "%s\t%t\t%s\n", s.Hostname, s.Enabled, s.Status)
		}
		w.Flush()
	},
}

var routeProtectionProtectCmd = &cobra.Command{
	Use:   "protect <project-id> <hostname>",
	Short: "Protect a route with a password challenge",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]
		hostname := args[1]
		pwdFile, _ := cmd.Flags().GetString("password-file")
		pwd, _ := cmd.Flags().GetString("password")

		password := pwd
		if pwdFile != "" {
			if pwdFile == "-" {
				reader := bufio.NewReader(os.Stdin)
				line, err := reader.ReadString('\n')
				if err != nil && err != io.EOF {
					fmt.Printf("Error reading password from stdin: %v\n", err)
					os.Exit(1)
				}
				password = strings.TrimRight(line, "\r\n")
			} else {
				content, err := os.ReadFile(pwdFile)
				if err != nil {
					fmt.Printf("Error reading password file: %v\n", err)
					os.Exit(1)
				}
				password = strings.TrimRight(string(content), "\r\n")
			}
		}

		if len(password) < 12 {
			fmt.Println("Error: Password must be at least 12 characters long")
			os.Exit(1)
		}

		client := getClient()
		res, err := client.ChangeRouteProtection(projectID, &models.RouteProtectionInput{
			Hostname: hostname,
			Action:   "protect",
			Password: password,
		})
		if err != nil {
			fmt.Printf("Error protecting route: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Route %s protected successfully (status: %s)\n", res.Hostname, res.Status)
	},
}

var routeProtectionUnlockCmd = &cobra.Command{
	Use:   "unlock <project-id> <hostname>",
	Short: "Remove password protection from a route",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]
		hostname := args[1]
		yes, _ := cmd.Flags().GetBool("yes")

		if !yes {
			fmt.Printf("Make %s public without password challenge? (y/N): ", hostname)
			var resp string
			fmt.Scanln(&resp)
			if strings.ToLower(strings.TrimSpace(resp)) != "y" {
				fmt.Println("Cancelled")
				return
			}
		}

		client := getClient()
		res, err := client.ChangeRouteProtection(projectID, &models.RouteProtectionInput{
			Hostname: hostname,
			Action:   "unlock",
		})
		if err != nil {
			fmt.Printf("Error unlocking route: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Route %s is now public (status: %s)\n", res.Hostname, res.Status)
	},
}

var routeProtectionRevokeCmd = &cobra.Command{
	Use:   "revoke <project-id> <hostname>",
	Short: "Revoke all active browser sessions for a route",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]
		hostname := args[1]

		client := getClient()
		res, err := client.ChangeRouteProtection(projectID, &models.RouteProtectionInput{
			Hostname: hostname,
			Action:   "revoke",
		})
		if err != nil {
			fmt.Printf("Error revoking sessions: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Sessions revoked for route %s (status: %s)\n", res.Hostname, res.Status)
	},
}

var routeProtectionRetryCmd = &cobra.Command{
	Use:   "retry <project-id> <hostname>",
	Short: "Retry applying route protection at edge",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		projectID := args[0]
		hostname := args[1]

		client := getClient()
		res, err := client.ChangeRouteProtection(projectID, &models.RouteProtectionInput{
			Hostname: hostname,
			Action:   "retry",
		})
		if err != nil {
			fmt.Printf("Error retrying edge protection: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Edge protection retried for route %s (status: %s)\n", res.Hostname, res.Status)
	},
}

func init() {
	routeProtectionProtectCmd.Flags().String("password", "", "Protection password (min 12 chars)")
	routeProtectionProtectCmd.Flags().String("password-file", "", "Password file or '-' for stdin")

	routeProtectionUnlockCmd.Flags().BoolP("yes", "y", false, "Confirm making the route public without prompt")

	routeProtectionCmd.AddCommand(
		routeProtectionListCmd,
		routeProtectionProtectCmd,
		routeProtectionUnlockCmd,
		routeProtectionRevokeCmd,
		routeProtectionRetryCmd,
	)

	projectRoutingCmd.AddCommand(routeProtectionCmd)
	projectCmd.AddCommand(projectRoutingCmd)
}
