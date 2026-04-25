package diff

import (
	"testing"

	"github.com/vsimon/macform/internal/registry"
	"github.com/vsimon/macform/internal/spec"
)

// mockProvider implements provider.Provider for testing.
type mockProvider struct {
	value string
	found bool
	err   error
}

func (m *mockProvider) Read() (string, bool, error) { return m.value, m.found, m.err }
func (m *mockProvider) Write(_ string) error        { return nil }
func (m *mockProvider) Delete() error               { return nil }

// fakeRegistry implements Registry for testing.
type fakeRegistry struct {
	sectionOrder []string
	keys         map[string][]registry.SettingDef
}

func (f *fakeRegistry) Sections() []string { return f.sectionOrder }
func (f *fakeRegistry) SectionKeys(section string) []registry.SettingDef {
	return f.keys[section]
}

func TestCompute_ActionAdd(t *testing.T) {
	reg := &fakeRegistry{
		sectionOrder: []string{"dock"},
		keys: map[string][]registry.SettingDef{
			"dock": {{SpecKey: "autohide", Provider: &mockProvider{found: false}, Type: "bool"}},
		},
	}
	s := spec.Spec{"dock": {"autohide": true}}

	entries, err := Compute(s, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Action != ActionAdd {
		t.Errorf("expected ActionAdd, got %v", e.Action)
	}
	if e.Section != "dock" {
		t.Errorf("expected section dock, got %s", e.Section)
	}
	if e.SpecKey != "autohide" {
		t.Errorf("expected specKey autohide, got %s", e.SpecKey)
	}
	if e.CurrentVal != "" {
		t.Errorf("expected empty CurrentVal, got %q", e.CurrentVal)
	}
	if e.DesiredVal != "true" {
		t.Errorf("expected DesiredVal=true, got %q", e.DesiredVal)
	}
}

func TestCompute_ActionChange(t *testing.T) {
	reg := &fakeRegistry{
		sectionOrder: []string{"dock"},
		keys: map[string][]registry.SettingDef{
			"dock": {{SpecKey: "tile-size", Provider: &mockProvider{value: "48", found: true}, Type: "int"}},
		},
	}
	s := spec.Spec{"dock": {"tile-size": 64}}

	entries, err := Compute(s, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Action != ActionChange {
		t.Errorf("expected ActionChange, got %v", e.Action)
	}
	if e.CurrentVal != "48" {
		t.Errorf("expected CurrentVal=48, got %q", e.CurrentVal)
	}
	if e.DesiredVal != "64" {
		t.Errorf("expected DesiredVal=64, got %q", e.DesiredVal)
	}
}

func TestCompute_ActionDelete(t *testing.T) {
	reg := &fakeRegistry{
		sectionOrder: []string{"dock"},
		keys: map[string][]registry.SettingDef{
			"dock": {{SpecKey: "autohide", Provider: &mockProvider{value: "true", found: true}, Type: "bool"}},
		},
	}
	s := spec.Spec{"dock": {"autohide": nil}}

	entries, err := Compute(s, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Action != ActionDelete {
		t.Errorf("expected ActionDelete, got %v", e.Action)
	}
	if e.CurrentVal != "true" {
		t.Errorf("expected CurrentVal=true, got %q", e.CurrentVal)
	}
	if e.DesiredVal != "" {
		t.Errorf("expected empty DesiredVal, got %q", e.DesiredVal)
	}
}

func TestCompute_ActionNone(t *testing.T) {
	reg := &fakeRegistry{
		sectionOrder: []string{"dock"},
		keys: map[string][]registry.SettingDef{
			"dock": {{SpecKey: "autohide", Provider: &mockProvider{value: "true", found: true}, Type: "bool"}},
		},
	}
	s := spec.Spec{"dock": {"autohide": true}}

	entries, err := Compute(s, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Action != ActionNone {
		t.Errorf("expected ActionNone, got %v", e.Action)
	}
	if e.CurrentVal != "true" {
		t.Errorf("expected CurrentVal=true, got %q", e.CurrentVal)
	}
}

func TestCompute_SectionOrder(t *testing.T) {
	reg := &fakeRegistry{
		sectionOrder: []string{"dock", "finder"},
		keys: map[string][]registry.SettingDef{
			"dock":   {{SpecKey: "autohide", Provider: &mockProvider{found: false}, Type: "bool"}},
			"finder": {{SpecKey: "show-hidden-files", Provider: &mockProvider{found: false}, Type: "bool"}},
		},
	}
	s := spec.Spec{
		"dock":   {"autohide": false},
		"finder": {"show-hidden-files": true},
	}

	entries, err := Compute(s, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Section != "dock" {
		t.Errorf("expected dock first, got %s", entries[0].Section)
	}
	if entries[1].Section != "finder" {
		t.Errorf("expected finder second, got %s", entries[1].Section)
	}
}

func TestCompute_ValueMap_ActionChange(t *testing.T) {
	reg := &fakeRegistry{
		sectionOrder: []string{"finder"},
		keys: map[string][]registry.SettingDef{
			"finder": {{
				SpecKey:  "default-view-style",
				Provider: &mockProvider{value: "Nlsv", found: true},
				Type:     "string",
				ValueMap: map[string]string{
					"icon": "icnv", "list": "Nlsv", "column": "clmv", "gallery": "Flwv",
				},
			}},
		},
	}
	s := spec.Spec{"finder": {"default-view-style": "icon"}}

	entries, err := Compute(s, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Action != ActionChange {
		t.Errorf("expected ActionChange, got %v", e.Action)
	}
	if e.CurrentVal != "list" {
		t.Errorf("expected CurrentVal=list (decoded from Nlsv), got %q", e.CurrentVal)
	}
	if e.DesiredVal != "icon" {
		t.Errorf("expected DesiredVal=icon, got %q", e.DesiredVal)
	}
}

func TestCompute_ValueMap_ActionNone(t *testing.T) {
	reg := &fakeRegistry{
		sectionOrder: []string{"finder"},
		keys: map[string][]registry.SettingDef{
			"finder": {{
				SpecKey:  "default-view-style",
				Provider: &mockProvider{value: "icnv", found: true},
				Type:     "string",
				ValueMap: map[string]string{
					"icon": "icnv", "list": "Nlsv", "column": "clmv", "gallery": "Flwv",
				},
			}},
		},
	}
	s := spec.Spec{"finder": {"default-view-style": "icon"}}

	entries, err := Compute(s, reg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Action != ActionNone {
		t.Errorf("expected ActionNone, got %v", entries[0].Action)
	}
}
