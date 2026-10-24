package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var aiModeFlag string

var aiCmd = &cobra.Command{
	Use:   "ai [prompt...]",
	Short: "Ask the instance-aware Codedock AI Copilot to inspect resources or propose actions",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runAIAsk(strings.Join(args, " "))
	},
}

var aiAskCmd = &cobra.Command{
	Use:   "ask [prompt...]",
	Short: "Send a prompt to the Codedock AI Copilot",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runAIAsk(strings.Join(args, " "))
	},
}

var aiConfirmCmd = &cobra.Command{
	Use:   "confirm <action-id>",
	Short: "Approve and execute a pending AI action by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := getClient()
		res, err := c.ConfirmAIAction(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error confirming action: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Action %s executed successfully.\n", args[0])
		if tool, ok := res["tool"].(string); ok && tool != "" {
			fmt.Printf("Tool: %s\n", tool)
		}
	},
}

func runAIAsk(prompt string) {
	c := getClient()
	resp, err := c.AskAI(prompt, aiModeFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error querying Codedock AI: %v\n", err)
		os.Exit(1)
	}

	if len(resp.Steps) > 0 {
		fmt.Println("Executed MCP Steps:")
		for _, s := range resp.Steps {
			fmt.Printf("  • [%s] %s — %s\n", s.Status, s.Tool, s.Label)
		}
		fmt.Println()
	}

	fmt.Println(resp.Reply)

	if len(resp.PendingActions) > 0 {
		fmt.Println("\nPending Actions (Require Approval):")
		for _, a := range resp.PendingActions {
			fmt.Printf("  • %s — %s (%s %s)\n", a.ActionID, a.Title, a.Method, a.Path)
			fmt.Printf("    Run: codedock ai confirm %s\n", a.ActionID)
		}
	}
}

func init() {
	aiCmd.PersistentFlags().StringVar(&aiModeFlag, "mode", "agent", "Assistant mode: 'agent' or 'plan'")
	aiCmd.AddCommand(aiAskCmd)
	aiCmd.AddCommand(aiConfirmCmd)
	rootCmd.AddCommand(aiCmd)
}
