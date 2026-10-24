package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Manage DNS and ingress tunnels",
}

var ingressCmd = &cobra.Command{
	Use:   "ingress",
	Short: "Manage Cloudflare ingress tunnels",
}

var ingressListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ingress tunnel connections",
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		connections, err := client.ListIngressConnections()
		if err != nil {
			fmt.Printf("Error listing ingress connections: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tSERVER\tNAME\tCONNECTOR\tSTATUS")
		for _, c := range connections {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", c.ID, c.ServerID, c.Name, c.Connector, c.Status)
		}
		w.Flush()
	},
}

var ingressConnectCmd = &cobra.Command{
	Use:   "connect <server-id> <tunnel-name>",
	Short: "Connect a Cloudflare ingress tunnel to a server",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		serverID := args[0]
		tunnelName := args[1]
		connector, _ := cmd.Flags().GetString("connector")
		if connector == "" {
			connector = "managed"
		}

		client := getClient()
		conn, err := client.ConnectIngress(serverID, tunnelName, connector)
		if err != nil {
			fmt.Printf("Error connecting ingress: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Ingress tunnel connected: %s (ID: %s, Status: %s)\n", conn.Name, conn.ID, conn.Status)
	},
}

var ingressRetryCmd = &cobra.Command{
	Use:   "retry <id>",
	Short: "Retry reconciliation for an ingress tunnel",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		client := getClient()
		conn, err := client.RetryIngress(id)
		if err != nil {
			fmt.Printf("Error retrying ingress: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Ingress reconciliation retried: %s (Status: %s)\n", conn.Name, conn.Status)
	},
}

var ingressDisconnectCmd = &cobra.Command{
	Use:   "disconnect <id>",
	Short: "Disconnect an ingress tunnel",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		client := getClient()
		if err := client.DisconnectIngress(id); err != nil {
			fmt.Printf("Error disconnecting ingress: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Ingress tunnel %s disconnected successfully\n", id)
	},
}

func init() {
	ingressConnectCmd.Flags().String("connector", "managed", "Connector type (managed or existing)")

	ingressCmd.AddCommand(ingressListCmd)
	ingressCmd.AddCommand(ingressConnectCmd)
	ingressCmd.AddCommand(ingressRetryCmd)
	ingressCmd.AddCommand(ingressDisconnectCmd)

	dnsCmd.AddCommand(ingressCmd)
	rootCmd.AddCommand(dnsCmd)
}
