/*
 * Code by Kevin Gillette (https://github.com/extemporalgenome): https://go.dev/play/p/U_9ejUmD4QJ
 */

package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
)

func Example() {
	want := Map{
		"arr_fail": b(`[{"x":[null]}]`),
		"arr_ok":   b(`[1.0,["\/",{"x":[]}]]`),
		"f_fail":   b(`false`),
		"f_ok":     b(`false`),
		"n_fail":   b(`null`),
		"n_ok":     b(`null`),
		"num_fail": b(`5.1`),
		"num_ok":   b(`5.3e1`),
		"obj_fail": b(`{"missing":[],"fail":["\/",null]}`),
		"obj_ok":   b(`{"ok":[]}`),
		"str_fail": b(`"pqr"`),
		"str_ok":   b(`"\/"`), // json forward slashes can be optionally escaped
		"t_fail":   b(`true`),
		"t_ok":     b(`true`),
	}

	got := Map{
		"arr_fail": b(`[{"x":[[]]},true]`),
		"arr_ok":   b(`[10e-1,["/",{"x":["extra"]}]]`),
		"f_fail":   b(`true`),
		"f_ok":     b(`false`),
		"n_fail":   b(`false`),
		"n_ok":     b(`null`),
		"num_fail": b(`5e1`),
		"num_ok":   b(`53`), // different format of same value
		"obj_fail": b(`{"fail":[0]}`),
		"obj_ok":   b(`{"ok":["extra"],"bonus":"field"}`),
		"str_fail": b(`"xyz"`),
		"str_ok":   b(`"/"`), // different format of same value
		"t_fail":   b(`100`),
		"t_ok":     b(`true`),
		"extra":    b(`"doesn't matter"`),
	}

	msgs := Check(got, want)
	if len(msgs) == 0 {
		fmt.Println("complete validation success!")
		return
	}

	fmt.Println("validation errors:")
	for _, fm := range msgs {
		fmt.Printf("- %s: %s\n", fm.Field, fm.Message)
	}

	// Output:
	// validation errors:
	// - arr_fail.0.x.0: mismatched types (got array, want null)
	// - f_fail: mismatched types (got true, want false)
	// - n_fail: mismatched types (got false, want null)
	// - num_fail: mismatched number values (got 50, want 5.1)
	// - obj_fail.fail: missing elements (got 1, want 2)
	// - obj_fail.fail.0: mismatched types (got number, want string)
	// - obj_fail.missing: missing key
	// - str_fail: mismatched strings (got "xyz", want "pqr")
	// - t_fail: mismatched types (got number, want true)
}

func b(s string) []byte {
	return []byte(s)
}

type Map map[string]json.RawMessage

type FieldMessage struct{ Field, Message string }

func (f FieldMessage) String() string {
	return f.Field + ": " + f.Message
}

// Check ensures that got has at least all keys in want,
// and that their values are compatible for a given key:
//
// - Scalars must be semantically equal.
// - Objects in got must be a compatible superset of those in want.
// - Arrays in got must have at least as many elements as want
//   and be element-wise compatible:
//   it's an exercise to the reader to implement the pretty expensive
//   "got must be a compatible superset of elements in want,
//   regardless of order."
//
// If validation succeeds completely, the returned slice will be nil;
// otherwise, it will contain validation failure messages.
//
func Check(got, want Map) []FieldMessage {
	var msgs []FieldMessage
	ff := newFailFormatter(&msgs)
	checkMap(ff, got, want)

	return msgs
}

type Type string

const (
	TypeNull   Type = "null"
	TypeFalse  Type = "false"
	TypeTrue   Type = "true"
	TypeNumber Type = "number"
	TypeString Type = "string"
	TypeArray  Type = "array"
	TypeObject Type = "object"
)

func ParseType(leading byte) Type {
	switch leading {
	case 'n':
		return TypeNull
	case 'f':
		return TypeFalse
	case 't':
		return TypeTrue
	case '"':
		return TypeString
	case '[':
		return TypeArray
	case '{':
		return TypeObject
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return TypeNumber
	}

	msg := "non-json byte given: %q (was json.Unmarshal/Decoder used?)"
	panic(fmt.Errorf(msg, leading))
}

