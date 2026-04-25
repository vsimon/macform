package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/vsimon/macform/internal/diff"
	"github.com/vsimon/macform/internal/registry"
)

const (
	colorCyan   = "\033[1;96m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorReset  = "\033[0m"
)

// Printer writes plan output to Out, optionally with ANSI color.
type Printer struct {
	NoColor bool
	Out     io.Writer
}

// PrintPlan prints the diff entries grouped by section.
// ActionNone entries are filtered out.
// If no changes remain, prints "No changes. System matches spec."
// Otherwise prints section headers, per-entry lines, and a summary.
func (p *Printer) PrintPlan(entries []diff.DiffEntry) {
	var changed []diff.DiffEntry
	for _, e := range entries {
		if e.Action != diff.ActionNone {
			changed = append(changed, e)
		}
	}

	if len(changed) == 0 {
		fmt.Fprintln(p.Out, "No changes. System matches spec.")
		return
	}

	bySec := map[string][]diff.DiffEntry{}
	for _, e := range changed {
		bySec[e.Section] = append(bySec[e.Section], e)
	}

	var toAdd, toChange, toDelete int

	p.colored(colorCyan, "macform will perform the following actions:\n")
	fmt.Fprintln(p.Out)

	for _, section := range registry.Sections() {
		secEntries, ok := bySec[section]
		if !ok {
			continue
		}

		maxLen := 0
		for _, e := range secEntries {
			if len(e.SpecKey) > maxLen {
				maxLen = len(e.SpecKey)
			}
		}

		p.colored(colorYellow, fmt.Sprintf("  ~ %s\n", section))

		for _, e := range secEntries {
			pad := strings.Repeat(" ", maxLen-len(e.SpecKey))
			switch e.Action {
			case diff.ActionAdd:
				toAdd++
				p.colored(colorGreen, fmt.Sprintf("    + %s:%s  (not set) -> %s\n", e.SpecKey, pad, e.DesiredVal))
			case diff.ActionChange:
				toChange++
				p.colored(colorYellow, fmt.Sprintf("    ~ %s:%s  %s -> %s\n", e.SpecKey, pad, e.CurrentVal, e.DesiredVal))
			case diff.ActionDelete:
				toDelete++
				p.colored(colorRed, fmt.Sprintf("    - %s:%s  %s -> (deleted)\n", e.SpecKey, pad, e.CurrentVal))
			}
		}
	}

	fmt.Fprintln(p.Out)
	fmt.Fprintf(p.Out, "Plan: %d to add, %d to change, %d to remove.\n", toAdd, toChange, toDelete)
}

func (p *Printer) colored(color, text string) {
	if p.NoColor {
		fmt.Fprint(p.Out, text)
	} else {
		fmt.Fprintf(p.Out, "%s%s%s", color, text, colorReset)
	}
}
