package main

import (
	"fmt"
	"os"

	"github.com/toolboxexamples/config-global/config"
	"github.com/toolboxexamples/config-global/printer"
)

func main() {
	// Set the global configs at the very begining of the program
	// Exit if cannot set up configs
	if err := config.SetupConfigs(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print the Configuration Table
	printer.PrintConfigTable()

	// Print actual environment
	printer.PrintEnvironmentVariables()

	// Print config items with helpher methods
	printer.PrintConfigs()

	// Print a config item by its name (if it is not found an empty string will be printed)
	printer.PrintConfig("APP_LOG_FORMAT_ERRORS")
	printer.PrintConfig("THIS_VARIABLE_IS_NOT_EXISTS")

	// Lookup a config item by its name, print a boolean indicating if it was found and the value (or an empty string if it is not found)
	printer.LookupConfigItem("APP_PORT")
	printer.LookupConfigItem("MY_FANCY_VARIABLE")

	fmt.Println("Creating .env.sample file...")
	if err := config.Global().CreateSampleFile(".env.sample"); err != nil {
		fmt.Printf("Failed to create .env.sample file: %s\n", err)
		os.Exit(1)
	}
}
