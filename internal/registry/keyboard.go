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
	{
		SpecKey:  "function-keys",
		Type:     "string",
		Provider: provider.NewDefaults("NSGlobalDomain", "com.apple.keyboard.fnState", "bool"),
		ValueMap: map[string]string{
			"special": "false", "standard": "true",
		},
	},
	{
		SpecKey:  "function-key-action",
		Type:     "string",
		Provider: provider.NewDefaults("com.apple.HIToolbox", "AppleFnUsageType", "int"),
		ValueMap: map[string]string{
			"do-nothing": "0", "change-input-source": "1", "show-emoji": "2", "start-dictation": "3",
		},
	},
	{
		SpecKey:  "auto-capitalize",
		Type:     "bool",
		Provider: provider.NewDefaults("NSGlobalDomain", "NSAutomaticCapitalizationEnabled", "bool"),
	},
	{
		SpecKey:  "auto-correct",
		Type:     "bool",
		Provider: provider.NewDefaults("NSGlobalDomain", "NSAutomaticSpellingCorrectionEnabled", "bool"),
	},
	{
		SpecKey:  "press-and-hold",
		Type:     "bool",
		Provider: provider.NewDefaults("NSGlobalDomain", "ApplePressAndHoldEnabled", "bool"),
	},
	{
		SpecKey:  "smart-dashes",
		Type:     "bool",
		Provider: provider.NewDefaults("NSGlobalDomain", "NSAutomaticDashSubstitutionEnabled", "bool"),
	},
	{
		SpecKey:  "double-space-period",
		Type:     "bool",
		Provider: provider.NewDefaults("NSGlobalDomain", "NSAutomaticPeriodSubstitutionEnabled", "bool"),
	},
	{
		SpecKey:  "keyboard-navigation",
		Type:     "bool",
		Provider: provider.NewDefaults("NSGlobalDomain", "AppleKeyboardUIMode", "int"),
		ValueMap: map[string]string{
			"true": "2", "false": "0",
		},
	},
}
