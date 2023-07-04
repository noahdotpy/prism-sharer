package cmd

import (
	"github.com/spf13/cobra"

	"github.com/noahdotpy/prism-sharer/core"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var (
	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply configuration changes",
		Run:   core.ApplyRun,
	}
)
