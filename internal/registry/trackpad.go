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
	},
	{
		SpecKey:  "tracking-speed",
		Type:     "float",
		Provider: provider.NewDefaults("NSGlobalDomain", "com.apple.trackpad.scaling", "float"),
	},
}
