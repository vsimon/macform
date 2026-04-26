package provider

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type osascriptProvider struct {
	settingName string
	readScript  string
	writeScript func(string) string
	defaultVal  string
}

var (
	osascriptRunner func(string) (string, error) = runOsascript
	retryDelay      time.Duration                = 500 * time.Millisecond
)

// NewOsascript returns a Provider that reads and writes via AppleScript GUI scripting.
// writeScript is called with the desired value and must return a complete AppleScript
// that performs the UI action (toggle checkbox, select dropdown, set text field, etc.).
// Go handles the read-before-write check, so writeScript only runs when the value
// actually needs to change.
func NewOsascript(settingName, readScript string, writeScript func(string) string, defaultVal string) Provider {
	return &osascriptProvider{
		settingName: settingName,
		readScript:  readScript,
		writeScript: writeScript,
		defaultVal:  defaultVal,
	}
}

func (p *osascriptProvider) Read() (string, bool, error) {
	out, err := runOsascriptWithRetry(p.readScript)
	if err != nil {
		return "", false, fmt.Errorf("%s: read: %w", p.settingName, err)
	}
	if out == "" {
		return "", false, nil
	}
	return out, true, nil
}

func (p *osascriptProvider) Write(value string) error {
	current, found, err := p.Read()
	if err != nil {
		return err
	}
	if found && current == value {
		return nil
	}
	if _, err := runOsascriptWithRetry(p.writeScript(value)); err != nil {
		return fmt.Errorf("%s: write: %w", p.settingName, err)
	}
	return nil
}

func (p *osascriptProvider) Delete() error {
	return p.Write(p.defaultVal)
}

func runOsascriptWithRetry(script string) (string, error) {
	out, err := osascriptRunner(script)
	if err == nil {
		return out, nil
	}
	time.Sleep(retryDelay)
	out, err = osascriptRunner(script)
	if err == nil {
		return out, nil
	}
	return "", fmt.Errorf("osascript: System Settings couldn't be reached after 2 attempts \u2014 close any open dialogs and try again (%w)", err)
}

func runOsascript(script string) (string, error) {
	out, err := exec.Command("osascript", "-e", script).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("osascript: %w \u2013 %s", err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}
