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
