package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"codedock/internal/models"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:     "servers",
	Aliases: []string{"server"},
	Short:   "Manage servers",
}

var serverListCmd = &cobra.Command{
	Use:   "list",
	Short: "List servers",
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		servers, err := client.ListServers()
		if err != nil {
			fmt.Printf("Error listing servers: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tHOST\tTYPE\tSTATUS")
		for _, s := range servers {
			host := s.IPAddress
			if host == "" {
				host = "127.0.0.1"
			}
			serverType := "Remote SSH"
			if s.IsLocal {
				serverType = "Local Docker"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", s.ID, s.Name, host, serverType, s.Status)
		}
		w.Flush()
	},
}

var serverCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a server",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		ip, _ := cmd.Flags().GetString("ip")
		port, _ := cmd.Flags().GetInt("port")
		user, _ := cmd.Flags().GetString("user")
		key, _ := cmd.Flags().GetString("key")
		isLocal, _ := cmd.Flags().GetBool("is-local")

		if name == "" {
			fmt.Println("Error: --name flag is required")
			os.Exit(1)
		}
		if ip == "" {
			ip = "127.0.0.1"
		}
		if port == 0 {
			port = 22
		}
		if user == "" {
			user = "root"
		}

		client := getClient()
		server, err := client.CreateServer(&models.CreateServerRequest{
			Name:      name,
			IPAddress: ip,
			IsLocal:   isLocal,
			SSHHost:   ip,
			SSHPort:   port,
			SSHUser:   user,
			SSHKey:    key,
		})
		if err != nil {
			fmt.Printf("Error creating server: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Server %s created successfully with ID: %s\n", server.Name, server.ID)
	},
}

var serverDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		if err := client.DeleteServer(args[0]); err != nil {
			fmt.Printf("Error deleting server: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Server %s deleted successfully\n", args[0])
	},
}

var serverNetworkCmd = &cobra.Command{
	Use:   "network [id]",
	Short: "Show a server's outbound network settings",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		settings, err := getClient().GetServerNetworkSettings(args[0])
		if err != nil {
			fmt.Printf("Error loading network settings: %v\n", err)
			os.Exit(1)
		}
		internet := "unknown"
		if settings.InternetAccess != nil {
			internet = "disabled"
			if *settings.InternetAccess {
				internet = "enabled"
			}
		}
		egress := "*"
		if len(settings.Egress) > 0 {
			egress = strings.Join(settings.Egress, ", ")
		}
		fmt.Printf("Server:           %s\n", args[0])
		fmt.Printf("Internet access:  %s\n", internet)
		fmt.Printf("Egress allowlist: %s\n", egress)
		if settings.OutboundIP != "" {
			fmt.Printf("Outbound gateway: %s\n", settings.OutboundIP)
		}
		if settings.OutboundMode != "" {
			fmt.Printf("Outbound mode:    %s\n", settings.OutboundMode)
		}
		if settings.Revision != "" {
			fmt.Printf("Revision:         %s\n", settings.Revision)
		}
	},
}

var serverNetworkSetCmd = &cobra.Command{
	Use:   "set [id]",
	Short: "Update a server's outbound network settings",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		current, err := client.GetServerNetworkSettings(args[0])
		if err != nil {
			fmt.Printf("Error loading network settings: %v\n", err)
			os.Exit(1)
		}
		internet := current.InternetAccess == nil || *current.InternetAccess
		if cmd.Flags().Changed("internet") {
			internet, _ = cmd.Flags().GetBool("internet")
		}
		egress := current.Egress
		if raw, _ := cmd.Flags().GetString("egress"); strings.TrimSpace(raw) != "" {
			egress = nil
			for _, entry := range strings.Split(raw, ",") {
				if trimmed := strings.ToLower(strings.TrimSpace(entry)); trimmed != "" {
					egress = append(egress, trimmed)
				}
			}
		}
		if confirm, _ := cmd.Flags().GetBool("confirm"); !confirm {
			fmt.Println("Refusing to change outbound settings without --confirm")
			os.Exit(1)
		}
		updated, err := client.UpdateServerNetworkSettings(args[0], &models.UpdateServerNetworkSettingsRequest{
			InternetAccess:         internet,
			ExpectedInternetAccess: current.InternetAccess == nil || *current.InternetAccess,
			Egress:                 egress,
			ExpectedRevision:       current.Revision,
			Confirm:                true,
		})
		if err != nil {
			fmt.Printf("Error updating network settings: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Network settings updated for server %s (%d egress entries)\n", args[0], len(updated.Egress))
	},
}

func init() {
	serverCreateCmd.Flags().String("name", "", "Server name")
	serverCreateCmd.Flags().String("ip", "", "Server IP address")
	serverCreateCmd.Flags().Int("port", 0, "SSH port")
	serverCreateCmd.Flags().String("user", "", "SSH username")
	serverCreateCmd.Flags().String("key", "", "SSH key")
	serverCreateCmd.Flags().Bool("is-local", false, "Local Docker server")
	serverNetworkSetCmd.Flags().Bool("internet", true, "Outbound internet access")
	serverNetworkSetCmd.Flags().String("egress", "", "Comma-separated egress allowlist (replaces the current list)")
	serverNetworkSetCmd.Flags().Bool("confirm", false, "Confirm the outbound change")
	serverNetworkCmd.AddCommand(serverNetworkSetCmd)
	serverCmd.AddCommand(serverListCmd)
	serverCmd.AddCommand(serverCreateCmd)
	serverCmd.AddCommand(serverDeleteCmd)
	serverCmd.AddCommand(serverNetworkCmd)
	rootCmd.AddCommand(serverCmd)
}
