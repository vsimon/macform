package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vsimon/macform/internal/diff"
	"github.com/vsimon/macform/internal/output"
	"github.com/vsimon/macform/internal/registry"
	"github.com/vsimon/macform/internal/spec"
)

var applyFiles []string
var autoApprove bool

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply spec to system settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(applyFiles) > 1 {
			return fmt.Errorf("--file may only be specified once")
		}
		var applyFile string
		if len(applyFiles) == 1 {
			applyFile = applyFiles[0]
		}

		path, err := spec.Resolve(applyFile)
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

		var toApply []diff.DiffEntry
		for _, e := range entries {
			if e.Action != diff.ActionNone {
				toApply = append(toApply, e)
			}
		}

		if len(toApply) == 0 {
			fmt.Print("No changes. macOS configuration is up-to-date.\n")
			return nil
		}

		printer := &output.Printer{NoColor: NoColor, Out: os.Stdout}
		printer.PrintPlan(entries)

		if !autoApprove {
			fmt.Println()
			fmt.Println("Do you want to perform these actions?")
			fmt.Println("  macform will perform the actions described above.")
			fmt.Println("  Only 'yes' will be accepted to approve.")
			fmt.Println()
			fmt.Print("  Enter a value: ")

			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input != "yes" {
				fmt.Fprintln(os.Stderr, "Apply cancelled.")
				return fmt.Errorf("apply cancelled")
			}
			fmt.Println()
		}

		var added, changed, removed int
		for _, e := range toApply {
			def, ok := registry.Lookup(e.Section, e.SpecKey)
			if !ok {
				return fmt.Errorf("applying %s.%s: setting not found in registry", e.Section, e.SpecKey)
			}

			switch e.Action {
			case diff.ActionAdd:
				encoded := registry.Encode(def, e.DesiredVal)
				if err := def.Provider.Write(encoded); err != nil {
					return fmt.Errorf("applying %s.%s: %w", e.Section, e.SpecKey, err)
				}
				added++
				fmt.Printf("  + %s: (not set) → %s\n", e.SpecKey, e.DesiredVal)
			case diff.ActionChange:
				encoded := registry.Encode(def, e.DesiredVal)
				if err := def.Provider.Write(encoded); err != nil {
					return fmt.Errorf("applying %s.%s: %w", e.Section, e.SpecKey, err)
				}
				changed++
				fmt.Printf("  ~ %s: %s → %s\n", e.SpecKey, e.CurrentVal, e.DesiredVal)
			case diff.ActionDelete:
				if err := def.Provider.Delete(); err != nil {
					return fmt.Errorf("applying %s.%s: %w", e.Section, e.SpecKey, err)
				}
				removed++
				fmt.Printf("  - %s: %s → (deleted)\n", e.SpecKey, e.CurrentVal)
			}
		}

		seen := map[string]bool{}
		for _, e := range toApply {
			def, ok := registry.Lookup(e.Section, e.SpecKey)
			if !ok {
				continue
			}
			if def.RestartProcess != "" && !seen[def.RestartProcess] {
				seen[def.RestartProcess] = true
				exec.Command("killall", def.RestartProcess).Run() //nolint:errcheck
				fmt.Printf("  $ killall %s\n", def.RestartProcess)
			}
		}

		fmt.Println()
		if NoColor {
			fmt.Printf("Apply complete! Resources: %d added, %d changed, %d removed.\n", added, changed, removed)
		} else {
			fmt.Printf("\033[1;32mApply complete! Resources: %d added, %d changed, %d removed.\033[0m\n", added, changed, removed)
		}

		return nil
	},
}

func init() {
	applyCmd.Flags().StringArrayVarP(&applyFiles, "file", "f", nil, "Path to spec file")
	applyCmd.Flags().BoolVar(&autoApprove, "auto-approve", false, "Skip confirmation prompt")
	rootCmd.AddCommand(applyCmd)
}
