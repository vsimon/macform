package registry

import "github.com/vsimon/macform/internal/provider"

var hotCornerActionMap = map[string]string{
	"no-op":                "1",
	"mission-control":      "2",
	"application-windows":  "3",
	"desktop":              "4",
	"start-screen-saver":   "5",
	"disable-screen-saver": "6",
	"dashboard":            "7",
	"put-display-to-sleep": "10",
	"launchpad":            "11",
	"notification-center":  "12",
	"lock-screen":          "13",
	"quick-note":           "14",
}

var hotCornerModifierMap = map[string]string{
	"none":    "0",
	"shift":   "131072",
	"control": "262144",
	"option":  "524288",
	"command": "1048576",
}

var hotCornerSettings = []SettingDef{
	{
		SpecKey: "top-left", Type: "string", RestartCommand: killDock,
		Provider: provider.NewDefaults("com.apple.dock", "wvous-tl-corner", "int"),
		ValueMap: hotCornerActionMap,
	},
	{
		SpecKey: "top-left-modifier", Type: "string", RestartCommand: killDock,
		Provider: provider.NewDefaults("com.apple.dock", "wvous-tl-modifier", "int"),
		ValueMap: hotCornerModifierMap,
	},
	{
		SpecKey: "top-right", Type: "string", RestartCommand: killDock,
		Provider: provider.NewDefaults("com.apple.dock", "wvous-tr-corner", "int"),
		ValueMap: hotCornerActionMap,
	},
	{
		SpecKey: "top-right-modifier", Type: "string", RestartCommand: killDock,
		Provider: provider.NewDefaults("com.apple.dock", "wvous-tr-modifier", "int"),
		ValueMap: hotCornerModifierMap,
	},
	{
		SpecKey: "bottom-left", Type: "string", RestartCommand: killDock,
		Provider: provider.NewDefaults("com.apple.dock", "wvous-bl-corner", "int"),
		ValueMap: hotCornerActionMap,
	},
	{
		SpecKey: "bottom-left-modifier", Type: "string", RestartCommand: killDock,
		Provider: provider.NewDefaults("com.apple.dock", "wvous-bl-modifier", "int"),
		ValueMap: hotCornerModifierMap,
	},
	{
		SpecKey: "bottom-right", Type: "string", RestartCommand: killDock,
		Provider: provider.NewDefaults("com.apple.dock", "wvous-br-corner", "int"),
		ValueMap: hotCornerActionMap,
	},
	{
		SpecKey: "bottom-right-modifier", Type: "string", RestartCommand: killDock,
		Provider: provider.NewDefaults("com.apple.dock", "wvous-br-modifier", "int"),
		ValueMap: hotCornerModifierMap,
	},
}
