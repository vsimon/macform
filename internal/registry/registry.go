package registry

import "github.com/vsimon/macform/internal/provider"

// SettingDef describes a single managed macOS setting.
// Scalar settings set Provider; list settings set ProviderFor (a factory called per item).
type SettingDef struct {
	SpecKey        string
	Provider       provider.Provider
	ProviderFor    func(string) provider.Provider
	Type           string
	ValueMap       map[string]string
	RestartProcess string
}

// defaultRegistry implements the Registry interface used by the diff package.
type defaultRegistry struct{}

// Default is the singleton registry for use in production callers.
var Default defaultRegistry

func (defaultRegistry) Sections() []string                { return Sections() }
func (defaultRegistry) SectionKeys(s string) []SettingDef { return SectionKeys(s) }

// sections holds all registered settings in deterministic order.
var sections = map[string][]SettingDef{
	"dock": {
		{
			SpecKey: "autohide", Type: "bool", RestartProcess: "Dock",
			Provider: provider.NewDefaults("com.apple.dock", "autohide", "bool"),
		},
		{
			SpecKey: "tile-size", Type: "int", RestartProcess: "Dock",
			Provider: provider.NewDefaults("com.apple.dock", "tilesize", "int"),
		},
		{
			SpecKey: "orientation", Type: "string", RestartProcess: "Dock",
			Provider: provider.NewDefaults("com.apple.dock", "orientation", "string"),
		},
		{
			SpecKey: "minimize-to-application", Type: "bool", RestartProcess: "Dock",
			Provider: provider.NewDefaults("com.apple.dock", "minimize-to-application", "bool"),
		},
		{
			SpecKey: "show-recents", Type: "bool", RestartProcess: "Dock",
			Provider: provider.NewDefaults("com.apple.dock", "show-recents", "bool"),
		},
		{
			SpecKey: "magnification", Type: "bool", RestartProcess: "Dock",
			Provider: provider.NewDefaults("com.apple.dock", "magnification", "bool"),
		},
		{
			SpecKey: "large-size", Type: "int", RestartProcess: "Dock",
			Provider: provider.NewDefaults("com.apple.dock", "largesize", "int"),
		},
		{
			SpecKey: "min-effect", Type: "string", RestartProcess: "Dock",
			Provider: provider.NewDefaults("com.apple.dock", "mineffect", "string"),
		},
		{
			SpecKey: "scroll-to-open", Type: "bool", RestartProcess: "Dock",
			Provider: provider.NewDefaults("com.apple.dock", "scroll-to-open", "bool"),
		},
		{
			SpecKey: "remove-apps", Type: "list", RestartProcess: "Dock",
			ProviderFor: provider.NewDockAppPresence,
		},
	},
	"finder": {
		{
			SpecKey: "show-hidden-files", Type: "bool", RestartProcess: "Finder",
			Provider: provider.NewDefaults("com.apple.finder", "AppleShowAllFiles", "bool"),
		},
		{
			SpecKey: "show-extensions", Type: "bool", RestartProcess: "Finder",
			Provider: provider.NewDefaults("NSGlobalDomain", "AppleShowAllExtensions", "bool"),
		},
		{
			SpecKey: "show-path-bar", Type: "bool", RestartProcess: "Finder",
			Provider: provider.NewDefaults("com.apple.finder", "ShowPathbar", "bool"),
		},
		{
			SpecKey: "show-status-bar", Type: "bool", RestartProcess: "Finder",
			Provider: provider.NewDefaults("com.apple.finder", "ShowStatusBar", "bool"),
		},
		{
			SpecKey: "default-view-style", Type: "string", RestartProcess: "Finder",
			Provider: provider.NewDefaults("com.apple.finder", "FXPreferredViewStyle", "string"),
			ValueMap: map[string]string{
				"icon": "icnv", "list": "Nlsv", "column": "clmv", "gallery": "Flwv",
			},
		},
		{
			SpecKey: "warn-on-extension-change", Type: "bool", RestartProcess: "Finder",
			Provider: provider.NewDefaults("com.apple.finder", "FXEnableExtensionChangeWarning", "bool"),
		},
		{
			SpecKey: "new-window-target", Type: "string", RestartProcess: "Finder",
			Provider: provider.NewDefaults("com.apple.finder", "NewWindowTarget", "string"),
			ValueMap: map[string]string{
				"recents": "PfAF", "home": "PfHm", "desktop": "PfDe", "documents": "PfDo", "computer": "PfCm", "volumes": "PfVo", "icloud-drive": "PfID",
			},
		},
	},
	"display": {
		{
			SpecKey:  "auto-brightness",
			Type:     "bool",
			Provider: provider.NewOsascript("auto-brightness", autoBrightnessReadScript, autoBrightnessWriteScript, "true"),
		},
	},
	"battery": {
		{
			SpecKey:  "slightly-dim-on-battery",
			Type:     "bool",
			Provider: provider.NewOsascript("slightly-dim-on-battery", slightlyDimReadScript, slightlyDimWriteScript, "true"),
		},
	},
	"control-center": controlCenterSettings,
	"trackpad":       trackpadSettings,
	"keyboard":       keyboardSettings,
	"hot-corners":    hotCornerSettings,
}

