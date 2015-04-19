package main

import (
	"fmt"
	"os"
)


func main() {
	err := runCmd("/usr/bin/say", []string{"say", "-v", "cello", "hi there"})
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
}


func runCmd(path string, cmd []string) error {
	// Run /bin/hostname and collect output.
	r, w, err := os.Pipe()
	defer r.Close()
	p, err := os.StartProcess(path, cmd, &os.ProcAttr{Files: []*os.File{nil, nil, nil}})
	if err != nil {
		return err
	}
	w.Close()

	_, err = p.Wait()
	if err != nil {
		return err
	}
	err = p.Kill()
	if err == nil {
		return err
	}
	return nil
}

