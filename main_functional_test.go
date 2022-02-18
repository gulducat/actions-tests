//go:build functional

package main

import (
	"os/exec"
	"testing"

	gorun "github.com/gulducat/go-run-programs/hcl"
)

func TestFunctional(t *testing.T) {
	// start programs
	stop, err := gorun.RunFromHCL("programs.hcl")
	defer stop()
	if err != nil {
		t.Fatal(err)
	}

	args := []string{
		"all", "consul", "nomad", "vault",
	}

	for _, a := range args {
		t.Run(a, func(t *testing.T) {
			out := runBinary(t, a)
			pretendToTest(t, out)
		})
	}
}

func runBinary(t *testing.T, args ...string) string {
	t.Log("running actions-tests with args:", args)
	bts, err := exec.Command("actions-tests", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("error running actions-tests: %s: %s", err, string(bts))
		return ""
	}
	return string(bts)
}

func pretendToTest(t *testing.T, output string) {
	t.Log("hello i am totally a test haha, got this output:\n>>>", output)
}
