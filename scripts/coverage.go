package main

import (
	"fmt"
	"log"
	"os"

	constants "github.com/toolboxconstants"
	coverage "github.com/toolboxcoverage"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("This script needs an argument: the path of the Go test result file")
	}
	calc := coverage.NewCalculator(os.Args[1], constants.RequiredTestCoveragePercent, []string{
		"scripts",
		"examples",
		"examples/migration",
		"examples/migration/statik",
		"examples/logger",
		"examples/logger/component",
		"examples/request-logger",
		"examples/cli",
		"examples/logger-global",
		"examples/config-global/printer",
		"examples/config-global",
		"examples/config-global/config",
		"examples/logger-global/component",
		"examples/logger-global/log",
	})
	if err := calc.Scan(); err != nil {
		log.Fatalf("Coverage Calculator Failed: %s", err)
	}
	fmt.Println(calc.Render())
	if !calc.IsCoveredEnough() {
		log.Fatalf("The required %v%% test coverage was not achieved!", constants.RequiredTestCoveragePercent)
	}
}
