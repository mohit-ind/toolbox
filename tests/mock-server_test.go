package tests

//////////////////////////////////////////////////////////////////////////////////////////////////|
//                                                                                               ||
// This test suite proves that the methodology described in this article:                        ||
// http://hassansin.github.io/Unit-Testing-http-client-in-Go                                     ||
//                                                                                               ||
// Namely: hijacking the http.Client by switching its RoundTripper with a custom one, is an easy //
// and convenient way to abstract the network layer away, and focus on the http.Request and     //
// http.Response objects in our rest client tests. This allows us to write simple unit         //
// tests against the outgoing http calls, proveing that they use then intended;               //
// method, uri, headers, body.                                                               //
// Also we can create fake responses with same capabilities and optional delay.             //
//                                                                                         //
////////////////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

/////////////////
// Test Suite //
///////////////

// HTTPTestSuite extends testify's Suite.
type HTTPTestSuite struct {
	suite.Suite
}

func (hts *HTTPTestSuite) basicResp() *http.Response {
	return &http.Response{
		Status:     "OK",
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("OK"))),
	}
}

func (hts *HTTPTestSuite) basicRoundTripFunc() RoundTripFunc {
	return RoundTripFunc(func(req *http.Request) *http.Response {
		return hts.basicResp()
	})
}

/////////////////////////
// RoundTripper Tests //
///////////////////////

func (hts *HTTPTestSuite) TestRoundTripper() {
	expectedResponse := hts.basicResp()
	fn := hts.basicRoundTripFunc()
	roundResp, err := fn.RoundTrip(&http.Request{})
	hts.NoError(err, "RoundTrip shouldn't return any error")
	hts.Equalf(expectedResponse, roundResp, "The Response should be: %+v Actual: %+v", expectedResponse, roundResp)
}

func (hts *HTTPTestSuite) TestNewclient() {
	cli := NewTestClient(hts.basicRoundTripFunc())
	hts.NotNil(cli, "HTTP Client should have been created")

	resp, err := cli.Get("/")
	hts.NoError(err, "The Client should perform the GET Request without any error")

	expectedResponse := hts.basicResp()
	hts.Equalf(expectedResponse, resp, "The Response should be: %+v Actual: %+v", expectedResponse, resp)
}

////////////////////////////////
// RequestRequirements Tests //
//////////////////////////////

func (hts *HTTPTestSuite) TestNewRequestRequirements() {
	req := NewRequestRequirements()
	hts.NotNil(req, "New RequestRequirements should have been created")
	hts.Equal(
		[]byte{},
		req.Body,
		"The default Body of the RequestRequirements should be an empty byte slice",
	)
	hts.Equal(
		http.Header{},
		req.Headers,
		"The default Headers of the RequestRequirements should be an empty http.Header",
	)
	hts.Equalf(
		"",
		req.Method,
		"The default Method of the RequestRequirements should be an empty string",
	)
	hts.Equalf(
		&url.URL{},
		req.URL,
		"The default URL of the RequestRequirements should be an empty *url.URL",
	)
}

func (hts *HTTPTestSuite) TestRequestRequirementsWithURL() {
	testURL, err := url.Parse("wss://www.testhost.com:8080/test/url?queryKey=queryVal#testHash")
	hts.NoError(err, "Test URL should have been parsed")

	req := NewRequestRequirements().WithURL(testURL)

	hts.Equalf(
		testURL,
		req.URL,
		"The RequestRequirement's URL should be: %+v Actual: %+v",
		testURL,
		req.URL,
	)
}

func (hts *HTTPTestSuite) TestRequestRequirementsWithMethod() {
	req := NewRequestRequirements().WithMethod("PUT")
	hts.Equal(
		"PUT",
		req.Method,
		"The RequestRequirement's Method should be 'PUT'",
	)
}

