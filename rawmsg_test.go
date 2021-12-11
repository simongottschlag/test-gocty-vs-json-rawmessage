package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRawCheck(t *testing.T) {
	rawCases := testToTestRawCases(t, testCases)

	for i, c := range rawCases {
		t.Logf("Test iteration %d: %s", i, c.testDescription)

		err := RawCheck(c.got, c.want)
		if c.expectedResult {
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
