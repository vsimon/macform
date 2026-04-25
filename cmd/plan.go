package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vsimon/macform/internal/diff"
	"github.com/vsimon/macform/internal/output"
	"github.com/vsimon/macform/internal/registry"
	"github.com/vsimon/macform/internal/spec"
)

var planFiles []string

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Show changes between spec and current system state",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(planFiles) > 1 {
			return fmt.Errorf("--file may only be specified once")
		}
		var planFile string
		if len(planFiles) == 1 {
			planFile = planFiles[0]
		}

		path, err := spec.Resolve(planFile)
		if err != nil {
			return err
		}

		s, err := spec.Load(path)
		if err != nil {
			return err
		}

		if err := spec.Validate(s); err != nil {
			return err
		}

		entries, err := diff.Compute(s, registry.Default)
		if err != nil {
			return err
		}

		printer := &output.Printer{NoColor: NoColor, Out: os.Stdout}
		printer.PrintPlan(entries)

		return nil
	},
}

func init() {
	planCmd.Flags().StringArrayVarP(&planFiles, "file", "f", nil, "Path to spec file")
	rootCmd.AddCommand(planCmd)
}
