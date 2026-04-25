package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// NoColor disables ANSI color output when true.
var NoColor bool

var rootCmd = &cobra.Command{
	Use:   "macform",
	Short: "Declarative macOS system settings management",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if os.Getenv("NO_COLOR") != "" {
			NoColor = true
		}
	},
	SilenceUsage: true,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&NoColor, "no-color", false, "Disable color output")
	rootCmd.SetVersionTemplate(fmt.Sprintf("macform %s (commit: %s, built: %s)\n", version, commit, date))
	rootCmd.Version = version
}
