package provider

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// defaultsRunner executes "defaults <args>" and returns combined output. Tests may replace this.
// For "read" operations, exit code 1 is translated to errKeyNotFound.
var defaultsRunner = func(args ...string) ([]byte, error) {
	out, err := exec.Command("defaults", args...).CombinedOutput()
	if err != nil && args[0] == "read" {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, errKeyNotFound
		}
	}
	return out, err
}

var errKeyNotFound = errors.New("defaults: key not found")

type defaultsProvider struct {
	domain string
	key    string
	typ    string
}

// NewDefaults returns a Provider bound to a single defaults key.
func NewDefaults(domain, key, typ string) Provider {
	return &defaultsProvider{domain: domain, key: key, typ: typ}
}

func (d *defaultsProvider) Read() (string, bool, error) {
	out, err := defaultsRunner("read", d.domain, d.key)
	if err != nil {
		if err == errKeyNotFound {
			return "", false, nil
		}
		return "", false, fmt.Errorf("defaults read %s %s: %w", d.domain, d.key, err)
	}
	return strings.TrimRight(string(out), "\n"), true, nil
}

func (d *defaultsProvider) Write(value string) error {
	flag, err := typeFlag(d.typ)
	if err != nil {
		return err
	}
	out, err := defaultsRunner("write", d.domain, d.key, flag, value)
	if err != nil {
		return fmt.Errorf("defaults write %s %s: %s: %w", d.domain, d.key, strings.TrimSpace(string(out)), err)
	}
	return nil
}

func (d *defaultsProvider) Delete() error {
	out, err := defaultsRunner("delete", d.domain, d.key)
	if err != nil {
		return fmt.Errorf("defaults delete %s %s: %s: %w", d.domain, d.key, strings.TrimSpace(string(out)), err)
	}
	return nil
}

func typeFlag(typ string) (string, error) {
	switch typ {
	case "bool":
		return "-bool", nil
	case "int":
		return "-int", nil
	case "float":
		return "-float", nil
	case "string":
		return "-string", nil
	default:
		return "", fmt.Errorf("unknown type %q", typ)
	}
}

type multiDefaultsProvider struct {
	primary      defaultsProvider
	extraDomains []string
}

// NewMultiDefaults returns a Provider that reads from primaryDomain and writes/deletes
// to primaryDomain and all extraDomains. Stops at the first error.
func NewMultiDefaults(primaryDomain string, extraDomains []string, key, typ string) Provider {
	return &multiDefaultsProvider{
		primary:      defaultsProvider{domain: primaryDomain, key: key, typ: typ},
		extraDomains: extraDomains,
	}
}

func (m *multiDefaultsProvider) Read() (string, bool, error) {
	return m.primary.Read()
}

func (m *multiDefaultsProvider) Write(value string) error {
	if err := m.primary.Write(value); err != nil {
		return err
	}
	for _, domain := range m.extraDomains {
		d := m.primary
		d.domain = domain
		if err := d.Write(value); err != nil {
			return err
		}
	}
	return nil
}

func (m *multiDefaultsProvider) Delete() error {
	if err := m.primary.Delete(); err != nil {
		return err
	}
	for _, domain := range m.extraDomains {
		d := m.primary
		d.domain = domain
		if err := d.Delete(); err != nil {
			return err
		}
	}
	return nil
}
