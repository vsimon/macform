package cmd

import (
	"testing"
)

func TestPlanCmd_MultipleFFiles_ReturnsError(t *testing.T) {
	planFiles = []string{"a.yaml", "b.yaml"}
	defer func() { planFiles = nil }()

	err := planCmd.RunE(planCmd, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "--file may only be specified once" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApplyCmd_MultipleFFiles_ReturnsError(t *testing.T) {
	applyFiles = []string{"a.yaml", "b.yaml"}
	defer func() { applyFiles = nil }()

	err := applyCmd.RunE(applyCmd, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "--file may only be specified once" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGenerateCmd_MultipleFFiles_ReturnsError(t *testing.T) {
	generateFiles = []string{"a.yaml", "b.yaml"}
	defer func() { generateFiles = nil }()

	err := generateCmd.RunE(generateCmd, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "--file may only be specified once" {
		t.Fatalf("unexpected error: %v", err)
	}
}
