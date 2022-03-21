package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// RoundTripFunc will replace the http.Client's Transport layer
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip implements the http.RoundTripper interface, so we can inject it to the http.Client
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns a new *http.Client with the supplied RoundTripFunc as Transport
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

// RequestRequirements is the collection of possible required Request values
// The MockServer will assert incoming Requests against this set
type RequestRequirements struct {
	URL     *url.URL
	Method  string
	Headers http.Header
	Body    []byte
}

// NewRequestRequirements creates an empty RequestRequirements object
func NewRequestRequirements() *RequestRequirements {
	return &RequestRequirements{
		URL:     &url.URL{},
		Headers: make(http.Header),
		Body:    []byte{},
	}
}

// WithURL adds the supplied URL to the RequestRequirements
func (rr *RequestRequirements) WithURL(URL *url.URL) *RequestRequirements {
	rr.URL = URL
	return rr
}

// WithMethod adds the supplied Method to the RequestRequirements
func (rr *RequestRequirements) WithMethod(Method string) *RequestRequirements {
	rr.Method = Method
	return rr
}

// WithQueryParams adds the supplied QueryParams to the RequestRequirements
func (rr *RequestRequirements) WithQueryParams(QueryParams url.Values) *RequestRequirements {
	rr.URL.RawQuery = QueryParams.Encode()
	return rr
}

// WithHeaders adds the supplied Headers to the RequestRequirements
func (rr *RequestRequirements) WithHeaders(Headers http.Header) *RequestRequirements {
	rr.Headers = Headers
	return rr
}

// WithBody adds the supplied Body to the RequestRequirements
func (rr *RequestRequirements) WithBody(Body []byte) *RequestRequirements {
	rr.Body = Body
	return rr
}

// ResponseValues are the possible Response values
// The MockServer will return these values after the Delay
type ResponseValues struct {
	Delay   time.Duration
	Status  int
	Headers http.Header
	Body    []byte
}

// NewResponseValues creates the default ResponseValues with http.StatusOK 200
func NewResponseValues() *ResponseValues {
	return &ResponseValues{
		Status:  http.StatusOK,
		Headers: make(http.Header),
		Body:    []byte{},
	}
}

// WithDelay adds the supplied Delay to the ResponseValues
func (rv *ResponseValues) WithDelay(Delay time.Duration) *ResponseValues {
	rv.Delay = Delay
	return rv
}

// WithStatus adds the supplied Status to the ResponseValues
func (rv *ResponseValues) WithStatus(Status int) *ResponseValues {
	rv.Status = Status
	return rv
}

// WithHeaders adds the supplied Headers to the ResponseValues
func (rv *ResponseValues) WithHeaders(Headers http.Header) *ResponseValues {
	rv.Headers = Headers
	return rv
}

// WithBody adds the supplied Body to the ResponseValues
func (rv *ResponseValues) WithBody(Body []byte) *ResponseValues {
	rv.Body = Body
	return rv
}

// MockServer is a structure instanciated with a testing context T, a set of RequestRequirements,
// and a set of RepsonseValues. With its HijackClient method it can switch the Transport of a
// http.Client to its RoundTripFunc. This function will assert the outgoing http.Requests of the
// client, and response with preset values.
type MockServer struct {
	T        *testing.T
	Request  *RequestRequirements
	Response *ResponseValues
}

// NewMockServer creates a new MockServer with the supplied teting.T
// empty RequestRequirements and basic ResponseValues
func NewMockServer(t *testing.T) *MockServer {
	return &MockServer{
		T:        t,
		Request:  NewRequestRequirements(),
		Response: NewResponseValues(),
	}
}

// WithRequestURL adds the supplied URL to the MockServer's RequestRequirements
func (ms *MockServer) WithRequestURL(URL *url.URL) *MockServer {
	ms.Request.WithURL(URL)
	return ms
}

// WithRequestMethod adds the supplied Method to the MockServer's RequestRequirements
func (ms *MockServer) WithRequestMethod(Method string) *MockServer {
	ms.Request.WithMethod(Method)
	return ms
}

// WithRequestQueryParams adds the supplied QueryParams to the MockServer's RequestRequirements
func (ms *MockServer) WithRequestQueryParams(QueryParams url.Values) *MockServer {
	ms.Request.WithQueryParams(QueryParams)
	return ms
}

