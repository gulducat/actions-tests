//go:build integration

// run a mock http server instead of real products

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

func TestIntegration(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200") // CLI default is https, and -dev doesn't do tls

	// run http server
	http.HandleFunc("/", mockHandler)
	go func() {
		log.Fatal(http.ListenAndServe("127.0.0.1:8200", nil))
	}()

	// GET it in Go
	resp, err := http.Get("http://127.0.0.1:8200")
	if err != nil {
		log.Fatal(err)
	}
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("out:", string(bts))

	// fake system CLI call
	fp := NewFakePATH(t)
	curlActually := "#!/bin/bash\ncurl -sSH \"Command: $(basename $0) $*\" localhost:8200"
	fp.AddScript("consul", curlActually)
	fp.AddScript("nomad", curlActually)
	fp.AddScript("vault", curlActually)
	out, err := exec.Command("actions-tests", "all").CombinedOutput()
	if err != nil {
		t.Error("hcdiag err:", err)
	}
	t.Log("hcdiag out:", string(out))
}

func mockHandler(w http.ResponseWriter, r *http.Request) {
	header := r.Header["Command"]
	cmd := ""
	if len(header) > 0 {
		cmd = header[0]
	}
	msg := fmt.Sprintf(`{"uri": "%s", "cmd":"%s"}`, r.RequestURI, cmd)
	log.Println(msg)
	w.Write([]byte(msg))
}

// for faking executables
type fakePATH struct {
	t   *testing.T
	dir string
}

func NewFakePATH(t *testing.T) *fakePATH {
	binDir := path.Join(t.TempDir(), "bin")
	curPath := os.Getenv("PATH")
	delim := ":"
	if runtime.GOOS == "windows" {
		delim = ";"
	}
	t.Setenv("PATH", binDir+delim+curPath)
	if err := os.Mkdir(binDir, 0755); err != nil {
		t.Fatal(err)
	}
	return &fakePATH{t: t, dir: binDir}
}

func (fp *fakePATH) AddScript(name, content string) {
	fPath := filepath.Join(fp.dir, name)
	f, err := os.OpenFile(fPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fp.t.Fatal(err)
	}
	_, err = f.Write([]byte(content))
	if err != nil {
		fp.t.Fatal(err)
	}
}