func checkValue(ff failFormatter, got, want json.RawMessage) {
	if len(got) == 0 {
		// this shouldn't happen if the Map is decoded through the json package.
		msg := "programmer bug: empty RawMessage for got.%s"
		panic(fmt.Errorf(msg, ff))
	}

	if len(want) == 0 {
		msg := "programmer bug: empty RawMessage for want.%s"
		panic(fmt.Errorf(msg, ff))
	}

	gotType := ParseType(got[0])
	wantType := ParseType(want[0])

	if gotType != wantType {
		msg := "mismatched types (got %s, want %s)"
		ff.Fail(msg, gotType, wantType)
		return
	}

	switch wantType {
	case TypeNull, TypeTrue, TypeFalse:
		// already handled: we're treating value checking as a type
		// compatibility problem (true and false types instead of bool).
		return
	case TypeNumber:
		checkNumber(ff, got, want)
	case TypeString:
		checkString(ff, got, want)
	case TypeArray:
		checkArray(ff, got, want)
	case TypeObject:
		checkObject(ff, got, want)
	}
}

func checkNumber(ff failFormatter, got, want json.RawMessage) {
	var gotN, wantN json.Number

	decodePair(ff, got, want, &gotN, &wantN)

	if gotN.String() == wantN.String() {
		// they happen to use the same representation
		return
	}

	gotI, err1 := gotN.Int64()
	wantI, err2 := wantN.Int64()
	if err1 == nil && err2 == nil {
		if gotI != wantI {
			msg := "mismatched number values (got %d, want %d)"
			ff.Fail(msg, gotI, wantI)
		}

		return
	}

	gotF, err1 := gotN.Float64()
	wantF, err2 := wantN.Float64()
	if err1 == nil && err2 == nil {
		if gotF != wantF {
			msg := "mismatched number values (got %g, want %g)"
			ff.Fail(msg, gotF, wantF)
		}

		return
	}

	// unknown issue (perhaps input too large to represent)
	msg := "mismatched number values (got %s, want %s)"
	ff.Fail(msg, got, want)
}

func checkString(ff failFormatter, got, want json.RawMessage) {
	var gotS, wantS string

	// decode to ensure alternate encodings of the same text
	// are properly accounted for.
	decodePair(ff, got, want, &gotS, &wantS)

	if gotS != wantS {
		msg := "mismatched strings (got %q, want %q)"
		ff.Fail(msg, gotS, wantS)
	}
}

func checkArray(ff failFormatter, got, want json.RawMessage) {
	var gotA, wantA []json.RawMessage

	decodePair(ff, got, want, &gotA, &wantA)

	gotN := len(gotA)
	wantN := len(wantA)

	// even if lengths mismatch, check as much as possible
	n := wantN
	if gotN < wantN {
		n = gotN
		msg := "missing elements (got %d, want %d)"
		ff.Fail(msg, gotN, wantN)
	}

	for i := range wantA[:n] {
		got = gotA[i]
		want = wantA[i]

		ef := ff.Field(strconv.Itoa(i))
		checkValue(ef, got, want)
	}
}

func checkObject(ff failFormatter, got, want json.RawMessage) {
	var gotM, wantM Map

	decodePair(ff, got, want, &gotM, &wantM)
	checkMap(ff, gotM, wantM)
}

func checkMap(ff failFormatter, got, want Map) {
	keys := make([]string, 0, len(want))
	for key := range want {
		keys = append(keys, key)
	}

	// ensure deterministic message output
	sort.Strings(keys)

	for _, key := range keys {
		kf := ff.Field(key)

		want := want[key]

		got, ok := got[key]
		if !ok {
			kf.Fail("missing key")
			continue
		}

		checkValue(kf, got, want)
	}
}

func decodePair(ff failFormatter, got, want []byte, gotDst, wantDst interface{}) {
	decode(ff, "got", got, gotDst)
	decode(ff, "want", want, wantDst)
}

func decode(ff failFormatter, source string, data []byte, dst interface{}) {
	err := json.Unmarshal(data, dst)
	if err != nil {
		msg := "programmer bug: invalid json for %s.%s (did you use json.Unmarshal/Decoder?)"
		panic(fmt.Errorf(msg, source, ff))
	}
}

func newFailFormatter(dst *[]FieldMessage) failFormatter {
	return failFormatter{dst: dst}
}

type failFormatter struct {
	dst   *[]FieldMessage
	field string
}

func (f failFormatter) String() string {
	return f.field
}

// Field returns a nested formatted, such that
// newFailFormatter(...).Field("x").Field("0").Field("y")
// yields a formatter with field "x.0.y"
//
// Field will panic if name is empty.
func (f failFormatter) Field(name string) failFormatter {
	if f.field == "" {
		f.field = name
	} else {
		f.field += "." + name
	}

	return f
}

// Fail records a message in the closed-over slice.
//
// Fail will panic if this failFormatter was not derived
// through a call to Field.
func (f failFormatter) Fail(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	*f.dst = append(*f.dst, FieldMessage{Field: f.field, Message: msg})
}
