package main

import (
	"github.com/sirupsen/logrus"
	"github.com/toolboxexamples/logger/component"
	logger "github.com/toolboxlogger"
)

func main() {
	mainLogger := logger.NewCommonLogger("Test App", "v1.3.2", "test", "localhost", true)

	mainLogger.Entry().Debug("Logger has been set up - debug level")
	mainLogger.Entry().Info("Logger has been set up - info level")
	mainLogger.Entry().Warn("Logger has been set up - warn level")
	mainLogger.Entry().Error("Logger has been set up - error level")

	componentLogger := mainLogger.NewComponentLogger("My fancy component")

	myComponent := component.NewCoolComponent(componentLogger)

	myComponent.DoSomething()

	if err := myComponent.CallSomethingElse(); err != nil {
		mainLogger.WithError(err).WithFields(logrus.Fields{
			"Extra field string": "Extra string value",
			"Extra field number": 12345,
			"Extra field bool":   true,
		}).Fatal("Fatal Error occurred")
	}
}
