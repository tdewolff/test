package test

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
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
		return ""
	}
	s := fmt.Sprintln(msgs...)
	s = s[:len(s)-1] // remove newline
	return ": " + s
}

func printable(s string) string {
	s = strings.Replace(s, "\n", `\n`, -1)
	s = strings.Replace(s, "\r", `\r`, -1)
	s = strings.Replace(s, "\t", `\t`, -1)
	s = strings.Replace(s, " ", "\u00B7", -1)
	return s
}

const (
	Red   = "31"
	Green = "32"
)

func color(color string, s interface{}) string {
	return fmt.Sprintf("\033[00;%sm%v\033[00m", color, s)
}

/////////////////////////////////////////////////////////////////

func Fail(t *testing.T, msgs ...interface{}) {
	t.Errorf("%s%s", trace(), message(msgs...))
}

func That(t *testing.T, condition bool, msgs ...interface{}) {
	if !condition {
		t.Errorf("%s%s", trace(), message(msgs...))
	}
}

func T(t *testing.T, got, wanted interface{}, msgs ...interface{}) {
	gotType := reflect.TypeOf(got)
	wantedType := reflect.TypeOf(wanted)
	if gotType != wantedType {
		t.Errorf("%s%s: type %v != %v", trace(), message(msgs...), color(Red, gotType), color(Green, wantedType))
	}
	if got != wanted {
		t.Errorf("%s%s: %v != %v", trace(), message(msgs...), color(Red, got), color(Green, wanted))
	}
}

func Bytes(t *testing.T, got, wanted []byte, msgs ...interface{}) {
	if !bytes.Equal(got, wanted) {
		gotString := printable(string(got))
		wantedString := printable(string(wanted))
		t.Errorf("%s%s\n%s\n%s", trace(), message(msgs...), color(Red, gotString), color(Green, wantedString))
	}
}

func String(t *testing.T, got, wanted string, msgs ...interface{}) {
	if got != wanted {
		gotString := printable(got)
		wantedString := printable(wanted)
		t.Errorf("%s%s\n%s\n%s", trace(), message(msgs...), color(Red, gotString), color(Green, wantedString))
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
