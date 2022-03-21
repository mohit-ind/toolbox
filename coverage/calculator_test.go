package coverage

import (
	"io/ioutil"
	"os"
	"testing"

	constants "github.com/toolboxconstants"

	"github.com/stretchr/testify/suite"
)

///////////
// Suite //
///////////

// CalculatorTestSuite extends testify's Suite.
type CalculatorTestSuite struct {
	suite.Suite
}

const failFile = `ok  	github.com/toolboxcli	0.120s	coverage: 🐸% of statements
`

const goldenFile = `ok  	github.com/toolboxcli	0.120s	coverage: 100.0% of statements
ok  	github.com/toolboxconfig	0.217s	coverage: 94.0% of statements
?   	github.com/toolboxconnectors	[no test files]
ok  	github.com/toolboxconstants	0.105s	coverage: 100.0% of statements
?   	github.com/toolboxcoverage	[no test files]
?   	github.com/toolboxexamples	[no test files]
ok  	github.com/toolboxlogger	0.355s	coverage: 100.0% of statements
ok  	github.com/toolboxmiddlewares	0.192s	coverage: 24.4% of statements
?   	github.com/toolboxmodels	[no test files]
ok  	github.com/toolboxrest	0.283s	coverage: 97.1% of statements
?   	github.com/toolboxscripts	[no test files]
ok  	github.com/toolboxtests	0.523s	coverage: 100.0% of statements
`

func (cts *CalculatorTestSuite) prepareTestResultFile(content []byte) string {
	// Create empty tempfile.
	tmpfile, err := ioutil.TempFile(os.TempDir(), "devops-testing")
	cts.NoError(err, "Temp file should have been created")
	cts.NoError(tmpfile.Close(), "Temp file should have been closed")
	// Write content into the tempfile.
	cts.NoError(ioutil.WriteFile(tmpfile.Name(), content, 0644), "The content should be written in the tempfile")
	return tmpfile.Name()
}

func (cts *CalculatorTestSuite) TestNewCalculator() {
	calc := NewCalculator("", 0, nil)
	cts.NotNil(calc, "New Calculator should have been created")
	cts.Error(calc.Scan(), "On empty ResultFile Calculator should fail to Scan")
	cts.Contains(calc.Render(), "NAN%", "Without a successful Scan, Render should return NAN%")
	cts.False(calc.IsCoveredEnough(), "Without a successfull Scan, the achieved test coverage should be 0")
}

func (cts *CalculatorTestSuite) TestCalculatorScanError() {
	testFile := cts.prepareTestResultFile([]byte(failFile))

	calc := NewCalculator(testFile, constants.RequiredTestCoveragePercent, []string{"examples", "models", "scripts"})

	cts.Error(calc.Scan(), "Wrong test result file format, should cause the Scan to fail")
}

func (cts *CalculatorTestSuite) TestCalculatorRender() {
	testFile := cts.prepareTestResultFile([]byte(goldenFile))

	calc := NewCalculator(testFile, constants.RequiredTestCoveragePercent, []string{"examples", "models", "scripts"})

	cts.NoError(calc.Scan(), "The test file should have been successfully scanned by the Calculator")

	res := calc.Render()

	requiredStrings := []string{
		"github.com/toolboxcoverage    | No Tests!",
		"github.com/toolboxmiddlewares | 24.40%",
		"ALL TOGETHER                 |  68.39%",
	}

	for _, requiredString := range requiredStrings {
		cts.Containsf(res, requiredString, "The rendered result should contain: %s", requiredString)
	}

	cts.False(calc.IsCoveredEnough(), "Based on the goldenFile the coverages should be under the required 85%")
}

// TestCalculator runs the whole test suite
func TestCalculator(t *testing.T) {
	suite.Run(t, new(CalculatorTestSuite))
}