func (hts *HTTPTestSuite) TestRequestRequirementsWithQueryParams() {
	queryParamas := url.Values{
		"Key1": []string{"value1"},
		"Key2": []string{"value2"},
	}
	req := NewRequestRequirements().WithQueryParams(queryParamas)
	hts.Equal(
		queryParamas,
		req.URL.Query(),
		"The RequestRequirement's Method should be 'PUT'",
	)
}

func (hts *HTTPTestSuite) TestRequestRequirementsWithHeaders() {
	headers := http.Header{
		"Content-Type":           []string{"test/html"},
		"Access-Control-Max-Age": []string{"2520"},
	}
	req := NewRequestRequirements().WithHeaders(headers)
	hts.Equal(
		headers,
		req.Headers,
		"The RequestRequirement's Headers should be: %+v Actual: %+v",
		headers,
		req.Headers,
	)
}

func (hts *HTTPTestSuite) TestRequestRequirementsWithBody() {
	body := []byte("This is the body")
	req := NewRequestRequirements().WithBody(body)
	hts.Equal(
		body,
		req.Body,
		"The RequestRequirement's Body should be: %+v Actual: %+v",
		body,
		req.Body,
	)
}

///////////////////////////
// ResponseValues Tests //
/////////////////////////

func (hts *HTTPTestSuite) TestNewResponseValues() {
	resp := NewResponseValues()
	hts.NotNil(resp, "ResponseValues should have been created")
	hts.Equal(http.StatusOK, resp.Status, "The default Response Status should be 200 http.StatusOK")
	hts.Equal(http.Header{}, resp.Headers, "The default Response Headers should be an empty http.Header")
	hts.Equal([]byte{}, resp.Body, "The default Response Body should be an empty byte slice")
	hts.Equal(time.Duration(0), resp.Delay, "The default Response Delay should be an Duration 0")
}

func (hts *HTTPTestSuite) TestResponseValuesWithDelay() {
	delay := time.Millisecond * 100
	resp := NewResponseValues().WithDelay(delay)
	hts.Equal(delay, resp.Delay, "The ResponseValue's Delay should be 100 milliseconds")
}

func (hts *HTTPTestSuite) TestResponseValuesWithStatus() {
	status := http.StatusTeapot
	resp := NewResponseValues().WithStatus(status)
	hts.Equal(status, resp.Status, "The ResponseValue's Status should be 418 I'm a teapot")
}

func (hts *HTTPTestSuite) TestResponseValuesWithHeaders() {
	headers := http.Header{
		"Content-Type":           []string{"test/html"},
		"Access-Control-Max-Age": []string{"2520"},
	}
	resp := NewResponseValues().WithHeaders(headers)
	hts.Equal(
		headers,
		resp.Headers,
		"The ResponseValue's Headers should be: %+v Actual: %+v",
		headers,
		resp.Headers,
	)
}

func (hts *HTTPTestSuite) TestResponseValuesWithBody() {
	body := []byte("This is the body")
	resp := NewResponseValues().WithBody(body)
	hts.Equal(
		body,
		resp.Body,
		"The ResponseValue's Body should be: %+v Actual: %+v",
		body,
		resp.Body,
	)
}

///////////////////////
// MockServer Tests //
/////////////////////

func (hts *HTTPTestSuite) TestNewMockServer() {
	srv := NewMockServer(hts.T())
	hts.NotNil(srv, "MockServer should have been created")
	hts.NotNil(srv.T, "the testing context's T should have been passed to the MockServer")
	hts.NotNil(srv.Request, "MockServer's RequestRequirements should have been created")
	hts.NotNil(srv.Response, "MockServer's ResponseValues should have been created")
}

func (hts *HTTPTestSuite) TestMockServerWithRequestURL() {
	testURL, err := url.Parse("wss://www.testhost.com:8080/test/url?queryKey=queryVal#testHash")
	hts.NoError(err, "Test URL should have been parsed")
	srv := NewMockServer(hts.T()).WithRequestURL(testURL)
	hts.Equalf(
		testURL,
		srv.Request.URL,
		"The MockServer's RequestRequirement's URL should be: %+v Actual: %+v",
		testURL,
		srv.Request.URL,
	)
}

