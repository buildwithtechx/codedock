package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"codedock/cmd/codedock/client"
	"codedock/cmd/codedock/config"
	"github.com/spf13/cobra"
)

var (
	mcpServerFlag string
	mcpTokenFlag  string
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Model Context Protocol (MCP) server bridge and tool inspection",
	Long:  `Connect local AI editors (Claude Desktop, Cursor, VS Code, Windsurf, Zed) to your Codedock instance via stdio or inspect available MCP tools.`,
	Run: func(cmd *cobra.Command, args []string) {
		runMCPStdioBridge()
	},
}

var mcpServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run a stdio JSON-RPC 2.0 bridge to the remote Codedock MCP server",
	Run: func(cmd *cobra.Command, args []string) {
		runMCPStdioBridge()
	},
}

var mcpStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display MCP server capabilities and transport status",
	Run: func(cmd *cobra.Command, args []string) {
		c := resolveMCPClient()
		status, err := c.GetMCPStatus()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error retrieving MCP status: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Server:           %s (%s)\n", status.ServerName, status.ServerVersion)
		fmt.Printf("Protocol Version: %s\n", status.ProtocolVersion)
		fmt.Printf("Enabled:          %v\n", status.Enabled)
		fmt.Printf("Tools Available:  %d\n", status.ToolCount)
		fmt.Printf("Prompts:          %d\n", status.PromptCount)
		fmt.Printf("Endpoint (HTTP):  %s/api/mcp\n", strings.TrimRight(c.BaseURL, "/"))
		fmt.Printf("Endpoint (SSE):   %s/api/mcp/sse\n", strings.TrimRight(c.BaseURL, "/"))
	},
}

var mcpToolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "List all MCP tools exposed by the connected Codedock instance",
	Run: func(cmd *cobra.Command, args []string) {
		c := resolveMCPClient()
		tools, err := c.ListMCPTools()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing MCP tools: %v\n", err)
			os.Exit(1)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "TOOL\tCATEGORY\tMETHOD\tDESCRIPTION")
		for _, t := range tools {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", t.Name, t.Category, t.Method, t.Description)
		}
		w.Flush()
	},
}

func resolveMCPClient() *client.Client {
	cfg, _ := config.Load()
	baseURL := mcpServerFlag
	if baseURL == "" {
		baseURL = os.Getenv("CODEDOCK_URL")
	}
	if baseURL == "" && cfg != nil {
		baseURL = cfg.ServerURL
	}
	token := mcpTokenFlag
	if token == "" {
		token = os.Getenv("CODEDOCK_TOKEN")
	}
	if token == "" && cfg != nil {
		token = cfg.Token
	}
	if baseURL == "" || token == "" {
		fmt.Fprintln(os.Stderr, "Error: Not authenticated. Run 'codedock login' or pass --server and --token.")
		os.Exit(1)
	}
	return client.NewClient(strings.TrimRight(baseURL, "/"), token)
}

func runMCPStdioBridge() {
	c := resolveMCPClient()
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		respBytes, err := c.SendMCPJSONRPC(line)
		if err != nil {
			errResp := fmt.Sprintf(`{"jsonrpc":"2.0","id":null,"error":{"code":-32603,"message":%q}}`, err.Error())
			fmt.Fprintln(os.Stdout, errResp)
			continue
		}
		if len(respBytes) > 0 {
			fmt.Fprintln(os.Stdout, string(bytes.TrimSpace(respBytes)))
		}
	}
}

func init() {
	mcpCmd.PersistentFlags().StringVar(&mcpServerFlag, "server", "", "Codedock instance URL (overrides saved config)")
	mcpCmd.PersistentFlags().StringVar(&mcpTokenFlag, "token", "", "Personal Access Token (overrides saved config)")

	mcpCmd.AddCommand(mcpServeCmd)
	mcpCmd.AddCommand(mcpStatusCmd)
	mcpCmd.AddCommand(mcpToolsCmd)
	rootCmd.AddCommand(mcpCmd)
}
