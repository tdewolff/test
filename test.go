package test

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unicode"
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
	return "\r    " + strings.Repeat(" ", len(fmt.Sprintf("%s:", trace2))) + "\r    " + trace3
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
	s = strings.Replace(s, "\x00", `\0`, -1)

	s2 := ""
	for _, r := range s {
		if !unicode.IsPrint(r) {
			s2 += fmt.Sprintf("\\x%X", r)
		} else {
			s2 += string(r)
		}
	}
	return s2
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
	t.Helper()
	t.Fatalf("%s%s", trace(), message(msgs...))
}

func Error(t *testing.T, err error, msgs ...interface{}) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s%s: %s", trace(), message(msgs...), color(Red, err.Error()))
	}
}

func That(t *testing.T, condition bool, msgs ...interface{}) {
	t.Helper()
	if !condition {
		t.Fatalf("%s%s: false", trace(), message(msgs...))
	}
}

func T(t *testing.T, got, wanted interface{}, msgs ...interface{}) {
	t.Helper()
	gotType := reflect.TypeOf(got)
	wantedType := reflect.TypeOf(wanted)
	if gotType != wantedType {
		t.Fatalf("%s%s: type %v != %v", trace(), message(msgs...), color(Red, gotType), color(Green, wantedType))
		return
	}
	if reflect.DeepEqual(got, wanted) {
		return
	}

	if wantedType != nil {
		if equals, ok := wantedType.MethodByName("Equals"); ok && equals.Type.NumIn() == 2 && equals.Type.NumOut() == 1 && equals.Type.In(0) == wantedType && equals.Type.In(1) == gotType && equals.Type.Out(0).Kind() == reflect.Bool && equals.Func.Call([]reflect.Value{reflect.ValueOf(wanted), reflect.ValueOf(got)})[0].Bool() {
			return
		}
	}
	t.Fatalf("%s%s: %v != %v", trace(), message(msgs...), color(Red, got), color(Green, wanted))
}

func Bytes(t *testing.T, got, wanted []byte, msgs ...interface{}) {
	t.Helper()
	if !bytes.Equal(got, wanted) {
		gotString := printable(string(got))
		wantedString := printable(string(wanted))
		t.Fatalf("%s%s:\n%s\n%s", trace(), message(msgs...), color(Red, gotString), color(Green, wantedString))
	}
}

func String(t *testing.T, got, wanted string, msgs ...interface{}) {
	t.Helper()
	if got != wanted {
		gotString := printable(got)
		wantedString := printable(wanted)
		t.Fatalf("%s%s:\n%s\n%s", trace(), message(msgs...), color(Red, gotString), color(Green, wantedString))
	}
}

func Float(t *testing.T, got, wanted float64, msgs ...interface{}) {
	t.Helper()
	if math.IsNaN(wanted) != math.IsNaN(got) || !math.IsNaN(wanted) && math.Abs(got-wanted) > 1e-6 {
		t.Fatalf("%s%s: %v != %v", trace(), message(msgs...), color(Red, got), color(Green, wanted))
	}
}

func FloatDiff(t *testing.T, got, wanted, diff float64, msgs ...interface{}) {
	t.Helper()
	if math.IsNaN(wanted) != math.IsNaN(got) || !math.IsNaN(wanted) && math.Abs(got-wanted) > diff {
		t.Fatalf("%s%s: %v != %v", trace(), message(msgs...), color(Red, got), color(Green, fmt.Sprintf("%v ± %v", wanted, diff)))
	}
}

func Minify(t *testing.T, input string, err error, got, wanted string, msgs ...interface{}) {
	t.Helper()
	inputString := printable(input)
	if err != nil {
		t.Fatalf("%s%s:\n%s\n%s", trace(), message(msgs...), inputString, color(Red, err.Error()))
		return
	}

	if got != wanted {
		gotString := printable(got)
		wantedString := printable(wanted)
		t.Fatalf("%s%s:\n%s\n%s\n%s", trace(), message(msgs...), inputString, color(Red, gotString), color(Green, wantedString))
	}
}
