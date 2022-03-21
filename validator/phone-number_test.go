package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDutchPhoneNumber(t *testing.T) {
	assert := require.New(t)

	testCases := map[string]struct {
		Input          string
		ExpectedOutput string
		ExpectedError  string
	}{
		"Empty Input": {
			ExpectedError: " is not a Dutch Mobile Number",
		},
		"Short Dutch": {
			Input:          "0634482527",
			ExpectedOutput: "0634482527",
		},
		"Short Dutch with spaces added": {
			Input:          "06 34 48 25 27",
			ExpectedOutput: "0634482527",
		},
		"Short Dutch with space and hyphens added": {
			Input:          "06 34-4825-27",
			ExpectedOutput: "0634482527",
		},
		"Long Dutch": {
			Input:          "+31634482527",
			ExpectedOutput: "0634482527",
		},
		"Long Dutch with spaces added": {
			Input:          "+31 6 34 48 25 27",
			ExpectedOutput: "0634482527",
		},
		"Long Dutch with a space and hyphens added": {
			Input:          "+31 6-3448-2527",
			ExpectedOutput: "0634482527",
		},
		"Too short": {
			Input:         "063448252",
			ExpectedError: "063448252 is not a Dutch Mobile Number",
		},
		"Too long number": {
			Input:         "06556843081",
			ExpectedError: "06556843081 is not a Dutch Mobile Number",
		},
		"Too long number, prefixed with +31": {
			Input:         "+316556843081",
			ExpectedError: "+316556843081 is not a Dutch Mobile Number",
		},
		"Not begining with 06": {
			Input:         "0534482527",
			ExpectedError: "0534482527 is not a Dutch Mobile Number",
		},
		"Prefixed with both +31 and 0": {
			Input:         "+310634482527",
			ExpectedError: "+310634482527 is not a Dutch Mobile Number",
		},
		"Invalid character": {
			Input:         "063A482527",
			ExpectedError: "063A482527 is not a Dutch Mobile Number",
		},
	}

	for testCaseName, testCase := range testCases {
		t.Logf("Testing ParseDutchPhoneNumber: %s", testCaseName)
		output, err := ParseDutchPhoneNumber(testCase.Input)
		if testCase.ExpectedError != "" {
			assert.EqualErrorf(
				err,
				testCase.ExpectedError,
				"Expected error: %s Actual error: %s",
				testCase.ExpectedError,
				err,
			)
		} else {
			assert.NoError(err, "There should be no error in this case")
		}
		assert.Equalf(
			testCase.ExpectedOutput,
			output,
			"Expected output: %s Actual: %s",
			testCase.ExpectedOutput,
			output,
		)
	}
}
