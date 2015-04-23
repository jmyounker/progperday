package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
)

func main() {
	st := CallStatus("/usr/bin/say", "-v", "cello", "hi there")
	if st != noErr {
		fmt.Printf("error: %d\n", st)
		return
	}
	fmt.Printf("status: %d\n", st)
}

var shell = "/bin/sh"
const noErr = 0


func Call(path string, arg ...string) {
	nonZeroPanic(CallStatus(path, arg...))
}

func CallSh(cmd string) {
	nonZeroPanic(CallStatusSh(cmd))
}

func CallStatus(path string, arg ...string) int {
	return checkExitErr(exec.Command(path, arg...).Run())
}

func CallStatusSh(cmd string) int {
	return checkExitErr(runShell(cmd, exec.Command(shell)))
}

func Output(path string, arg ...string) string {
	out, st := OutputStatus(path, arg...)
	nonZeroPanic(st)
	return out
}

func OutputSh(cmd string) string {
	out, st := OutputStatusSh(cmd)
	nonZeroPanic(st)
	return out
}

func OutputStatus(path string, arg ...string) (string, int) {
	c := exec.Command(path, arg...)
	out, err := c.Output()
	st := checkExitErr(err)
	return string(out), st
}

func OutputStatusSh(cmd string) (string, int) {
	c := exec.Command(shell)

	var b bytes.Buffer;
	c.Stdout = &b

	st := checkExitErr(runShell(cmd, c))

	return string(b.Bytes()), st
}

func OutputError(path string, arg ...string) (string, string) {
	out, err, st := OutputErrorStatus(path, arg...)
	nonZeroPanic(st)
	return out, err
}

func OutputErrorSh(cmd string) (string, string) {
	out, err, st := OutputErrorStatusSh(cmd)
	nonZeroPanic(st)
	return out, err
}

func OutputErrorStatus(path string, arg ...string) (string, string, int) {
	c := exec.Command(path, arg...)

	var b, e bytes.Buffer;
	c.Stdout = &b
	c.Stderr = &e

	st := checkExitErr(c.Run())
	return string(b.Bytes()), string(e.Bytes()), st
}

func OutputErrorStatusSh(cmd string) (string, string, int) {
	c := exec.Command(shell)

	var b, e bytes.Buffer;
	c.Stdout = &b
	c.Stderr = &e

	st := checkExitErr(runShell(cmd, c))
	return string(b.Bytes()), string(e.Bytes()), st
}

func CombinedOutput(path string, arg ...string) string {
	out, st := CombinedOutputStatus(path, arg...)
	nonZeroPanic(st)
	return out
}

func CombinedOutputSh(cmd string) string {
	out, st := CombinedOutputStatusSh(cmd)
	nonZeroPanic(st)
	return out
}

func CombinedOutputStatus(path string, arg ...string) (string, int) {
	c := exec.Command(path, arg...)
	out, err := c.CombinedOutput()
	return string(out), checkExitErr(err)
}

func CombinedOutputStatusSh(cmd string) (string, int) {
	c := exec.Command(shell)

	var b bytes.Buffer;
	c.Stdout = &b
	c.Stderr = &b

	st := checkExitErr(runShell(cmd, c))
	return string(b.Bytes()), st
}

func runShell(cmd string, c *exec.Cmd) error {
	in, err := c.StdinPipe()
	if err != nil {
		return err
	}
	if err := c.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(cmd)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	return c.Wait()
}

func checkExitErr(err error) int {
	if err == nil {
		return 0
	}

	if exit, ok := err.(*exec.ExitError); ok {
		return exit.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	}
	panic(fmt.Sprintf("error: %s", err))
}

func nonZeroPanic(status int) {
	if status != noErr {
		panic(fmt.Sprintf("exit code: %d", status))
	}
}
