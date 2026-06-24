//go:build integration

package integration

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

var macformBin string

func TestMain(m *testing.M) {
	macformBin = os.Getenv("MACFORM")
	if macformBin == "" {
		panic("MACFORM env var not set — pass the path to the macform binary")
	}
	os.Exit(m.Run())
}

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		Setup: func(env *testscript.Env) error {
			env.Setenv("MACFORM", macformBin)
			return nil
		},
		Condition: func(cond string) (bool, error) {
			switch cond {
			case "builtin-display":
				out, _ := exec.Command("system_profiler", "SPDisplaysDataType").Output()
				return bytes.Contains(out, []byte("Built-in")), nil
			case "has-battery":
				out, _ := exec.Command("system_profiler", "SPPowerDataType").Output()
				return bytes.Contains(out, []byte("Battery Information")), nil
			}
			return false, nil
		},
	})
}
