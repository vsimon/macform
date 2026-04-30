package provider

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var plistBuddyRunner = func(args ...string) ([]byte, error) {
	cmd := exec.Command("/usr/libexec/PlistBuddy", args...)
	cmd.Stderr = io.Discard
	return cmd.Output()
}

func dockPlistPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "Preferences", "com.apple.dock.plist")
}

type dockAppPresenceProvider struct {
	bundleID string
}

// NewDockAppPresence returns a Provider that checks and removes a single Dock app by bundle ID.
func NewDockAppPresence(bundleID string) Provider {
	return &dockAppPresenceProvider{bundleID: bundleID}
}

// Read returns ("present", true, nil) if the app is in the Dock, ("", false, nil) if not.
func (p *dockAppPresenceProvider) Read() (string, bool, error) {
	out, err := plistBuddyRunner("-c", "Print :persistent-apps:", dockPlistPath())
	if err != nil {
		return "", false, nil
	}
	target := "bundle-identifier = " + p.bundleID
	for _, line := range strings.Split(string(out), "\n") {
		if strings.TrimSpace(line) == target {
			return "present", true, nil
		}
	}
	return "", false, nil
}

func (p *dockAppPresenceProvider) Write(_ string) error { return nil }

// Delete removes all Dock entries matching the bundle ID.
func (p *dockAppPresenceProvider) Delete() error {
	path := dockPlistPath()
	var toDelete []int
	for i := 0; ; i++ {
		out, err := plistBuddyRunner("-c", fmt.Sprintf("Print :persistent-apps:%d:tile-data:bundle-identifier", i), path)
		if err != nil {
			break
		}
		if strings.TrimSpace(string(out)) == p.bundleID {
			toDelete = append(toDelete, i)
		}
	}
	for i := len(toDelete) - 1; i >= 0; i-- {
		if _, err := plistBuddyRunner("-c", fmt.Sprintf("Delete :persistent-apps:%d", toDelete[i]), path); err != nil {
			return fmt.Errorf("removing %s from Dock: %w", p.bundleID, err)
		}
	}
	return nil
}
