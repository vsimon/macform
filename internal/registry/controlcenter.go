package registry

import "github.com/vsimon/macform/internal/provider"

var controlCenterSettings = []SettingDef{
	{
		SpecKey:  "battery-show-percentage",
		Type:     "bool",
		Provider: provider.NewDefaults("com.apple.controlcenter", "BatteryShowPercentage", "bool"),
	},
}
