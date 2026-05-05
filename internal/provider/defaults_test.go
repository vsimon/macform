package provider

import (
	"testing"
)

func TestDefaults_Read_NormalizesBoolInts(t *testing.T) {
	for _, tc := range []struct{ raw, want string }{
		{"0\n", "false"},
		{"1\n", "true"},
		{"false\n", "false"},
		{"true\n", "true"},
	} {
		orig := defaultsRunner
		raw := tc.raw
		defaultsRunner = func(args ...string) ([]byte, error) { return []byte(raw), nil }
		t.Cleanup(func() { defaultsRunner = orig })

		p := NewDefaults("NSGlobalDomain", "TestKey", "bool")
		val, _, err := p.Read()
		if err != nil {
			t.Fatalf("raw=%q: Read failed: %v", tc.raw, err)
		}
		if val != tc.want {
			t.Errorf("raw=%q: got %q, want %q", tc.raw, val, tc.want)
		}
	}
}

func TestMultiDefaults_Read_UsesPrimaryDomain(t *testing.T) {
	orig := defaultsRunner
	defaultsRunner = func(args ...string) ([]byte, error) {
		if args[0] == "read" && args[1] == "primary" {
			return []byte("primary-value\n"), nil
		}
		return []byte("extra-value\n"), nil
	}
	t.Cleanup(func() { defaultsRunner = orig })

	p := NewMultiDefaults("primary", []string{"extra"}, "TestKey", "string")
	val, found, err := p.Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if !found {
		t.Fatal("Read: key not found")
	}
	if val != "primary-value" {
		t.Errorf("got %q, want %q", val, "primary-value")
	}
}

func TestMultiDefaults_Write_WritesToAllDomains(t *testing.T) {
	var gotDomains []string
	orig := defaultsRunner
	defaultsRunner = func(args ...string) ([]byte, error) {
		if args[0] == "write" {
			gotDomains = append(gotDomains, args[1])
		}
		return nil, nil
	}
	t.Cleanup(func() { defaultsRunner = orig })

	p := NewMultiDefaults("primary", []string{"extra"}, "TestKey", "string")
	if err := p.Write("hello"); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	wantDomains := []string{"primary", "extra"}
	if len(gotDomains) != len(wantDomains) {
		t.Fatalf("wrote to %d domains, want %d: %v", len(gotDomains), len(wantDomains), gotDomains)
	}
	for i, want := range wantDomains {
		if gotDomains[i] != want {
			t.Errorf("domain[%d]: got %q, want %q", i, gotDomains[i], want)
		}
	}
}

func TestMultiDefaults_Delete_DeletesFromAllDomains(t *testing.T) {
	var gotDomains []string
	orig := defaultsRunner
	defaultsRunner = func(args ...string) ([]byte, error) {
		if args[0] == "delete" {
			gotDomains = append(gotDomains, args[1])
		}
		return nil, nil
	}
	t.Cleanup(func() { defaultsRunner = orig })

	p := NewMultiDefaults("primary", []string{"extra"}, "TestKey", "string")
	if err := p.Delete(); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	wantDomains := []string{"primary", "extra"}
	if len(gotDomains) != len(wantDomains) {
		t.Fatalf("deleted from %d domains, want %d: %v", len(gotDomains), len(wantDomains), gotDomains)
	}
	for i, want := range wantDomains {
		if gotDomains[i] != want {
			t.Errorf("domain[%d]: got %q, want %q", i, gotDomains[i], want)
		}
	}
}
