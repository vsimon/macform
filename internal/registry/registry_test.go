package registry

import "testing"

func TestLookup_Found(t *testing.T) {
	def, ok := Lookup("dock", "tile-size")
	if !ok {
		t.Fatal("expected to find dock/tile-size")
	}
	if def.SpecKey != "tile-size" {
		t.Errorf("got SpecKey %q, want %q", def.SpecKey, "tile-size")
	}
	if def.Type != "int" {
		t.Errorf("got Type %q, want %q", def.Type, "int")
	}
	if def.Provider == nil {
		t.Error("expected non-nil Provider")
	}
}

func TestLookup_NotFound(t *testing.T) {
	_, ok := Lookup("dock", "nonexistent")
	if ok {
		t.Error("expected not to find dock/nonexistent")
	}

	_, ok = Lookup("nonexistent", "tile-size")
	if ok {
		t.Error("expected not to find nonexistent section")
	}
}

func TestEncodeDecode_ValueMapped(t *testing.T) {
	def, ok := Lookup("finder", "default-view-style")
	if !ok {
		t.Fatal("expected to find finder/default-view-style")
	}

	cases := []struct{ spec, sys string }{
		{"icon", "icnv"},
		{"list", "Nlsv"},
		{"column", "clmv"},
		{"gallery", "Flwv"},
	}
	for _, c := range cases {
		if got := Encode(def, c.spec); got != c.sys {
			t.Errorf("Encode(%q) = %q, want %q", c.spec, got, c.sys)
		}
		if got := Decode(def, c.sys); got != c.spec {
			t.Errorf("Decode(%q) = %q, want %q", c.sys, got, c.spec)
		}
	}
}

func TestEncodeDecode_NoValueMap(t *testing.T) {
	def, ok := Lookup("dock", "orientation")
	if !ok {
		t.Fatal("expected to find dock/orientation")
	}
	if got := Encode(def, "bottom"); got != "bottom" {
		t.Errorf("Encode passthrough: got %q, want %q", got, "bottom")
	}
	if got := Decode(def, "bottom"); got != "bottom" {
		t.Errorf("Decode passthrough: got %q, want %q", got, "bottom")
	}
}

func TestAllDockSettingsRegistered(t *testing.T) {
	expected := []string{
		"autohide", "tile-size", "orientation", "minimize-to-application",
		"show-recents", "magnification", "large-size", "min-effect", "scroll-to-open",
	}
	for _, key := range expected {
		if _, ok := Lookup("dock", key); !ok {
			t.Errorf("dock/%s not registered", key)
		}
	}
}

func TestAllFinderSettingsRegistered(t *testing.T) {
	expected := []string{
		"show-hidden-files", "show-extensions", "show-path-bar",
		"show-status-bar", "default-view-style", "warn-on-extension-change", "new-window-target",
	}
	for _, key := range expected {
		if _, ok := Lookup("finder", key); !ok {
			t.Errorf("finder/%s not registered", key)
		}
	}
}

func TestAllDisplaySettingsRegistered(t *testing.T) {
	expected := []string{"auto-brightness"}
	for _, key := range expected {
		def, ok := Lookup("display", key)
		if !ok {
			t.Errorf("display/%s not registered", key)
			continue
		}
		if def.Provider == nil {
			t.Errorf("display/%s: Provider is nil", key)
		}
	}
}

func TestAllTrackpadSettingsRegistered(t *testing.T) {
	expected := []string{"tap-to-click", "tracking-speed"}
	for _, key := range expected {
		def, ok := Lookup("trackpad", key)
		if !ok {
			t.Errorf("trackpad/%s not registered", key)
			continue
		}
		if def.Provider == nil {
			t.Errorf("trackpad/%s: Provider is nil", key)
		}
	}
}

func TestSectionsIncludesTrackpad(t *testing.T) {
	for _, s := range Sections() {
		if s == "trackpad" {
			return
		}
	}
	t.Error("trackpad not in Sections()")
}

func TestAllKeyboardSettingsRegistered(t *testing.T) {
	expected := []string{"repeat-rate", "repeat-delay", "function-keys", "function-key-action"}
	for _, key := range expected {
		def, ok := Lookup("keyboard", key)
		if !ok {
			t.Errorf("keyboard/%s not registered", key)
			continue
		}
		if def.Provider == nil {
			t.Errorf("keyboard/%s: Provider is nil", key)
		}
	}
}

func TestSectionsIncludesKeyboard(t *testing.T) {
	for _, s := range Sections() {
		if s == "keyboard" {
			return
		}
	}
	t.Error("keyboard not in Sections()")
}

func TestSettingDef_UserNote_Preserved(t *testing.T) {
	def, ok := Lookup("keyboard", "function-keys")
	if !ok {
		t.Fatal("expected to find keyboard/function-keys")
	}
	if len(def.UserNote) == 0 {
		t.Error("expected UserNote to be non-empty for function-keys")
	}
}
