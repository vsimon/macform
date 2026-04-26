package registry

import "github.com/vsimon/macform/internal/provider"

var controlCenterSettings = []SettingDef{
	{
		SpecKey:  "show-battery",
		Type:     "bool",
		Provider: provider.NewDefaults("com.apple.controlcenter", "Battery", "int"),
		ValueMap: map[string]string{"true": "18", "false": "24"},
	},
	{
		SpecKey:  "show-wifi",
		Type:     "bool",
		Provider: provider.NewDefaults("com.apple.controlcenter", "WiFi", "int"),
		ValueMap: map[string]string{"true": "18", "false": "24"},
	},
	{
		SpecKey:  "show-battery-percentage",
		Type:     "bool",
		Provider: provider.NewDefaults("com.apple.controlcenter", "BatteryShowPercentage", "bool"),
	},
}
