package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtyCheck(t *testing.T) {
	for i, c := range testCases {
		t.Logf("Test iteration %d: %s", i, c.testDescription)

		err := ctyCheck(c.requiredClaims, c.tokenClaims)

		if c.expectedResult {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

func TestCheck(t *testing.T) {
	rawCases := testToTestRawCases(t, testCases)

	for i, c := range rawCases {
		t.Logf("Test iteration %d: %s", i, c.testDescription)

		msgs := Check(c.got, c.want)

		var err error
		if len(msgs) != 0 {
			err = fmt.Errorf("validation errors")
		}

		if c.expectedResult {
			if len(msgs) != 0 {
				t.Logf("\n\tGot: %s\n\tWant: %s", c.got, c.want)
				for _, fm := range msgs {
					t.Logf("- %s: %s\n", fm.Field, fm.Message)
				}
			}
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

func TestCheckExample(t *testing.T) {
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
	require.Len(t, msgs, 9)
}

func testToTestRawCases(tb testing.TB, cases []testCase) []testRawCase {
	tb.Helper()

	var rawCases []testRawCase

	for _, c := range testCases {
		wantBytes, err := json.Marshal(c.requiredClaims)
		require.NoError(tb, err)

		want := Map{}
		err = json.Unmarshal(wantBytes, &want)
		require.NoError(tb, err)

		gotBytes, err := json.Marshal(c.tokenClaims)
		require.NoError(tb, err)

		got := Map{}
		err = json.Unmarshal(gotBytes, &got)
		require.NoError(tb, err)

		rawCases = append(rawCases, testRawCase{
			testDescription: c.testDescription,
			want:            want,
			got:             got,
			expectedResult:  c.expectedResult,
		})
	}

	return rawCases
}

type testRawCase struct {
	testDescription string
	want            Map
	got             Map
	expectedResult  bool
}
