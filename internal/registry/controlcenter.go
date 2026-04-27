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
		SpecKey:  "show-bluetooth",
		Type:     "bool",
		Provider: provider.NewDefaults("com.apple.controlcenter", "Bluetooth", "int"),
		ValueMap: map[string]string{"true": "18", "false": "24"},
	},
	{
		SpecKey:  "show-sound",
		Type:     "string",
		Provider: provider.NewCurrentHostDefaults("com.apple.controlcenter", "Sound", "int"),
		ValueMap: map[string]string{"always": "18", "when-active": "2", "never": "24"},
	},
	{
		SpecKey:  "show-spotlight",
		Type:     "bool",
		Provider: provider.NewDefaults("com.apple.Spotlight", "MenuItemHidden", "int"),
		ValueMap: map[string]string{"true": "0", "false": "1"},
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