func (hts *HTTPTestSuite) TestMockServerWithRequestMethod() {
	srv := NewMockServer(hts.T()).WithRequestMethod("PUT")
	hts.Equal(
		"PUT",
		srv.Request.Method,
		"The MockServer's RequestRequirement's Method should be 'PUT'",
	)
}

func (hts *HTTPTestSuite) TestMockServerWithRequestQueryParams() {
	queryParamas := url.Values{
		"Key1": []string{"value1"},
		"Key2": []string{"value2"},
	}
	srv := NewMockServer(hts.T()).WithRequestQueryParams(queryParamas)
	hts.Equal(
		queryParamas,
		srv.Request.URL.Query(),
		"The MockServer's RequestRequirement's Method should be 'PUT'",
	)
}

func (hts *HTTPTestSuite) TestMockServerWithRequestHeaders() {
	headers := http.Header{
		"Content-Type":           []string{"test/html"},
		"Access-Control-Max-Age": []string{"2520"},
	}
	srv := NewMockServer(hts.T()).WithRequestHeaders(headers)
	hts.Equal(
		headers,
		srv.Request.Headers,
		"The MockServer's RequestRequirement's Headers should be: %+v Actual: %+v",
		headers,
		srv.Request.Headers,
	)
}

func (hts *HTTPTestSuite) TestMockServerWithRequestBody() {
	body := []byte("This is the body")
	srv := NewMockServer(hts.T()).WithRequestBody(body)
	hts.Equal(
		body,
		srv.Request.Body,
		"The MockServer's RequestRequirement's Body should be: %+v Actual: %+v",
		body,
		srv.Request.Body,
	)
}

func (hts *HTTPTestSuite) TestMockServerWithResponseDelay() {
	delay := time.Millisecond * 100
	srv := NewMockServer(hts.T()).WithResponseDelay(delay)
	hts.Equal(delay, srv.Response.Delay, "The MockServer's ResponseValue's Delay should be 100 milliseconds")
}

func (hts *HTTPTestSuite) TestMockServerWithResponseStatus() {
	status := http.StatusTeapot
	srv := NewMockServer(hts.T()).WithResponseStatus(status)
	hts.Equal(status, srv.Response.Status, "The MockServer's ResponseValue's Status should be 418 I'm a teapot")
}

func (hts *HTTPTestSuite) TestMockServerWithResponseHeaders() {
	headers := http.Header{
		"Content-Type":           []string{"test/html"},
		"Access-Control-Max-Age": []string{"2520"},
	}
	srv := NewMockServer(hts.T()).WithResponseHeaders(headers)
	hts.Equal(
		headers,
		srv.Response.Headers,
		"The MockServer's ResponseValue's Headers should be: %+v Actual: %+v",
		headers,
		srv.Response.Headers,
	)
}

func (hts *HTTPTestSuite) TestMockServerWithResponseBody() {
	body := []byte("This is the body")
	srv := NewMockServer(hts.T()).WithResponseBody(body)
	hts.Equal(
		body,
		srv.Response.Body,
		"The MockServer's ResponseValue's Body should be: %+v Actual: %+v",
		body,
		srv.Response.Body,
	)
}

func (hts *HTTPTestSuite) TestMockServerHijackClient() {
	cli := &http.Client{}
	srv := NewMockServer(hts.T())
	srv.HijackClient(cli)
	_, ok := cli.Transport.(RoundTripFunc)
	hts.True(ok, "The MockServer should have successfully hijacked the http.Client")
}

//////////////////////////////
// Putting all together... //
////////////////////////////

