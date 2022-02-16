//go:build functional

package main

import (
	"context"
	"fmt"
	"os/exec"
	"testing"
	"time"
)

type testProgram struct {
	Name  string
	Run   []string
	Check []string
	Test  func(*testing.T, string)
}

func (tp testProgram) String() string {
	return tp.Name
}

func TestFunctional(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200") // default is https, and -dev doesn't do tls

	// programs to run for us to test against
	programs := []testProgram{
		{
			Name:  "consul",
			Run:   []string{"consul", "agent", "-dev"},
			Check: []string{"consul", "members"},
			Test:  pretendToTest, // testConsul
		},
		{
			Name:  "nomad",
			Run:   []string{"nomad", "agent", "-dev"},
			Check: []string{"nomad", "node", "status"},
			Test:  pretendToTest, // testNomad
		},
		{
			Name:  "vault",
			Run:   []string{"vault", "server", "-dev"},
			Check: []string{"vault", "status"},
			Test:  pretendToTest, // testVault
		},
	}
	runPrograms(t, programs...)

	t.Run("all", func(t *testing.T) {
		out := runBinary(t, "all")
		pretendToTest(t, out)
	})

	for _, p := range programs {
		t.Run(p.Name, func(t *testing.T) {
			out := runBinary(t, p.Name)
			p.Test(t, out)
		})
	}
}

func runWithContext(ctx context.Context, stop context.CancelFunc, command ...string) {
	out, err := exec.CommandContext(ctx, command[0], command[1:]...).CombinedOutput()
	if err != nil {
		fmt.Printf("XXXXX err running %s:\n  XXXoutput: %s\n  XXXerror: %s\nXXXXX\n", command, out, err)
		stop() // stop everything now
	}
}

func runPrograms(t *testing.T, programs ...testProgram) {
	t.Log("starting:", programs)

	// run programs in the background
	ctx, stop := context.WithCancel(context.Background())
	t.Cleanup(stop) // stop them after test is done
	for _, p := range programs {
		go runWithContext(ctx, stop, p.Run...)
	}

	// wait for programs to start
	var (
		bts []byte
		err error
	)
	for _, p := range programs {
		for x := 0; x < 20; x++ { // TODO: 20 seconds max wait time?
			bts, err = exec.Command(p.Check[0], p.Check[1:]...).CombinedOutput()
			if err == nil {
				break
			}
			time.Sleep(time.Second)
		}
		if err != nil {
			t.Fatalf("error during test setup: '%s' error: %s: %s", p.Check, err, string(bts))
		}
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
