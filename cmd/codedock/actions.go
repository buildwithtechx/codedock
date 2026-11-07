package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var actionsCmd = &cobra.Command{
	Use:   "actions",
	Short: "Manage CI workflows, execution runs, and runners",
}

var actionsWorkflowsCmd = &cobra.Command{
	Use:   "workflows",
	Short: "List action workflows",
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		workflows, err := client.ListWorkflows()
		if err != nil {
			fmt.Printf("Error listing workflows: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTRIGGER\tENABLED")
		for _, wf := range workflows {
			fmt.Fprintf(w, "%s\t%s\t%s\t%t\n", wf.ID, wf.Name, wf.Trigger, wf.Enabled)
		}
		w.Flush()
	},
}

var actionsRunsCmd = &cobra.Command{
	Use:   "runs",
	Short: "List action execution runs",
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		runs, err := client.ListRuns()
		if err != nil {
			fmt.Printf("Error listing action runs: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\t#\tSTATUS\tTRIGGERED_BY\tDURATION")
		for _, r := range runs {
			fmt.Fprintf(w, "%s\t%s\t#%d\t%s\t%s\t%ds\n", r.ID, r.Name, r.Number, r.Status, r.TriggeredBy, r.DurationS)
		}
		w.Flush()
	},
}

var actionsRunRetryCmd = &cobra.Command{
	Use:   "retry <run-id>",
	Short: "Retry an action run",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		client := getClient()
		run, err := client.RetryRun(id)
		if err != nil {
			fmt.Printf("Error retrying action run: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Action run retried: %s (#%d, Status: %s)\n", run.Name, run.Number, run.Status)
	},
}

var actionsRunCancelCmd = &cobra.Command{
	Use:   "cancel <run-id>",
	Short: "Cancel an active action run",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		client := getClient()
		run, err := client.CancelRun(id)
		if err != nil {
			fmt.Printf("Error cancelling action run: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Action run cancelled: %s (#%d, Status: %s)\n", run.Name, run.Number, run.Status)
	},
}

var actionsRunnersCmd = &cobra.Command{
	Use:   "runners",
	Short: "List action execution runners",
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		runners, err := client.ListRunners()
		if err != nil {
			fmt.Printf("Error listing runners: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tMODE\tSTATUS\tENABLED")
		for _, r := range runners {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%t\n", r.ID, r.Name, r.Mode, r.Status, r.Enabled)
		}
		w.Flush()
	},
}

func init() {
	actionsRunsCmd.AddCommand(actionsRunRetryCmd)
	actionsRunsCmd.AddCommand(actionsRunCancelCmd)

	actionsCmd.AddCommand(actionsWorkflowsCmd)
	actionsCmd.AddCommand(actionsRunsCmd)
	actionsCmd.AddCommand(actionsRunnersCmd)

	rootCmd.AddCommand(actionsCmd)
}
