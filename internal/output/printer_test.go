package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vsimon/macform/internal/diff"
)

func newPrinter(buf *bytes.Buffer) *Printer {
	return &Printer{NoColor: true, Out: buf}
}

func TestPrintPlan_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	p := newPrinter(&buf)

	entries := []diff.DiffEntry{
		{Section: "dock", SpecKey: "autohide", Action: diff.ActionNone, CurrentVal: "true", DesiredVal: "true"},
		{Section: "finder", SpecKey: "show-hidden-files", Action: diff.ActionNone, CurrentVal: "false", DesiredVal: "false"},
	}

	p.PrintPlan(entries)

	got := buf.String()
	want := "No changes. macOS configuration is up-to-date.\n"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestPrintPlan_NoChanges_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	p := newPrinter(&buf)

	p.PrintPlan(nil)

	got := buf.String()
	want := "No changes. macOS configuration is up-to-date.\n"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestPrintPlan_MixedActions(t *testing.T) {
	var buf bytes.Buffer
	p := newPrinter(&buf)

	entries := []diff.DiffEntry{
		// dock section
		{Section: "dock", SpecKey: "autohide", Action: diff.ActionAdd, CurrentVal: "", DesiredVal: "true"},
		{Section: "dock", SpecKey: "tile-size", Action: diff.ActionChange, CurrentVal: "48", DesiredVal: "64"},
		{Section: "dock", SpecKey: "orientation", Action: diff.ActionNone, CurrentVal: "bottom", DesiredVal: "bottom"},
		// finder section
		{Section: "finder", SpecKey: "show-hidden-files", Action: diff.ActionDelete, CurrentVal: "true", DesiredVal: ""},
	}

	p.PrintPlan(entries)

	got := buf.String()
	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")

	// Expected line count: 1 section header (dock) + 2 action lines + 1 section header (finder) + 1 action line + 1 summary = 6
	wantLines := []string{
		"macform will perform the following actions:",
		"",
		"  ~ dock",
		"    + autohide:   (not set) -> true",
		"    ~ tile-size:  48 -> 64",
		"  ~ finder",
		"    - show-hidden-files:  true -> (deleted)",
		"",
		"Plan: 1 to add, 1 to change, 1 to remove.",
	}

	if len(lines) != len(wantLines) {
		t.Fatalf("expected %d lines, got %d:\n%s", len(wantLines), len(lines), got)
	}

	for i, want := range wantLines {
		if lines[i] != want {
			t.Errorf("line %d: expected %q, got %q", i+1, want, lines[i])
		}
	}
}

func TestPrintPlan_OnlyAdds(t *testing.T) {
	var buf bytes.Buffer
	p := newPrinter(&buf)

	entries := []diff.DiffEntry{
		{Section: "dock", SpecKey: "autohide", Action: diff.ActionAdd, CurrentVal: "", DesiredVal: "true"},
		{Section: "dock", SpecKey: "tile-size", Action: diff.ActionAdd, CurrentVal: "", DesiredVal: "48"},
	}

	p.PrintPlan(entries)

	got := buf.String()
	if !strings.Contains(got, "  ~ dock") {
		t.Error("expected dock section header")
	}
	if !strings.Contains(got, "    + autohide:   (not set) -> true") {
		t.Error("expected autohide add line")
	}
	if !strings.Contains(got, "    + tile-size:  (not set) -> 48") {
		t.Error("expected tile-size add line")
	}
	if !strings.Contains(got, "Plan: 2 to add, 0 to change, 0 to remove.") {
		t.Errorf("expected summary, got:\n%s", got)
	}
}

func TestPrintPlan_SectionOrder(t *testing.T) {
	var buf bytes.Buffer
	p := newPrinter(&buf)

	// Provide finder before dock to verify output is still dock-first.
	entries := []diff.DiffEntry{
		{Section: "finder", SpecKey: "show-hidden-files", Action: diff.ActionAdd, CurrentVal: "", DesiredVal: "true"},
		{Section: "dock", SpecKey: "autohide", Action: diff.ActionAdd, CurrentVal: "", DesiredVal: "false"},
	}

	p.PrintPlan(entries)

	got := buf.String()
	dockIdx := strings.Index(got, "  ~ dock")
	finderIdx := strings.Index(got, "  ~ finder")

	if dockIdx == -1 || finderIdx == -1 {
		t.Fatalf("missing section headers: %s", got)
	}
	if dockIdx > finderIdx {
		t.Errorf("expected dock before finder in output")
	}
}

func TestPrintPlan_Color(t *testing.T) {
	var buf bytes.Buffer
	p := &Printer{NoColor: false, Out: &buf}

	entries := []diff.DiffEntry{
		{Section: "dock", SpecKey: "autohide", Action: diff.ActionAdd, CurrentVal: "", DesiredVal: "true"},
	}

	p.PrintPlan(entries)

	got := buf.String()
	if !strings.Contains(got, colorGreen) {
		t.Errorf("expected green ANSI code for add, got: %q", got)
	}
	if !strings.Contains(got, colorReset) {
		t.Errorf("expected ANSI reset code, got: %q", got)
	}
}