func (hts *HTTPTestSuite) TestMockServerAgainstVariousCases() {
	testCases := map[string]struct {
		requestURI     string
		requestMethod  string
		requestHeader  http.Header
		requestBody    []byte
		responseDelay  time.Duration
		responseStatus int
		responseHeader http.Header
		responseBody   []byte
	}{
		"Basic 200 OK": {
			requestMethod: "GET",
		},
		"Basic 200 OK With URI": {
			requestURI:    "http://john_cena:pass123@wwf.org:8080/admin/pages?paginator=10&lowdash=off#asdf123",
			requestMethod: "GET",
		},
		"201 Created With URI": {
			requestURI:     "http://john_cena:pass123@wwf.org:8080/admin/pages?paginator=10&lowdash=off#asdf123",
			requestMethod:  "GET",
			responseStatus: http.StatusCreated,
		},
		"POST 201 Created With URI": {
			requestURI:     "http://john_cena:pass123@wwf.org:8080/admin/pages?paginator=10&lowdash=off#asdf123",
			requestMethod:  "POST",
			responseStatus: http.StatusCreated,
		},
		"Extra Request headers": {
			requestURI:    "example.com",
			requestMethod: "GET",
			requestHeader: http.Header{
				"Content-Type":           []string{"test/html"},
				"Access-Control-Max-Age": []string{"2520"},
			},
		},
		"Request Body": {
			requestURI:    "example.com",
			requestMethod: "PUT",
			requestBody:   []byte("I'm the request body"),
		},
		"Request with some Delay": {
			requestURI:    "example.com",
			requestMethod: "OPTIONS",
			responseDelay: time.Millisecond * 100,
		},
		"Response status": {
			requestMethod:  "GET",
			responseStatus: http.StatusTeapot,
		},
		"Response header": {
			requestMethod: "GET",
			responseHeader: http.Header{
				"Content-Type": []string{"application/json"},
			},
		},
		"Response body": {
			requestMethod: "GET",
			responseBody:  []byte(`{"json-key":"json-value"}`),
		},
	}

	for testCaseName, testCase := range testCases {
		hts.T().Logf("MockServer Test: %s", testCaseName)

		reqURL, err := url.Parse(testCase.requestURI)
		hts.NoError(err, "The test case's requestURI should have been parsed")

		srv := NewMockServer(hts.T()).
			WithRequestURL(reqURL).
			WithRequestMethod(testCase.requestMethod).
			WithRequestHeaders(testCase.requestHeader).
			WithRequestBody(testCase.requestBody).
			WithResponseDelay(testCase.responseDelay).
			WithResponseStatus(testCase.responseStatus).
			WithResponseHeaders(testCase.responseHeader).
			WithResponseBody(testCase.responseBody)

		cli := &http.Client{}

		srv.HijackClient(cli)

		req, err := http.NewRequest(
			testCase.requestMethod,
			testCase.requestURI,
			ioutil.NopCloser(bytes.NewReader(testCase.requestBody)),
		)

		hts.NoError(err, "The test http Request should have been created without any error")

		req.Header = testCase.requestHeader

		beforeReq := time.Now()
		resp, err := cli.Do(req)
		hts.NoError(err, "The test response http Response should have arrived")
		afterReq := time.Now()
		hts.Greater(
			int64(afterReq.Sub(beforeReq)),
			int64(testCase.responseDelay),
			"The test http Response should have arrived a bit later than the delay",
		)

		hts.NotNil(resp, "A http Response should have been returned")
		hts.Equalf(
			testCase.responseStatus,
			resp.StatusCode,
			"Required response status: %d Actual: %d",
			testCase.responseStatus,
			resp.StatusCode,
		)
		hts.Equalf(
			testCase.responseHeader,
			resp.Header,
			"Required response header: %v Actual: %v",
			testCase.responseHeader,
			resp.Header,
		)
		if len(testCase.responseBody) > 0 {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			hts.NoError(err, "The response body should be readable")
			hts.Equalf(
				testCase.responseBody,
				bodyBytes,
				"Required response body: %s Actual: %s",
				testCase.responseBody,
				bodyBytes,
			)
		}
	}
}

/////////////////////////
// Run the Test Suite //
///////////////////////
func TestHTTP(t *testing.T) {
	suite.Run(t, new(HTTPTestSuite))
}
