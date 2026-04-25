package registry

import "testing"

func TestDisplaySettingsScript_DifferentActionsProduceDifferentScripts(t *testing.T) {
	read := displaySettingsScript("return value of")
	click := displaySettingsScript("click")
	if read == click {
		t.Error("expected different scripts for different actions")
	}
}

func TestAutoBrightnessWriteScript_IgnoresArgument(t *testing.T) {
	a := autoBrightnessWriteScript("true")
	b := autoBrightnessWriteScript("false")
	if a != b {
		t.Error("write script should be identical regardless of argument")
	}
}
