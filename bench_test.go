package main

import (
	"testing"
)

func BenchmarkCtyCheck(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		for _, c := range testCases {
			// b.Logf("Test iteration %d: %s", i, c.testDescription)
			err := ctyCheck(c.requiredClaims, c.tokenClaims)

			if c.expectedResult {
				if err != nil {
					b.FailNow()
				}
			} else {
				if err == nil {
					b.FailNow()
				}
			}
		}
	}
}

func BenchmarkRawCheck(b *testing.B) {
	b.ReportAllocs()

	rawCases := testToTestRawCases(b, testCases)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for _, c := range rawCases {
			// b.Logf("Test iteration %d: %s", i, c.testDescription)
			err := RawCheck(c.got, c.want)

			if c.expectedResult {
				if err != nil {
					b.FailNow()
				}
			} else {
				if err == nil {
					b.FailNow()
				}
			}
		}
	}
}
