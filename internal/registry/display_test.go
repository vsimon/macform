package registry

import (
	"strings"
	"testing"
)

func TestDisplaySettingsScript_ContainsDialogDismissal(t *testing.T) {
	script := displaySettingsScript("return value of")
	if !strings.Contains(script, "count of sheets of window 1") {
		t.Error("script missing sheet count check for dialog dismissal")
	}
	if !strings.Contains(script, "key code 53") {
		t.Error("script missing Escape keystroke (key code 53)")
	}
	if !strings.Contains(script, "System Settings\" to activate") {
		t.Error("script missing activate call — required for accessibility tree to render")
	}
}

func TestDisplaySettingsScript_ContainsElementPolling(t *testing.T) {
	script := displaySettingsScript("return value of")
	if !strings.Contains(script, "repeat while waited < 10") {
		t.Error("script missing element polling loop")
	}
	if !strings.Contains(script, "set cb to missing value") {
		t.Error("script missing checkbox variable initialization")
	}
	if !strings.Contains(script, "exit repeat") {
		t.Error("script missing exit repeat on checkbox found")
	}
}

func TestDisplaySettingsScript_ContainsBothCheckboxPaths(t *testing.T) {
	script := displaySettingsScript("return value of")
	if !strings.Contains(script, "group 3 of splitter group 1") {
		t.Error("script missing Tahoe checkbox path (group 3)")
	}
	if !strings.Contains(script, "group 2 of splitter group 1") {
		t.Error("script missing legacy checkbox path (group 2)")
	}
}

func TestDisplaySettingsScript_NoFixedDelayTwo(t *testing.T) {
	script := displaySettingsScript("return value of")
	if strings.Contains(script, "delay 2") {
		t.Error("script still contains hard-coded 'delay 2' — should use element polling instead")
	}
}

func TestDisplaySettingsScript_ContainsInnerAction(t *testing.T) {
	script := displaySettingsScript("return value of")
	if !strings.Contains(script, "return value of") {
		t.Error("script missing innerAction")
	}
	script2 := displaySettingsScript("click")
	if !strings.Contains(script2, "click") {
		t.Error("script missing innerAction for click variant")
	}
}
