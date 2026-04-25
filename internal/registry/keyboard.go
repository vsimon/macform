package registry

import "github.com/vsimon/macform/internal/provider"

var keyboardSettings = []SettingDef{
	{
		SpecKey:  "repeat-rate",
		Type:     "int",
		Provider: provider.NewDefaults("NSGlobalDomain", "KeyRepeat", "int"),
	},
	{
		SpecKey:  "repeat-delay",
		Type:     "int",
		Provider: provider.NewDefaults("NSGlobalDomain", "InitialKeyRepeat", "int"),
	},
}
