package printer

import (
	"fmt"
	"os"

	constants "github.com/toolboxconstants"

	"github.com/toolboxexamples/config-global/config"
)

func PrintConfigTable() {
	fmt.Printf("\nConfig Table:\n%+v\n", config.Global().DumpTable())
}

func PrintEnvironmentVariables() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	fmt.Printf("Environment Variables:\n")
	fmt.Println(constants.APP_PORT, ":", os.Getenv(constants.APP_PORT))
	fmt.Println(constants.EC2_ID, ":", os.Getenv(constants.EC2_ID))
	fmt.Println("Actual hostname", ":", hostname)
	fmt.Println(constants.APP_ENV, ":", os.Getenv(constants.APP_ENV))
	fmt.Println(constants.APP_DB_SECRET_NAME, ":", os.Getenv(constants.APP_DB_SECRET_NAME))
	fmt.Println(constants.APP_DEBUG, ":", os.Getenv(constants.APP_DEBUG))
	fmt.Println(constants.APP_LOG_LEVEL, ":", os.Getenv(constants.APP_LOG_LEVEL))
	fmt.Println(constants.APP_LOG_DEV, ":", os.Getenv(constants.APP_LOG_DEV))
	fmt.Println(constants.APP_LOG_FORMAT_ERRORS, ":", os.Getenv(constants.APP_LOG_FORMAT_ERRORS))
}

func PrintConfigs() {
	fmt.Printf("\nConfiguration Items:\n")
	fmt.Println(constants.APP_PORT, ":", config.Global().Port())
	fmt.Println("Hostname", ":", config.Global().Hostname())
	fmt.Println(constants.APP_ENV, ":", config.Global().Env())
	fmt.Println(constants.APP_DB_SECRET_NAME, ":", config.Global().DBSecretName())
	fmt.Println(constants.APP_LOG_LEVEL, ":", config.Global().LogLevel())
	fmt.Println("Logrus log level", ":", config.Global().LogrusLogLevel())
	fmt.Println("IsDebug", ":", config.Global().IsDebug())
	fmt.Println("IsDev", ":", config.Global().IsDev())
	fmt.Println("IsTest", ":", config.Global().IsTest())
	fmt.Println("IsStaging", ":", config.Global().IsStaging())
	fmt.Println("IsAcceptance", ":", config.Global().IsAcceptance())
	fmt.Println("IsProduction", ":", config.Global().IsProduction())
}

func PrintConfig(name string) {
	fmt.Printf("\nGet %s : %s\n", name, config.Global().Get(name))
}

func LookupConfigItem(name string) {
	val, found := config.Global().Lookup(name)
	fmt.Printf("\nLookup Config Item by Name: %s Found: %t Value: %s\n", name, found, val)
}
