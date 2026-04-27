package provider

import (
	"fmt"
	"strings"
)

var trackpadDomains = []string{
	"com.apple.AppleMultitouchTrackpad",
	"com.apple.driver.AppleBluetoothMultitouch.trackpad",
}

type draggingStyleProvider struct{}

func NewDraggingStyle() Provider { return draggingStyleProvider{} }

func (draggingStyleProvider) Read() (string, bool, error) {
	out, err := defaultsRunner("read", trackpadDomains[0], "TrackpadThreeFingerDrag")
	if err == nil && strings.TrimSpace(string(out)) == "1" {
		return "three-finger-drag", true, nil
	}

	out, err = defaultsRunner("read", trackpadDomains[0], "Dragging")
	if err != nil || strings.TrimSpace(string(out)) != "1" {
		return "", false, nil
	}

	out, err = defaultsRunner("read", trackpadDomains[0], "DragLock")
	if err == nil && strings.TrimSpace(string(out)) == "1" {
		return "with-drag-lock", true, nil
	}
	return "without-drag-lock", true, nil
}

func (draggingStyleProvider) Write(value string) error {
	var dragging, dragLock, threeFinger string
	switch value {
	case "with-drag-lock":
		dragging, dragLock, threeFinger = "1", "1", "0"
	case "without-drag-lock":
		dragging, dragLock, threeFinger = "1", "0", "0"
	case "three-finger-drag":
		dragging, dragLock, threeFinger = "0", "0", "1"
	default:
		return fmt.Errorf("unknown dragging style %q", value)
	}
	writes := []struct{ key, val string }{
		{"Dragging", dragging},
		{"DragLock", dragLock},
		{"TrackpadThreeFingerDrag", threeFinger},
	}
	for _, domain := range trackpadDomains {
		for _, w := range writes {
			if out, err := defaultsRunner("write", domain, w.key, "-int", w.val); err != nil {
				return fmt.Errorf("defaults write %s %s: %s: %w", domain, w.key, strings.TrimSpace(string(out)), err)
			}
		}
	}
	return nil
}

func (draggingStyleProvider) Delete() error {
	keys := []string{"Dragging", "DragLock", "TrackpadThreeFingerDrag"}
	for _, domain := range trackpadDomains {
		for _, key := range keys {
			defaultsRunner("delete", domain, key) //nolint: ignore key-not-found
		}
	}
	return nil
}
