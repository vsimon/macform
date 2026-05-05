package registry

import "github.com/vsimon/macform/internal/provider"

var trackpadSettings = []SettingDef{
	{
		SpecKey: "tap-to-click",
		Type:    "bool",
		Provider: provider.NewMultiDefaults(
			"com.apple.AppleMultitouchTrackpad",
			[]string{"com.apple.driver.AppleBluetoothMultitouch.trackpad"},
			"Clicking", "bool",
		),
		UserNote: []string{"# tap-to-click requires logout or restart to take full effect"},
	},
	{
		SpecKey:  "tracking-speed",
		Type:     "float",
		Provider: provider.NewDefaults("NSGlobalDomain", "com.apple.trackpad.scaling", "float"),
		UserNote: []string{"# tracking-speed requires logout or restart to take full effect"},
	},
	{
		SpecKey:  "dragging-style",
		Type:     "string",
		Provider: provider.NewDraggingStyle(),
		UserNote: []string{"# dragging-style requires logout or restart to take full effect"},
	},
	{
		SpecKey:  "natural-scrolling",
		Type:     "bool",
		Provider: provider.NewDefaults("NSGlobalDomain", "com.apple.swipescrolldirection", "bool"),
	},
}
