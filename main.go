package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	os.Setenv("VAULT_ADDR", "http://127.0.0.1:8200") // no https b/c -dev

	// run programs in the background
	ctx, stop := context.WithCancel(context.Background())
	defer stop()
	go exec.CommandContext(ctx, "consul", "agent", "-dev").Run()
	go exec.CommandContext(ctx, "nomad", "agent", "-dev").Run()
	go exec.CommandContext(ctx, "vault", "server", "-dev").Run()
	// lazy sleep for programs to start, might not be long enough for nomad..
	time.Sleep(time.Second * 8)

	// test em out
	exitCode := 0
	for _, c := range [][]string{
		{"consul", "members"},
		{"nomad", "node", "status"},
		{"vault", "status"},
	} {
		fmt.Println("running", c)
		bts, err := exec.Command(c[0], c[1:]...).CombinedOutput()
		if err != nil {
			exitCode++
			fmt.Println("err:", err)
		}
		fmt.Println(string(bts))
	}
	os.Exit(exitCode)
}
