package test

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"runtime"
	"strings"
	"testing"
)

// ErrPlain is the default error that is returned for functions in this package.
var ErrPlain = errors.New("error")

////////////////////////////////////////////////////////////////

func fileline(i int) string {
	_, file, line, ok := runtime.Caller(i)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}

func trace() string {
	trace2 := fileline(2)
	trace3 := fileline(3)
	return "\r\t" + strings.Repeat(" ", len(fmt.Sprintf("%s:", trace2))) + "\r\t" + trace3
}

func message(msgs ...interface{}) string {
	if len(msgs) == 0 {
		return "\n"
	}
	return ": " + fmt.Sprintln(msgs...)
}

func printable(s string) string {
	s = strings.Replace(s, "\n", `\n`, -1)
	s = strings.Replace(s, "\r", `\r`, -1)
	s = strings.Replace(s, "\t", `\t`, -1)
	return s
}

/////////////////////////////////////////////////////////////////

func That(t *testing.T, condition bool, msgs ...interface{}) {
	if !condition {
		t.Errorf("%s%s", trace(), message(msgs...))
	}
}

func Error(t *testing.T, err, expected error, msgs ...interface{}) {
	if err != expected {
		t.Errorf("%s%s   error: %v\nexpected: %v\n", trace(), message(msgs...), err, expected)
	}
}

func Int(t *testing.T, output, expected int, msgs ...interface{}) {
	if output != expected {
		t.Errorf("%s%s  output: %d\nexpected: %d\n", trace(), message(msgs...), output, expected)
	}
}

func Float(t *testing.T, output, expected float64, msgs ...interface{}) {
	if math.Abs(output-expected) > 1e-10 {
		t.Errorf("%s%s  output: %f\nexpected: %f\n", trace(), message(msgs...), output, expected)
	}
}

func String(t *testing.T, output, expected string, msgs ...interface{}) {
	if output != expected {
		t.Errorf("%s%s  output: %s\nexpected: %s\n", trace(), message(msgs...), printable(output), printable(expected))
	}
}

func Bytes(t *testing.T, output, expected []byte, msgs ...interface{}) {
	if !bytes.Equal(output, expected) {
		t.Errorf("%s%s  output: %s\nexpected: %s\n", trace(), message(msgs...), printable(string(output)), printable(string(expected)))
	}
}

func Minify(t *testing.T, input string, err error, output, expected string, msgs ...interface{}) {
	if err != nil {
		t.Errorf("%s%s   given: %s\n   error: %v\n", trace(), message(msgs...), printable(input), err)
	}
	if output != expected {
		t.Errorf("%s%s   given: %s\n  output: %s\nexpected: %s\n", trace(), message(msgs...), printable(input), printable(output), printable(expected))
	}
}
