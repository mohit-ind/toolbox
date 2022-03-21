package coverage

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

// Calculator is responsible for opening and scanning Go test result file,
// collecting each package's test coverage percentage while skipping the ones
// matching any string in the Calculator's skip list.
// After scanning Calculator can Render a table of the results, and tell
// if the requiredPercentage is fulfilled or not.
type Calculator struct {
	resultFile      string
	requiredPercent float64
	skip            []string
	coverages       map[string]float64
	sum             float64
}

// NewCalculator creates a new Calculator with the supplied ResultFile path, RequiredPercentage, and Skip list.
func NewCalculator(ResultFile string, RequiredPercent float64, Skip []string) *Calculator {
	return &Calculator{
		resultFile:      ResultFile,
		requiredPercent: RequiredPercent,
		skip:            Skip,
		coverages:       make(map[string]float64),
	}
}

// coverage returns the sum of all coverages per the number of checked packages.
func (calc *Calculator) coverage() float64 {
	return calc.sum / float64(len(calc.coverages))
}

// IsCoveredEnough returns true if the coverage is equals to or more than the requiredPercent.
func (calc *Calculator) IsCoveredEnough() bool {
	return calc.coverage() >= calc.requiredPercent
}

// contains check if s string contains any of the strings in a slice.
func (calc *Calculator) contains(s string, elems []string) bool {
	for _, elem := range elems {
		if strings.HasSuffix(s, elem) {
			return true
		}
	}
	return false
}

// Scan tries to open the Calculator's resultFile and scan its lines for test results,
// it populates the Calculator coverages list with packages names and coverage percentages.
// Scan may return an optional error if the resultFile cannot be opened, the percentage string
// cannot be converted to float64 or the line-scanner fails.
// Scan skips every packages which name contains any of the strings in the Calculator's skip list.
func (calc *Calculator) Scan() error {
	file, err := os.Open(calc.resultFile)
	if err != nil {
		return errors.Wrap(err, "Failed to open test result file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 1 {
			split := strings.Split(line, "\t")
			packageName := split[1]
			if calc.contains(packageName, calc.skip) {
				continue
			}
			if strings.HasPrefix(split[0], "?") {
				calc.coverages[packageName] = 0
				continue
			}
			covStr := split[3]
			percentStr := strings.TrimPrefix(covStr, "coverage: ")
			percentStr = strings.TrimSuffix(percentStr, "% of statements")
			percent, err := strconv.ParseFloat(percentStr, 64)
			if err != nil {
				return errors.Wrapf(err, "Failed to convert %s to float64", percentStr)
			}
			calc.coverages[packageName] = percent
			calc.sum += percent
		}
	}

	return scanner.Err()
}

// Render creates a olekukonko/tablewriter table with all the checked package names and percentages,
// it adds a summary as footer and returns the table as a string.
func (calc *Calculator) Render() string {
	data := [][]string{}
	for key, val := range calc.coverages {
		percentage := "No Tests!"
		if val > 0 {
			percentage = fmt.Sprintf("%.2f%%", val)
		}
		data = append(data, []string{key, percentage})
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Package", "Coverage"})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowSeparator("-")
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(data)
	table.SetFooter([]string{"ALL TOGETHER", fmt.Sprintf("%.2f%%", calc.coverage())})
	table.Render()

	return tableString.String()
}
