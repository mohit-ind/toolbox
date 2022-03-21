package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	toolLog "github.com/toolboxlogger"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/suite"
)

// RequestLoggerSuite extends testify's Suite.
type RequestLoggerSuite struct {
	suite.Suite
}

// TestRequestLoggerSuite runs all test cases against the LoggingMiddleware
func (rls *RequestLoggerSuite) TestRequestLogger() {

	getBasicTestHandler := func(status int) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			_, err := w.Write([]byte("Testing\n"))
			rls.NoError(err, "Response should be written")
		})
	}

	getDoubleWriteHandler := func(status int) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			w.WriteHeader(http.StatusTeapot)
			_, err := w.Write([]byte("Testing Double Header\n"))
			rls.NoError(err, "Response should be written")
		})
	}

	panicingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("PANIC!")
	})

	testCases := map[string]struct {
		msg         string
		statusCode  int
		logLvl      logrus.Level
		extraFields logrus.Fields
		method      string
		path        string
		handler     string
	}{
		"200 OK": {
			statusCode: http.StatusOK,
			logLvl:     logrus.InfoLevel,
			method:     "GET",
			path:       "/",
			msg:        "Request",
		},
		"Extra Fields": {
			statusCode: http.StatusOK,
			logLvl:     logrus.InfoLevel,
			method:     "GET",
			path:       "/",
			extraFields: logrus.Fields{
				"Int Field":    1,
				"String Field": "I am a string",
				"Bool Field":   false,
			},
			msg: "Request",
		},
		"404 Not Found": {
			statusCode: http.StatusNotFound,
			logLvl:     logrus.WarnLevel,
			method:     "GET",
			path:       "/cheese",
			msg:        "Request",
		},
		"Double write header": {
			statusCode: http.StatusCreated,
			logLvl:     logrus.InfoLevel,
			method:     "POST",
			path:       "/articles/123",
			handler:    "double",
			msg:        "Request",
		},
		"Panic!": {
			statusCode: http.StatusInternalServerError,
			logLvl:     logrus.ErrorLevel,
			method:     "UPDATE",
			path:       "/elephants",
			handler:    "panic",
			msg:        "Recovered from panic",
		},
		"Panic! And Extra Fields": {
			statusCode: http.StatusInternalServerError,
			logLvl:     logrus.ErrorLevel,
			method:     "UPDATE",
			path:       "/elephants",
			extraFields: logrus.Fields{
				"Int Field":    1,
				"String Field": "I am a string",
				"Bool Field":   false,
			},
			handler: "panic",
			msg:     "Recovered from panic",
		},
	}

	// create a null logger, which does not writes the log to stdout, but has a hook which we can check for log entries
	logger, hook := test.NewNullLogger()

	// iterate over test cases
	for name, tCase := range testCases {
		// log test case name (this will only be visible if a test fails)
		rls.T().Logf("Request Logger Middleware: %s", name)

		// create a test request
		req, err := http.NewRequest(tCase.method, tCase.path, nil)
		rls.NoError(err, "Test request should be created")

		// create a response recorder
		rr := httptest.NewRecorder()

		// depending on the test case choose a testHandler
		testHandler := getBasicTestHandler(tCase.statusCode)
		switch tCase.handler {
		case "double":
			testHandler = getDoubleWriteHandler(tCase.statusCode)
		case "panic":
			testHandler = panicingHandler
		}

		appLogger := toolLog.NewLogger(logger, tCase.extraFields)
		// create the logging middleware and wrap the testHandler with it
		handler := LoggingMiddleware(appLogger)(testHandler)

		// issue the test request
		rls.NotPanics(func() { handler.ServeHTTP(rr, req) }, "The panic should have been caught by the middleware")

		// check response status code
		rls.Equalf(tCase.statusCode, rr.Code, "handler returned wrong status code: want %v got %v", tCase.statusCode, rr.Code)

		// check the log entry
		rls.Equal(1, len(hook.Entries), "There should be only one log entry now")
		rls.Equalf(tCase.logLvl, hook.LastEntry().Level, "The log level should be %v", tCase.logLvl)
		rls.Containsf(hook.LastEntry().Message, tCase.msg, "The log message should contain: %s", tCase.msg)

		// check extra fields, if there are any
		for key, val := range tCase.extraFields {
			rls.Equalf(val, hook.LastEntry().Data[key], "Log entry should have key: %s value: %v", key, val)
		}

		// if no panic occurred, check the request log fields
		if tCase.handler != "panic" {
			rls.Equalf(tCase.statusCode, hook.LastEntry().Data["status"], "Status in the log message should be %d", tCase.statusCode)
			rls.Equalf(tCase.method, hook.LastEntry().Data["method"], "Method in the log message should be %s", tCase.method)
			rls.Equalf(tCase.path, hook.LastEntry().Data["path"], "Path in the log message should be %d", tCase.path)
			val := hook.LastEntry().Data["duration"].(time.Duration)
			rls.True(int64(val) > int64(time.Duration(time.Nanosecond*100)), "Duration in the log message should be greater than 100ns")
		}

		// reset the test hook, check if it is empty now
		hook.Reset()
		rls.Nil(hook.LastEntry())
	}

}

// TestRequestLogger runs the suite
func Test_RequestLogger(t *testing.T) {
	suite.Run(t, new(RequestLoggerSuite))
}
