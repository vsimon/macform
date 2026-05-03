package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vsimon/macform/internal/registry"
	"gopkg.in/yaml.v3"
)

var generateFiles []string
var generateForce bool

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Snapshot current system settings to a spec file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(generateFiles) > 1 {
			return fmt.Errorf("--file may only be specified once")
		}
		var generateFile string
		if len(generateFiles) == 1 {
			generateFile = generateFiles[0]
		}

		var outPath string
		if generateFile != "" {
			outPath = generateFile
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			outPath = filepath.Join(cwd, "macform.yaml")
		}

		if _, err := os.Stat(outPath); err == nil && !generateForce {
			fmt.Fprintf(os.Stderr, "File %s already exists.\n", outPath)
			fmt.Print("Overwrite? Only 'yes' will be accepted: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			input := strings.TrimSpace(scanner.Text())
			if input != "yes" {
				fmt.Fprintln(os.Stderr, "Generate cancelled.")
				return fmt.Errorf("generate cancelled")
			}
		}

		out := make(map[string]map[string]interface{})

		for _, section := range registry.Sections() {
			sectionMap := make(map[string]interface{})
			for _, def := range registry.SectionKeys(section) {
				if def.Type == "list" {
					continue
				}
				sysVal, found, err := def.Provider.Read()
				if err != nil {
					return err
				}
				if !found {
					continue
				}
				specVal := registry.Decode(&def, sysVal)
				if specVal == sysVal && len(def.ValueMap) > 0 {
					fmt.Fprintf(os.Stderr, "warning: %s.%s: no mapping for system value %q, using as-is\n", section, def.SpecKey, sysVal)
				}
				sectionMap[def.SpecKey] = coerceType(def.Type, specVal)
			}
			if len(sectionMap) > 0 {
				out[section] = sectionMap
			}
		}

		yamlBytes, err := yaml.Marshal(out)
		if err != nil {
			return err
		}

		dir := filepath.Dir(outPath)
		tmp, err := os.CreateTemp(dir, "macform-*.yaml")
		if err != nil {
			return err
		}
		if _, err := tmp.Write(yamlBytes); err != nil {
			tmp.Close()
			os.Remove(tmp.Name())
			return err
		}
		if err := tmp.Close(); err != nil {
			os.Remove(tmp.Name())
			return err
		}
		if err := os.Rename(tmp.Name(), outPath); err != nil {
			os.Remove(tmp.Name())
			return err
		}

		fmt.Fprintf(os.Stdout, "Generated %s\n", outPath)
		return nil
	},
}

func coerceType(typ, val string) interface{} {
	switch typ {
	case "bool":
		return val == "1" || val == "true"
	case "int":
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	case "float":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return val
}

func init() {
	generateCmd.Flags().StringArrayVarP(&generateFiles, "file", "f", nil, "Output path for spec file")
	generateCmd.Flags().BoolVar(&generateForce, "force", false, "Overwrite existing file without prompting")
	rootCmd.AddCommand(generateCmd)
}
