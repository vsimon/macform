package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vsimon/macform/internal/registry"
	"gopkg.in/yaml.v3"
)

// Spec is a parsed macform spec file: section → key → value.
// A nil value signals a delete (YAML null).
type Spec map[string]map[string]interface{}

// Load parses a YAML spec file at path, preserving null values as nil.
func Load(path string) (Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading spec file: %w", err)
	}
	var raw map[string]map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing spec file: %w", err)
	}
	return Spec(raw), nil
}

// Resolve returns the spec file path to use, applying FR-2 resolution:
//  1. flagPath if non-empty
//  2. macform.yaml in the current working directory
func Resolve(flagPath string) (string, error) {
	if flagPath != "" {
		if _, err := os.Stat(flagPath); err != nil {
			return "", fmt.Errorf("spec file not found: %s", flagPath)
		}
		return flagPath, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	candidate := filepath.Join(cwd, "macform.yaml")
	if _, err := os.Stat(candidate); err != nil {
		return "", fmt.Errorf("No spec file found. Run 'macform generate' to create one, or pass --file.")
	}
	return candidate, nil
}

// Validate checks all sections and keys in spec against the registry.
// Returns a combined error for all issues found.
func Validate(s Spec) error {
	var errs []string

	validSections := map[string]bool{}
	for _, sec := range registry.Sections() {
		validSections[sec] = true
	}

	for section, keys := range s {
		if !validSections[section] {
			errs = append(errs, fmt.Sprintf("unknown section %q", section))
			continue
		}
		for key, val := range keys {
			def, ok := registry.Lookup(section, key)
			if !ok {
				suggestion := didYouMean(key, registry.SectionKeys(section))
				msg := fmt.Sprintf("unknown key %q in section %q", key, section)
				if suggestion != "" {
					msg += fmt.Sprintf(" (did you mean %q?)", suggestion)
				}
				errs = append(errs, msg)
				continue
			}
			// nil is valid (delete)
			if val == nil {
				continue
			}
			if err := checkType(def, val); err != nil {
				errs = append(errs, fmt.Sprintf("%s.%s: %s", section, key, err))
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("spec validation failed:\n  %s", strings.Join(errs, "\n  "))
	}
	return nil
}

// checkType validates that val matches the expected type for def.
// For string settings with a ValueMap, it warns (non-fatal) on unrecognized values.
func checkType(def *registry.SettingDef, val interface{}) error {
	switch def.Type {
	case "bool":
		if _, ok := val.(bool); !ok {
			return fmt.Errorf("expected bool, got %T", val)
		}
	case "int":
		switch val.(type) {
		case int, int64, float64:
			// go-yaml v3 decodes integers as int
		default:
			return fmt.Errorf("expected int, got %T", val)
		}
	case "float":
		switch val.(type) {
		case float64, int, int64:
		default:
			return fmt.Errorf("expected float, got %T", val)
		}
	case "list":
		list, ok := val.([]interface{})
		if !ok {
			return fmt.Errorf("expected list, got %T", val)
		}
		for _, item := range list {
			if _, ok := item.(string); !ok {
				return fmt.Errorf("list items must be strings, got %T", item)
			}
		}
	case "string":
		s, ok := val.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", val)
		}
		// Warn (non-fatal) on unrecognized enum values
		if len(def.ValueMap) > 0 {
			if _, valid := def.ValueMap[s]; !valid {
				// Non-fatal: just note it — callers may inspect stderr
				fmt.Fprintf(os.Stderr, "warning: %s: unrecognized value %q\n", def.SpecKey, s)
			}
		}
	}
	return nil
}

// didYouMean returns the closest matching SpecKey from defs, or "" if none close.
func didYouMean(input string, defs []registry.SettingDef) string {
	best := ""
	bestDist := len(input) + 1
	for _, d := range defs {
		if strings.HasPrefix(d.SpecKey, input[:min(len(input), len(d.SpecKey))]) {
			return d.SpecKey
		}
		if dist := editDistance(input, d.SpecKey); dist < bestDist && dist <= 3 {
			bestDist = dist
			best = d.SpecKey
		}
	}
	return best
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// editDistance computes the Levenshtein distance between two strings.
func editDistance(a, b string) int {
	la, lb := len(a), len(b)
	dp := make([][]int, la+1)
	for i := range dp {
		dp[i] = make([]int, lb+1)
		dp[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		dp[0][j] = j
	}
	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + minOf3(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
			}
		}
	}
	return dp[la][lb]
}

func minOf3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
