package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/toolboxexamples/logger-global/component"
	"github.com/toolboxexamples/logger-global/log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Failed to load configs: %s\n", err)
		os.Exit(1)
	}

	log.Entry().Debug("Global logger has been set up - debug level")
	log.Entry().Info("Global logger has been set up - info level")
	log.Entry().Warn("Global logger has been set up - warn level")
	log.Entry().Error("Global logger has been set up - error level")

	componentLogger := log.Global().NewComponentLogger("My fancy component")

	myComponent := component.NewCoolComponent(componentLogger)

	myComponent.DoSomething()

	if err := myComponent.CallSomethingElse(); err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"Extra field string": "Extra string value",
			"Extra field number": 12345,
			"Extra field bool":   true,
		}).Fatal("Fatal Error occurred")
	}
}
