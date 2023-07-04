package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply configuration changes",
	Run:   applyRun,
}

func applyRun(cmd *cobra.Command, args []string) {
	// TODO: Implement apply command
}
