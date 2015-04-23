package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// Call

func TestCall(t *testing.T) {
	Call("/bin/echo", "foo")
}

func TestCallFailure(t *testing.T) {
	wantPanic(t, func() {Call("/bin/sh", "/bar")})
}

func TestCallSh(t *testing.T) {
	CallSh("/bin/echo \"foo\" | echo ")
}

func TestCallShFailure(t *testing.T) {
	wantPanic(t, func() { CallSh("/bin/sh baz") })
}


// CallStatus

func TestCallStatus(t *testing.T) {
	wantNoError(t, CallStatus("/bin/echo", "foo"))
}

func TestCallStatusFailure(t *testing.T) {
	wantError(t, CallStatus("/bin/sh", "baz"))
}

func TestCallStatusSh(t *testing.T) {
	wantNoError(t, CallStatusSh("/bin/echo foo"))
}

func TestCallStatusShFailure(t *testing.T) {
	wantError(t, CallStatusSh("/bin/sh baz"))
}


// Output

func TestOutput(t *testing.T) {
	wantEqual(t, "foo\n", Output("/bin/echo", "foo"))
}

func TestOutputFailure(t *testing.T) {
	wantPanic(t, func() { Output("/bin/sh", "foo") })
}

func TestOutputSh(t *testing.T) {
	wantEqual(t, "foo\n", OutputSh("/bin/echo \"echo foo\" | /bin/sh"))
}

func TestOutputShFailure(t *testing.T) {
	wantPanic(t, func() { OutputSh("/bin/sh baz") })
}


// OutputStatus

func TestOutputStatus(t *testing.T) {
	out, st := OutputStatus("/bin/echo", "foo")
	wantEqual(t, "foo\n", out)
	wantNoError(t, st)
}

func TestOutputStatusFailure(t *testing.T) {
	out, st := OutputStatus("/bin/sh", "foo")
	wantEqual(t, "", out)
	wantError(t, st)
}

func TestOutputStatusSh(t *testing.T) {
	out, st := OutputStatusSh("/bin/echo \"echo foo\" | /bin/sh")
	wantEqual(t, "foo\n", out)
	wantNoError(t, st)
}

func TestOutputStatusShFailure(t *testing.T) {
	out, st := OutputStatusSh("/bin/echo foo; /bin/sh baz")
	wantEqual(t, "foo\n", out)
	wantError(t, st)
}


// OutputError

func TestOutputError(t *testing.T) {
	f := shellScript(t, "echo \"bar\" >& 2; echo \"foo\"")
	defer f.Close()
	out, err := OutputError("/bin/sh", f.Name())
	wantEqual(t, "foo\n", out)
	wantEqual(t, "bar\n", err)
}

func TestOutputErrorFailure(t *testing.T) {
	f := shellScript(t, "echo \"bar\" >& 2; echo \"foo\"; foo")
	defer f.Close()
	wantPanic(t, func() { OutputError("/bin/sh", f.Name()) })
}

func TestOutputErrorSh(t *testing.T) {
	out, err := OutputErrorSh("echo \"bar\" >& 2; echo \"foo\"")
	wantEqual(t, "foo\n", out)
	wantEqual(t, "bar\n", err)
}

func TestOutputErrorShFailure(t *testing.T) {
	wantPanic(t, func() { OutputErrorSh("echo \"foo\"; echo bar >& 2; foo") })
}

// OutputErrorStatus

func TestOutputErrorStatus(t *testing.T) {
	f := shellScript(t, "echo \"bar\" >& 2; echo \"foo\"")
	defer f.Close()
	out, err, st := OutputErrorStatus("/bin/sh", f.Name())
	wantEqual(t, "foo\n", out)
	wantEqual(t, "bar\n", err)
	wantNoError(t, st)
}

func TestOutputErrorStatusFailure(t *testing.T) {
	f := shellScript(t, "echo \"bar\" >& 2; echo \"foo\"; foo")
	defer f.Close()
	out, err, st := OutputErrorStatus("/bin/sh", f.Name())
	wantEqual(t, "foo\n", out)
	wantContains(t, "bar", err)
	wantContains(t, "not found", err)
	wantError(t, st)
}

