package test

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

// ErrPlain is the default error that is returned for functions in this package.
var ErrPlain = errors.New("error")

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

func printable(s string) string {
	s = strings.Replace(s, "\n", `\n`, -1)
	s = strings.Replace(s, "\r", `\r`, -1)
	s = strings.Replace(s, "\t", `\t`, -1)
	return s
}

func Assert(t *testing.T, condition bool, msgs ...string) {
	if !condition {
		t.Errorf("%s: %s\n", trace(), strings.Join(msgs, " "))
	}
}

func Error(t *testing.T, err, expected error, msgs ...string) {
	if err != expected {
		t.Errorf("%s: %s\n   error: %v\nexpected: %v\n", trace(), strings.Join(msgs, " "), err, expected)
	}
}

func String(t *testing.T, output, expected string, msgs ...string) {
	if output != expected {
		t.Errorf("%s: %s\nminified: %s\nexpected: %s\n", trace(), strings.Join(msgs, " "), printable(output), printable(expected))
	}
}

func Minify(t *testing.T, input string, err error, output, expected string, msgs ...string) {
	if err != nil {
		t.Errorf("%s: %s\n   given: %s\n   error: %v\n", trace(), strings.Join(msgs, " "), printable(input), err)
	}
	if output != expected {
		t.Errorf("%s: %s\n   given: %s\nminified: %s\nexpected: %s\n", trace(), strings.Join(msgs, " "), printable(input), printable(output), printable(expected))
	}
}