// ExpandedRegistry is a flattened registry produced by Expand: list settings are replaced
// by per-item scalar SettingDefs so that diff and apply need no special-case handling.
type ExpandedRegistry struct {
	sections map[string][]SettingDef
}

func (r *ExpandedRegistry) Sections() []string                      { return Sections() }
func (r *ExpandedRegistry) SectionKeys(section string) []SettingDef { return r.sections[section] }
func (r *ExpandedRegistry) Lookup(section, specKey string) (*SettingDef, bool) {
	for i := range r.sections[section] {
		if r.sections[section][i].SpecKey == specKey {
			return &r.sections[section][i], true
		}
	}
	return nil, false
}

// Expand returns a flat spec and an ExpandedRegistry suitable for passing to diff.Compute.
// List settings (ProviderFor != nil) are expanded: each item becomes a nil-valued spec entry
// (triggering ActionDelete) backed by a per-item Provider from ProviderFor.
// The input map is not modified.
func Expand(s map[string]map[string]interface{}) (map[string]map[string]interface{}, *ExpandedRegistry) {
	flatSpec := make(map[string]map[string]interface{})
	reg := &ExpandedRegistry{sections: make(map[string][]SettingDef)}

	for _, section := range Sections() {
		sectionSpec, inSpec := s[section]
		if !inSpec {
			reg.sections[section] = SectionKeys(section)
			continue
		}

		flatSpec[section] = make(map[string]interface{})
		var defs []SettingDef

		for _, def := range SectionKeys(section) {
			if def.ProviderFor == nil {
				if v, ok := sectionSpec[def.SpecKey]; ok {
					flatSpec[section][def.SpecKey] = v
				}
				defs = append(defs, def)
			} else {
				items, _ := sectionSpec[def.SpecKey].([]interface{})
				for _, item := range items {
					id, _ := item.(string)
					flatSpec[section][id] = nil
					defs = append(defs, SettingDef{
						SpecKey:        id,
						Type:           "string",
						RestartProcess: def.RestartProcess,
						Provider:       def.ProviderFor(id),
					})
				}
			}
		}

		reg.sections[section] = defs
	}

	return flatSpec, reg
}

// Lookup finds a SettingDef by section and spec key in the static registry.
func Lookup(section, specKey string) (*SettingDef, bool) {
	defs, ok := sections[section]
	if !ok {
		return nil, false
	}
	for i := range defs {
		if defs[i].SpecKey == specKey {
			return &defs[i], true
		}
	}
	return nil, false
}

// SectionKeys returns all settings for a given section, in registration order.
func SectionKeys(section string) []SettingDef {
	return sections[section]
}

// Sections returns the ordered list of section names.
func Sections() []string {
	return []string{"dock", "finder", "display", "battery", "control-center", "trackpad", "keyboard", "hot-corners"}
}

// Encode converts a spec value to its system (defaults) representation.
func Encode(def *SettingDef, specVal string) string {
	if len(def.ValueMap) > 0 {
		if sysVal, ok := def.ValueMap[specVal]; ok {
			return sysVal
		}
	}
	return specVal
}

// Decode converts a system value back to its spec representation.
func Decode(def *SettingDef, sysVal string) string {
	if len(def.ValueMap) > 0 {
		for specVal, sv := range def.ValueMap {
			if sv == sysVal {
				return specVal
			}
		}
	}
	if def.Type == "bool" {
		switch sysVal {
		case "0":
			return "false"
		case "1":
			return "true"
		}
	}
	return sysVal
}