func TestOutputErrorStatusSh(t *testing.T) {
	out, err, st := OutputErrorStatusSh("echo \"bar\" >& 2; echo \"foo\"")
	wantEqual(t, "foo\n", out)
	wantEqual(t, "bar\n", err)
	wantNoError(t, st)
}

func TestOutputErrorStatusShFailure(t *testing.T) {
	out, err, st := OutputErrorStatusSh("echo \"foo\"; echo bar >& 2; foo")
	wantEqual(t, "foo\n", out)
	wantContains(t, "bar", err)
	wantContains(t, "not found", err)
	wantError(t, st)
}


// CombinedOutput

func TestCombinedOutput(t *testing.T) {
	f := shellScript(t, "echo \"foo\"; echo bar >& 2")
	defer f.Close()
	out := CombinedOutput("/bin/sh", f.Name())
	wantContains(t, "foo", out)
	wantContains(t, "bar", out)
}

func TestCombinedOutputFailure(t *testing.T) {
	wantPanic(t, func() { CombinedOutput("/bin/sh", "baz") })
}

func TestCombinedOutputSh(t *testing.T) {
	out := CombinedOutputSh("echo \"foo\"; echo bar >& 2")
	wantContains(t, "foo", out)
	wantContains(t, "bar", out)
}

func TestCombinedOutputShFailure(t *testing.T) {
	wantPanic(t, func() { CombinedOutputSh("baz") })
}


// CombinedOutputStatus

func TestCombinedOutputStatus(t *testing.T) {
	f := shellScript(t, "echo \"foo\"; echo bar >& 2")
	defer f.Close()
	out, st := CombinedOutputStatus("/bin/sh", f.Name())
	wantContains(t, "foo", out)
	wantContains(t, "bar", out)
	wantNoError(t, st)
}

func TestCombinedOutputStatusFailure(t *testing.T) {
	f := shellScript(t, "echo \"foo\"; echo bar >& 2; baz")
	defer f.Close()
	out, st := CombinedOutputStatus("/bin/sh", f.Name())
	wantContains(t, "foo", out)
	wantContains(t, "bar", out)
	wantContains(t, "baz", out)
	wantContains(t, "not found", out)
	wantError(t, st)
}

func TestCombinedOutputStatusSh(t *testing.T) {
	out, st := CombinedOutputStatusSh("echo \"foo\"; echo bar >& 2")
	wantContains(t, "foo", out)
	wantContains(t, "bar", out)
	wantNoError(t, st)
}

func TestCombinedOutputStatusFailureSh(t *testing.T) {
	out, st := CombinedOutputStatusSh("echo \"foo\"; echo bar >& 2; baz")
	wantContains(t, "foo", out)
	wantContains(t, "bar", out)
	wantContains(t, "baz", out)
	wantContains(t, "not found", out)
	wantError(t, st)
}


// Helpers

func shellScript(t *testing.T, cmd string) *os.File {
	f, err := ioutil.TempFile("", "writeShellScript")
	if err != nil {
		t.Fatalf("could not open temp script file: %s", err)
	}
	f.Write([]byte(cmd))
	return f
}

func wantPanic(t *testing.T, f func()) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	f()
}

func wantError(t *testing.T, st int) {
	if st == noErr {
		t.Fatalf("wanted error")
	}
}

func wantNoError(t *testing.T, st int) {
	if st != noErr {
		t.Fatalf("wanted no error")
	}
}

func wantEqual(t *testing.T, want, got string) {
	if want != got {
		t.Fatalf("wanted %q, but got %q", want, got)
	}
}

func wantNotEqual(t *testing.T, want, got string) {
	if want == got {
		t.Fatalf("dit not want %q, but got %q", want, got)
	}
}

func wantContains(t *testing.T, want, got string) {
	if !strings.Contains(got, want) {
		t.Fatalf("wanted to find %q in the string %q", want, got)
	}
}
