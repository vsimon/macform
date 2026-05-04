package registry

import "github.com/vsimon/macform/internal/provider"

var appearanceSettings = []SettingDef{
	{
		SpecKey:  "show-scroll-bars",
		Type:     "string",
		Provider: provider.NewDefaults("NSGlobalDomain", "AppleShowScrollBars", "string"),
		ValueMap: map[string]string{
			"automatic": "Automatic", "when-scrolling": "WhenScrolling", "always": "Always",
		},
	},
}
