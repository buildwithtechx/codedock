package main

import (
	"fmt"
	"os"
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

func init() {
	serverCreateCmd.Flags().String("name", "", "Server name")
	serverCreateCmd.Flags().String("ip", "", "Server IP address")
	serverCreateCmd.Flags().Int("port", 0, "SSH port")
	serverCreateCmd.Flags().String("user", "", "SSH username")
	serverCreateCmd.Flags().String("key", "", "SSH key")
	serverCreateCmd.Flags().Bool("is-local", false, "Local Docker server")
	serverCmd.AddCommand(serverListCmd)
	serverCmd.AddCommand(serverCreateCmd)
	serverCmd.AddCommand(serverDeleteCmd)
	rootCmd.AddCommand(serverCmd)
}
