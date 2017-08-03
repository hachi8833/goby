// +build linux

package vm

import (
	"os"
	"testing"
)

func TestCallingPluginFunction(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	p = import "github.com/goby-lang/goby/test_fixtures/import_test/plugin/plugin.go"
	p.send("Foo", "!")
	p.send("Baz")
	`

	v := initTestVM()
	// We don't test the result here for two reasons:
	// - If it doesn't work it'll returns error or panic
	// - It's hard to test a plugin obj
	v.testEval(t, input)
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func TestCallingPluginFunctionWithReturnValue(t *testing.T) {
	skipPluginTestIfEnvNotSet(t)

	input := `
	p = import "github.com/goby-lang/goby/test_fixtures/import_test/plugin/plugin.go"
	p.send("Bar")
	`

	v := initTestVM()
	evaluated := v.testEval(t, input)
	checkExpected(t, 0, evaluated, "Bar")
	v.checkCFP(t, 0, 0)
	v.checkSP(t, 0, 1)
}

func skipPluginTestIfEnvNotSet(t *testing.T) {
	if os.Getenv("TEST_PLUGIN") == "" {
		t.Skip("skipping plugin related tests")
	}
}