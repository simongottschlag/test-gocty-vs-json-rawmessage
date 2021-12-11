/*
 * Original code by Kevin Gillette (https://github.com/extemporalgenome): https://go.dev/play/p/U_9ejUmD4QJ
 */

package main

import (
	"encoding/json"
	"fmt"
	"sort"
)

type Map map[string]json.RawMessage

func RawCheck(got, want Map) error {
	return rawCheckMap(got, want)
}

func rawCheckValue(got, want json.RawMessage) error {
	if len(got) == 0 {
		return fmt.Errorf("got was of length 0")
	}

	if len(want) == 0 {
		return fmt.Errorf("want was of length 0")
	}

	gotType := getRawType(got)
	wantType := getRawType(want)

	if gotType != wantType {
		return fmt.Errorf("mismatched types (got %s, want %s)", gotType, wantType)
	}

	switch wantType {
	case RawTypeNull, RawTypeTrue, RawTypeFalse:
		return fmt.Errorf("%s should already have been checked", wantType)
	case RawTypeNumber:
		return rawCheckNumber(got, want)
	case RawTypeString:
		return rawCheckString(got, want)
	case RawTypeArray:
		return rawCheckArray(got, want)
	case RawTypeObject:
		return rawCheckObject(got, want)
	case RawTypeUnknown:
		return fmt.Errorf("unknown type: %s", wantType)
	default:
		return fmt.Errorf("unknown type: %s", wantType)
	}
}

func rawCheckNumber(got, want json.RawMessage) error {
	var gotN, wantN json.Number

	err := rawDecodePair(got, want, &gotN, &wantN)
	if err != nil {
		return err
	}

	if gotN.String() == wantN.String() {
		// they happen to use the same representation
		return nil
	}

	gotI, err1 := gotN.Int64()
	wantI, err2 := wantN.Int64()
	if err1 == nil && err2 == nil {
		if gotI != wantI {
			return fmt.Errorf("mismatched number values (got %d, want %d)", gotI, wantI)
		}

		return nil
	}

	gotF, err1 := gotN.Float64()
	wantF, err2 := wantN.Float64()
	if err1 == nil && err2 == nil {
		if gotF != wantF {
			return fmt.Errorf("mismatched number values (got %g, want %g)", gotF, wantF)
		}

		return nil
	}

	// unknown issue (perhaps input too large to represent)
	return fmt.Errorf("mismatched number values (got %s, want %s)", got, want)
}

func rawCheckString(got, want json.RawMessage) error {
	var gotS, wantS string

	// decode to ensure alternate encodings of the same text
	// are properly accounted for.
	err := rawDecodePair(got, want, &gotS, &wantS)
	if err != nil {
		return err
	}

	if gotS != wantS {
		return fmt.Errorf("mismatched strings (got %q, want %q)", got, want)
	}

	return nil
}

func rawCheckArray(got, want json.RawMessage) error {
	var gotA, wantA []json.RawMessage

	err := rawDecodePair(got, want, &gotA, &wantA)
	if err != nil {
		return err
	}

	for i := range wantA {
		err := rawCheckArrayContains(gotA, wantA[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func rawCheckArrayContains(gotA []json.RawMessage, want json.RawMessage) error {
	for i := range gotA {
		err := rawCheckValue(gotA[i], want)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("could not find want %q in gotA %q", want, gotA)
}

func rawCheckObject(got, want json.RawMessage) error {
	var gotM, wantM Map

	rawDecodePair(got, want, &gotM, &wantM)
	return rawCheckMap(gotM, wantM)
}

func rawCheckMap(got, want Map) error {
	keys := make([]string, 0, len(want))
	for key := range want {
		keys = append(keys, key)
	}

	// ensure deterministic message output
	sort.Strings(keys)

	for _, key := range keys {
		want := want[key]

		got, ok := got[key]
		if !ok {
			return fmt.Errorf("missing key %q", key)
		}

		err := rawCheckValue(got, want)
		if err != nil {
			return err
		}
	}

	return nil
}

func rawDecodePair(got, want []byte, gotDst, wantDst interface{}) error {
	err := rawDecode("got", got, gotDst)
	if err != nil {
		return err
	}

	err = rawDecode("want", want, wantDst)
	if err != nil {
		return err
	}

	return nil
}

func rawDecode(source string, data []byte, dst interface{}) error {
	err := json.Unmarshal(data, dst)
	if err != nil {
		return fmt.Errorf("invalid json received in rawDecode: %w", err)
	}

	return nil
}

type RawType string

const (
	RawTypeNull    RawType = "null"
	RawTypeFalse   RawType = "false"
	RawTypeTrue    RawType = "true"
	RawTypeNumber  RawType = "number"
	RawTypeString  RawType = "string"
	RawTypeArray   RawType = "array"
	RawTypeObject  RawType = "object"
	RawTypeUnknown RawType = "unknown"
)

func getRawType(v json.RawMessage) RawType {
	if len(v) == 0 {
		return RawTypeUnknown
	}

	switch v[0] {
	case 'n':
		return RawTypeNull
	case 'f':
		return RawTypeFalse
	case 't':
		return RawTypeTrue
	case '"':
		return RawTypeString
	case '[':
		return RawTypeArray
	case '{':
		return RawTypeObject
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return RawTypeNumber
	}

	return RawTypeUnknown
}
