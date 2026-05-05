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

		flatSpec, expandedReg := registry.Expand(s)
		entries, err := diff.Compute(flatSpec, expandedReg)
		if err != nil {
			return err
		}

		var toApply []diff.DiffEntry
		for _, e := range entries {
			if e.Action != diff.ActionNone {
				toApply = append(toApply, e)
			}
		}

		printer := &output.Printer{NoColor: NoColor, Out: os.Stdout}

		if len(toApply) == 0 {
			printer.PrintPlan(entries)
			return nil
		}

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
			def, ok := expandedReg.Lookup(e.Section, e.SpecKey)
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
			case diff.ActionChange:
				encoded := registry.Encode(def, e.DesiredVal)
				if err := def.Provider.Write(encoded); err != nil {
					return fmt.Errorf("applying %s.%s: %w", e.Section, e.SpecKey, err)
				}
				changed++
			case diff.ActionDelete:
				if err := def.Provider.Delete(); err != nil {
					return fmt.Errorf("applying %s.%s: %w", e.Section, e.SpecKey, err)
				}
				removed++
			}
		}

		seenKill := map[string]bool{}
		for _, e := range toApply {
			def, ok := expandedReg.Lookup(e.Section, e.SpecKey)
			if !ok {
				continue
			}
			if def.RestartProcess != "" && !seenKill[def.RestartProcess] {
				seenKill[def.RestartProcess] = true
				exec.Command("killall", def.RestartProcess).Run() //nolint:errcheck
			}
		}

		notes := map[string][]string{}
		for _, e := range toApply {
			if def, ok := expandedReg.Lookup(e.Section, e.SpecKey); ok {
				if len(def.UserNote) > 0 {
					notes[e.Section+"/"+e.SpecKey] = def.UserNote
				}
			}
		}
		printer.PrintAudit(toApply, notes)

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
