package spec

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "spec*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	path := writeTemp(t, `
dock:
  autohide: true
  tile-size: 48
finder:
  show-hidden-files: false
`)
	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if s["dock"]["autohide"] != true {
		t.Errorf("expected dock.autohide=true, got %v", s["dock"]["autohide"])
	}
}

func TestLoad_NullPreserved(t *testing.T) {
	path := writeTemp(t, `
dock:
  autohide: null
`)
	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	val, ok := s["dock"]["autohide"]
	if !ok {
		t.Fatal("expected dock.autohide key to be present")
	}
	if val != nil {
		t.Errorf("expected nil for null value, got %v", val)
	}
}

func TestResolve_FlagPath(t *testing.T) {
	path := writeTemp(t, "dock:\n  autohide: true\n")
	got, err := Resolve(path)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if got != path {
		t.Errorf("got %q, want %q", got, path)
	}
}

func TestResolve_CWD(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "macform.yaml")
	if err := os.WriteFile(specPath, []byte("dock:\n  autohide: true\n"), 0644); err != nil {
		t.Fatal(err)
	}
	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	got, err := Resolve("")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	// Resolve symlinks for comparison (macOS /var → /private/var)
	gotReal, _ := filepath.EvalSymlinks(got)
	wantReal, _ := filepath.EvalSymlinks(specPath)
	if gotReal != wantReal {
		t.Errorf("got %q, want %q", got, specPath)
	}
}

func TestResolve_NotFound(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer os.Chdir(orig)
	os.Chdir(dir)

	_, err := Resolve("")
	if err == nil {
		t.Fatal("expected error when no spec file found")
	}
}

func TestValidate_UnknownSection(t *testing.T) {
	s := Spec{"unknown": {"autohide": true}}
	if err := Validate(s); err == nil {
		t.Fatal("expected error for unknown section")
	}
}

func TestValidate_UnknownKey(t *testing.T) {
	s := Spec{"dock": {"nonexistent-key": true}}
	err := Validate(s)
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestValidate_UnknownKey_DidYouMean(t *testing.T) {
	s := Spec{"dock": {"autohid": true}} // close to "autohide"
	err := Validate(s)
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestValidate_TypeMismatch(t *testing.T) {
	s := Spec{"dock": {"autohide": "yes"}} // should be bool
	if err := Validate(s); err == nil {
		t.Fatal("expected type mismatch error")
	}
}

func TestValidate_NullIsValid(t *testing.T) {
	s := Spec{"dock": {"autohide": nil}}
	if err := Validate(s); err != nil {
		t.Errorf("null should be valid (delete), got: %v", err)
	}
}

func TestValidate_ValidSpec(t *testing.T) {
	s := Spec{
		"dock":   {"autohide": true, "tile-size": 48},
		"finder": {"show-hidden-files": false, "default-view-style": "list"},
	}
	if err := Validate(s); err != nil {
		t.Errorf("expected valid spec to pass, got: %v", err)
	}
}
