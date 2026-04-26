package provider

import (
	"errors"
	"strings"
	"testing"
)

func TestRunOsascriptWithRetry_SucceedsOnFirstAttempt(t *testing.T) {
	calls := 0
	orig := osascriptRunner
	origDelay := retryDelay
	t.Cleanup(func() { osascriptRunner = orig; retryDelay = origDelay })
	retryDelay = 0
	osascriptRunner = func(_ string) (string, error) {
		calls++
		return "true", nil
	}

	out, err := runOsascriptWithRetry("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "true" {
		t.Errorf("got %q, want %q", out, "true")
	}
	if calls != 1 {
		t.Errorf("runner called %d times, want 1", calls)
	}
}

func TestRunOsascriptWithRetry_SucceedsOnSecondAttempt(t *testing.T) {
	calls := 0
	orig := osascriptRunner
	origDelay := retryDelay
	t.Cleanup(func() { osascriptRunner = orig; retryDelay = origDelay })
	retryDelay = 0
	osascriptRunner = func(_ string) (string, error) {
		calls++
		if calls == 1 {
			return "", errors.New("execution error: Invalid index (-1719)")
		}
		return "true", nil
	}

	out, err := runOsascriptWithRetry("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "true" {
		t.Errorf("got %q, want %q", out, "true")
	}
	if calls != 2 {
		t.Errorf("runner called %d times, want 2", calls)
	}
}

func TestRunOsascriptWithRetry_FriendlyErrorAfterBothFail(t *testing.T) {
	orig := osascriptRunner
	origDelay := retryDelay
	t.Cleanup(func() { osascriptRunner = orig; retryDelay = origDelay })
	retryDelay = 0
	osascriptRunner = func(_ string) (string, error) {
		return "", errors.New("execution error: Can't get scroll area 2 of group 1")
	}

	_, err := runOsascriptWithRetry("test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "System Settings couldn't be reached after 2 attempts") {
		t.Errorf("unexpected error: %q", err.Error())
	}
	if !strings.Contains(err.Error(), "Can't get scroll area 2") {
		t.Errorf("expected underlying error in message, got: %q", err.Error())
	}
}
