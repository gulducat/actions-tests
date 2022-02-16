package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("Hi there:", os.Args)
	cmd := os.Args[1]
	cmdMap := map[string][]string{
		"consul": {"members"},
		"nomad":  {"node", "status"},
		"vault":  {"status"},
	}
	for c, args := range cmdMap {
		if c == cmd || cmd == "all" {
			runCommand(c, args...)
		}
	}
}

func runCommand(c string, args ...string) {
	bts, err := exec.Command(c, args...).CombinedOutput()
	if err != nil {
		fmt.Println(c, "eeeee:", err)
	}
	fmt.Println(string(bts))
}

//func mainBlah() {
//	os.Setenv("VAULT_ADDR", "http://127.0.0.1:8200") // no https b/c -dev
//
//	// run programs in the background
//	// ctx, stop := context.WithCancel(context.Background())
//	// defer stop()
//	// // go exec.CommandContext(ctx, "consul", "agent", "-dev").Run()
//	// // go exec.CommandContext(ctx, "nomad", "agent", "-dev").Run()
//	// // go exec.CommandContext(ctx, "vault", "server", "-dev").Run()
//	// go runWithContext(ctx, stop, "consul", "agent", "-dev")
//	// go runWithContext(ctx, stop, "nomad", "agent", "-dev")
//	// go runWithContext(ctx, stop, "vault", "server", "-dev")
//	// // lazy sleep for programs to start, might not be long enough for nomad..
//	// time.Sleep(time.Second * 8)
//	// test em out
//	exitCode := 0
//	for _, c := range [][]string{
//		{"consul", "members"},
//		{"nomad", "node", "status"},
//		{"vault", "status"},
//	} {
//		fmt.Println("running", c)
//		bts, err := exec.Command(c[0], c[1:]...).CombinedOutput()
//		if err != nil {
//			exitCode++
//			fmt.Println("err:", err)
//		}
//		fmt.Println(string(bts))
//	}
//	time.Sleep(time.Second * 10)
//	os.Exit(exitCode)
//}

// func runWithContext(ctx context.Context, stop context.CancelFunc, command ...string) {
// 	out, err := exec.CommandContext(ctx, command[0], command[1:]...).CombinedOutput()
// 	if err != nil {
// 		fmt.Printf("XXXXX err running %s:\n  XXXoutput: %s\n  XXXerror: %s\nXXXXX\n", command, out, err)
// 		stop()
// 	}
// }

// path := path.Join(os.Getenv("GOPATH"), "bin")
// if err := download(c[0], "latest", path); err != nil {
// 	fmt.Println("uh oh:", err)
// }
// continue

// func download(product, version, dir string) error {
// 	if dir == "" {
// 		dir = "."
// 	}
// 	fmt.Println("downloading", product, version, "to", dir)
// 	releasesURL := hbVar.ReleasesURL + "/" + product + "/index.json"
// 	fmt.Println(releasesURL)
// 	p, err := hb.NewProduct(releasesURL)
// 	if err != nil {
// 		fmt.Println("bad product", err)
// 		return err
// 	}
// 	v, err := p.GetVersion(version)
// 	if err != nil {
// 		fmt.Println("bad version", err)
// 		return err
// 	}
// 	b := v.GetBuildForLocal()
// 	// fmt.Println(p, v, b)
// 	if _, err := b.DownloadAndExtract(dir, product); err != nil {
// 		fmt.Println("shit install", err)
// 		return err
// 	}
// 	return nil
// }
