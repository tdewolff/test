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

// Epsilon is used for floating point comparison.
var Epsilon = 1e-10

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

func message(msgs ...any) string {
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
			if r <= 0xFF {
				s2 += fmt.Sprintf("\\x%02X", r)
			} else if r <= 0xFFFF {
				s2 += fmt.Sprintf("\\u%04X", r)
			} else {
				s2 += fmt.Sprintf("\\U%08X", r)
			}
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

func color(color string, s any) string {
	return fmt.Sprintf("\033[00;%sm%v\033[00m", color, s)
}

func floatEqual(a, b, epsilon float64) bool {
	// use mix of relative and absolute difference for large and small numbers respectively
	// see: https://stackoverflow.com/a/32334103
	if a == b {
		return true
	}
	diff := math.Abs(a - b)
	norm := math.Min(math.Abs(a)+math.Abs(b), math.MaxFloat64)
	return diff < epsilon*math.Max(1.0, norm)
}

/////////////////////////////////////////////////////////////////

func Fail(t *testing.T, msgs ...any) {
	t.Helper()
	t.Fatalf("%s%s", trace(), message(msgs...))
}

func Error(t *testing.T, err error, msgs ...any) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s%s: %s", trace(), message(msgs...), color(Red, err.Error()))
	}
}

func That(t *testing.T, condition bool, msgs ...any) {
	t.Helper()
	if !condition {
		t.Fatalf("%s%s: false", trace(), message(msgs...))
	}
}

func equalsInterface(got, wanted reflect.Value) (bool, bool) {
	if equals, ok := wanted.Type().MethodByName("Equals"); ok && equals.Type.NumIn() == 2 && equals.Type.NumOut() == 1 && equals.Type.In(0) == wanted.Type() && equals.Type.In(1) == got.Type() && equals.Type.Out(0).Kind() == reflect.Bool {
		return equals.Func.Call([]reflect.Value{wanted, got})[0].Bool(), true
	}
	return false, false
}

func T(t *testing.T, got, wanted any, msgs ...any) {
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
		gotValue := reflect.ValueOf(got)
		wantedValue := reflect.ValueOf(wanted)
		if equals, ok := equalsInterface(gotValue, wantedValue); ok && equals {
			return
		} else if wantedValue.Kind() == reflect.Slice {
			if gotValue.Len() == wantedValue.Len() {
				i := 0
				for ; i < wantedValue.Len(); i++ {
					if equals, ok := equalsInterface(gotValue.Index(i), wantedValue.Index(i)); !ok {
						break
					} else if !equals {
						break
					}
				}
				if i == wantedValue.Len() {
					return
				}
			}
		}
	}
	t.Fatalf("%s%s: %v != %v", trace(), message(msgs...), color(Red, got), color(Green, wanted))
}

func Bytes(t *testing.T, got, wanted []byte, msgs ...any) {
	t.Helper()
	if !bytes.Equal(got, wanted) {
		gotString := printable(string(got))
		wantedString := printable(string(wanted))
		t.Fatalf("%s%s:\n%s\n%s", trace(), message(msgs...), color(Red, gotString), color(Green, wantedString))
	}
}

func String(t *testing.T, got, wanted string, msgs ...any) {
	t.Helper()
	if got != wanted {
		gotString := printable(got)
		wantedString := printable(wanted)
		t.Fatalf("%s%s:\n%s\n%s", trace(), message(msgs...), color(Red, gotString), color(Green, wantedString))
	}
}

func Float(t *testing.T, got, wanted float64, msgs ...any) {
	t.Helper()
	if math.IsNaN(wanted) != math.IsNaN(got) || !math.IsNaN(wanted) && !floatEqual(got, wanted, Epsilon) {
		t.Fatalf("%s%s: %v != %v", trace(), message(msgs...), color(Red, got), color(Green, wanted))
	}
}

func Floats(t *testing.T, got, wanted []float64, msgs ...any) {
	t.Helper()
	equal := len(got) == len(wanted)
	if equal {
		for i := range got {
			if math.IsNaN(wanted[i]) != math.IsNaN(got[i]) || !math.IsNaN(wanted[i]) && !floatEqual(got[i], wanted[i], Epsilon) {
				equal = false
				break
			}
		}
	}
	if !equal {
		t.Fatalf("%s%s: %v != %v", trace(), message(msgs...), color(Red, got), color(Green, wanted))
	}
}

func FloatDiff(t *testing.T, got, wanted, epsilon float64, msgs ...any) {
	t.Helper()
	if math.IsNaN(wanted) != math.IsNaN(got) || !math.IsNaN(wanted) && !floatEqual(got, wanted, epsilon) {
		t.Fatalf("%s%s: %v != %v", trace(), message(msgs...), color(Red, got), color(Green, fmt.Sprintf("%v ± %v", wanted, epsilon)))
	}
}

func Minify(t *testing.T, input string, err error, got, wanted string, msgs ...any) {
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
