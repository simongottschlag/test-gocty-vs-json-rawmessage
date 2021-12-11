/*
 * Original code by Kevin Gillette (https://github.com/extemporalgenome): https://go.dev/play/p/U_9ejUmD4QJ
 */

package main

import (
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
