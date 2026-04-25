package diff

import (
	"fmt"

	"github.com/vsimon/macform/internal/registry"
	"github.com/vsimon/macform/internal/spec"
)

// Registry provides ordered access to SettingDefs. Implemented by the real
// registry (via registry.Default) and by fakes in tests.
type Registry interface {
	Sections() []string
	SectionKeys(section string) []registry.SettingDef
}

// Action describes what the plan will do to a setting.
type Action int

const (
	ActionNone   Action = iota
	ActionAdd
	ActionChange
	ActionDelete
)

// DiffEntry represents a single setting in the plan output.
type DiffEntry struct {
	Section    string
	SpecKey    string
	Action     Action
	CurrentVal string
	DesiredVal string
}

// Compute calculates the diff between the spec and the current system state.
func Compute(s spec.Spec, reg Registry) ([]DiffEntry, error) {
	var entries []DiffEntry

	for _, section := range reg.Sections() {
		sectionKeys, ok := s[section]
		if !ok {
			continue
		}

		for _, def := range reg.SectionKeys(section) {
			specVal, ok := sectionKeys[def.SpecKey]
			if !ok {
				continue
			}

			sysVal, found, err := def.Provider.Read()
			if err != nil {
				return nil, fmt.Errorf("reading %s.%s: %w", section, def.SpecKey, err)
			}

			currentDecoded := ""
			if found {
				currentDecoded = registry.Decode(&def, sysVal)
			}

			if specVal == nil {
				entries = append(entries, DiffEntry{
					Section:    section,
					SpecKey:    def.SpecKey,
					Action:     ActionDelete,
					CurrentVal: currentDecoded,
				})
				continue
			}

			desired := fmt.Sprintf("%v", specVal)

			if !found {
				entries = append(entries, DiffEntry{
					Section:    section,
					SpecKey:    def.SpecKey,
					Action:     ActionAdd,
					DesiredVal: desired,
				})
				continue
			}

			if currentDecoded != desired {
				entries = append(entries, DiffEntry{
					Section:    section,
					SpecKey:    def.SpecKey,
					Action:     ActionChange,
					CurrentVal: currentDecoded,
					DesiredVal: desired,
				})
			} else {
				entries = append(entries, DiffEntry{
					Section:    section,
					SpecKey:    def.SpecKey,
					Action:     ActionNone,
					CurrentVal: currentDecoded,
					DesiredVal: desired,
				})
			}
		}
	}

	return entries, nil
}