// WithRequestHeaders adds the supplied Headers to the MockServer's RequestRequirements
func (ms *MockServer) WithRequestHeaders(Headers http.Header) *MockServer {
	ms.Request.WithHeaders(Headers)
	return ms
}

// WithRequestBody adds the supplied Body to the MockServer's RequestRequirements
func (ms *MockServer) WithRequestBody(Body []byte) *MockServer {
	ms.Request.WithBody(Body)
	return ms
}

// WithResponseDelay adds the supplied Delay to the MockServer's ResponseValues
func (ms *MockServer) WithResponseDelay(Delay time.Duration) *MockServer {
	ms.Response.WithDelay(Delay)
	return ms
}

// WithResponseStatus adds the supplied Status to the MockServer's ResponseValues
func (ms *MockServer) WithResponseStatus(Status int) *MockServer {
	ms.Response.WithStatus(Status)
	return ms
}

// WithResponseHeaders adds the supplied Headers to the MockServer's ResponseValues
func (ms *MockServer) WithResponseHeaders(Headers http.Header) *MockServer {
	ms.Response.WithHeaders(Headers)
	return ms
}

// WithResponseBody adds the supplied Body to the MockServer's ResponseValues
func (ms *MockServer) WithResponseBody(Body []byte) *MockServer {
	ms.Response.WithBody(Body)
	return ms
}

// HijackClient switches the http.Client's Transport to the MockServer's RoundTrip Function
func (ms *MockServer) HijackClient(Client *http.Client) {
	Client.Transport = ms.GetRoundTripperFn()
}

// GetRoundTripperFn builds up the MockServer's Testing RoundTripFunc from its RequestRequirements and ResponseValues.
// This function will assert incoming Requests against the MockServer's RequestRequirements,
// causing the testing context to fail if any of the checks fails.
// Then it waits for a custom response delay and sends out the MockServer's ResponseValues.
func (ms *MockServer) GetRoundTripperFn() RoundTripFunc {
	return func(req *http.Request) *http.Response {
		assert := require.New(ms.T)

		assert.Equalf(
			ms.Request.URL.Scheme,
			req.URL.Scheme,
			"Required Request Scheme: %s Actual: %s",
			ms.Request.URL.Scheme,
			req.URL.Scheme,
		)

		assert.Equalf(
			ms.Request.URL.Path,
			req.URL.Path,
			"Required Request Path: %s Actual: %s",
			ms.Request.URL.Path,
			req.URL.Path,
		)

		assert.Equalf(
			ms.Request.URL.Host,
			req.URL.Host,
			"Required Request Host: %s Actual: %s",
			ms.Request.URL.Host,
			req.URL.Host,
		)

		assert.Equalf(
			ms.Request.URL.RawQuery,
			req.URL.RawQuery,
			"Required Request RawQuery: %s Actual: %s",
			ms.Request.URL.RawQuery,
			req.URL.RawQuery,
		)

		assert.Equalf(
			ms.Request.Method,
			req.Method,
			"Required Request Method: %s Actual: %s",
			ms.Request.Method,
			req.Method,
		)

		if len(ms.Request.Body) > 0 {
			reqBody, err := ioutil.ReadAll(req.Body)
			assert.NoError(err, "The Request Body should be readable")
			assert.Equalf(
				ms.Request.Body,
				reqBody,
				"Required Request Body: %s Actual: %s",
				ms.Request.Body,
				reqBody,
			)
		}

		if len(ms.Request.Headers) > 0 {
			assert.Equalf(
				ms.Request.Headers,
				req.Header,
				"Required Request Headers: %+v Actual: %+v",
				ms.Request.Headers,
				req.Header,
			)
		}

		reqValues := req.URL.Query()
		if len(reqValues) > 0 {
			assert.Equalf(
				ms.Request.URL.Query(),
				req.URL.Query(),
				"Required Request URL Query Values: %+v Actual: %+v",
				ms.Request.URL.Query(),
				req.URL.Query(),
			)
		}

		if ms.Response.Delay > 0 {
			time.Sleep(ms.Response.Delay)
		}

		return &http.Response{
			Status:        fmt.Sprintf("%d %s", ms.Response.Status, http.StatusText(ms.Response.Status)),
			StatusCode:    ms.Response.Status,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Body:          ioutil.NopCloser(bytes.NewBuffer(ms.Response.Body)),
			ContentLength: int64(len(ms.Response.Body)),
			Request:       req,
			Header:        ms.Response.Headers,
		}
	}
}
