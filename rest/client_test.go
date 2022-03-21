package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	tests "github.com/toolboxtests"

	"github.com/stretchr/testify/suite"

	constants "github.com/toolboxconstants"
)

///////////
// Suite //
///////////

// RestTestSuite extends testify's Suite.
type RestTestSuite struct {
	suite.Suite
}

func (rts *RestTestSuite) TestNewClient() {
	cli := NewClient("test_url", 0)
	rts.NotNil(cli, "Client should have been created")
	rts.Equal("test_url", cli.BaseURL, "The Base URL of the Client should be 'test_url'")
	rts.Equal(
		time.Duration(constants.DefaultHTTPClientTimeoutSec)*time.Second,
		cli.HTTPClient.Timeout,
		"On timeoutSec = 0 input, NewClient should create a Client with the DefaultHTTPClientTimeoutSec",
	)
}

func (rts *RestTestSuite) TestNewClientTimeout() {
	cli := NewClient("test_url", 1)
	rts.Equal(
		time.Second,
		cli.HTTPClient.Timeout,
		"On timeoutSec = 1 input, NewClient should create a Client with 1 sec timeout",
	)
}

func (rts *RestTestSuite) TestGet() {
	cli := NewClient("http://example.com/some/path", 1)
	req, err := cli.MakeNewRequest("GET", nil, nil, nil, nil)
	rts.NoError(err, "New Request should have been created without any error")
	rts.NotNil(req, "New Request should have been created")

	reqURL, err := url.Parse("http://example.com/some/path")
	rts.NoError(err, "the URL should have been parsed")
	srv := tests.NewMockServer(rts.T()).
		WithRequestURL(reqURL).
		WithRequestMethod("GET")

	srv.HijackClient(cli.HTTPClient)

	resp, err := cli.Do(req, nil)
	rts.NoError(err, "The Client should have sent the Request without any error")
	rts.NotNil(resp, "Response should have been returned")
	rts.Equal(http.StatusOK, resp.StatusCode, "Response should be 200 OK")
}

func (rts *RestTestSuite) TestInvalidMethod() {
	cli := NewClient("http://example.com/some/path", 1)
	_, err := cli.MakeNewRequest("Methód", nil, nil, nil, nil)
	rts.Error(err, "MakeNewRequest should throw an error in case of invalid method")
}

func (rts *RestTestSuite) TestPost() {
	testBody := struct {
		Name     string `json:"name"`
		Age      int    `json:"age"`
		IsLegend bool   `json:"is_legend"`
	}{
		Name:     "John Carmack",
		Age:      50,
		IsLegend: true,
	}

	buf := new(bytes.Buffer)
	rts.NoError(json.NewEncoder(buf).Encode(testBody), "The testBody should be encoded as JSON")

	cli := NewClient("http://example.com/some/path", 1)
	req, err := cli.MakeNewRequest("POST", testBody, nil, nil, nil)
	rts.NoError(err, "New Request should have been created without any error")
	rts.NotNil(req, "New Request should have been created")

	reqURL, err := url.Parse("http://example.com/some/path")
	rts.NoError(err, "the URL should have been parsed")
	srv := tests.NewMockServer(rts.T()).
		WithRequestURL(reqURL).
		WithRequestMethod("POST").
		WithRequestBody(buf.Bytes())

	srv.HijackClient(cli.HTTPClient)

	resp, err := cli.Do(req, nil)
	rts.NoError(err, "The Client should have sent the Request without any error")
	rts.NotNil(resp, "Response should have been returned")
	rts.Equal(http.StatusOK, resp.StatusCode, "Response should be 200 OK")
}

func (rts *RestTestSuite) TestGetJSON() {
	type testBody struct {
		Name     string `json:"name"`
		Age      int    `json:"age"`
		IsLegend bool   `json:"is_legend"`
	}

	myTestBody := testBody{
		Name:     "John Carmack",
		Age:      50,
		IsLegend: true,
	}

	buf := new(bytes.Buffer)
	rts.NoError(json.NewEncoder(buf).Encode(myTestBody), "The testBody should be encoded as JSON")

	cli := NewClient("http://example.com/some/path", 1)
	req, err := cli.MakeNewRequest("GET", nil, nil, nil, nil)
	rts.NoError(err, "New Request should have been created without any error")
	rts.NotNil(req, "New Request should have been created")

	reqURL, err := url.Parse("http://example.com/some/path")
	rts.NoError(err, "the URL should have been parsed")
	srv := tests.NewMockServer(rts.T()).
		WithRequestURL(reqURL).
		WithRequestMethod("GET").
		WithResponseBody(buf.Bytes())

	srv.HijackClient(cli.HTTPClient)

	var respBody testBody
	resp, err := cli.Do(req, &respBody)
	rts.NoError(err, "The Client should have sent the Request without any error")
	rts.NotNil(resp, "Response should have been returned")
	rts.Equal(http.StatusOK, resp.StatusCode, "Response should be 200 OK")
	rts.Equal(myTestBody, respBody, "The client should return the test body")
}
func (rts *RestTestSuite) TestPostInvalidBody() {
	testBody := struct {
		InvalidField chan int `json:"invalid_field"`
	}{
		InvalidField: make(chan int),
	}

	cli := NewClient("http://example.com/some/path", 1)
	_, err := cli.MakeNewRequest("POST", testBody, nil, nil, nil)
	rts.Error(err, "MakeNewRequest should return an error in invalid body")
}

func (rts *RestTestSuite) TestSetRequestHeaders() {
	cli := NewClient("http://example.com/some/path", 1)
	req, err := cli.MakeNewRequest("GET", nil, nil, map[string]string{
		"Content-Type":           "test/html",
		"Access-Control-Max-Age": "2520",
	}, nil)
	rts.NoError(err, "New Request should have been created without any error")
	rts.NotNil(req, "New Request should have been created")

	reqURL, err := url.Parse("http://example.com/some/path")
	rts.NoError(err, "the URL should have been parsed")
	srv := tests.NewMockServer(rts.T()).
		WithRequestURL(reqURL).
		WithRequestMethod("GET").
		WithRequestHeaders(http.Header{
			"Content-Type":           []string{"test/html"},
			"Access-Control-Max-Age": []string{"2520"},
		})

	srv.HijackClient(cli.HTTPClient)

	resp, err := cli.Do(req, nil)
	rts.NoError(err, "The Client should have sent the Request without any error")
	rts.NotNil(resp, "Response should have been returned")
	rts.Equal(http.StatusOK, resp.StatusCode, "Response should be 200 OK")
}

func (rts *RestTestSuite) TestAddRequestHeaders() {
	cli := NewClient("http://example.com/some/path", 1)
	req, err := cli.MakeNewRequest("GET", nil, nil, nil, map[string]string{
		"Content-Type":           "test/html",
		"Access-Control-Max-Age": "2520",
	})
	rts.NoError(err, "New Request should have been created without any error")
	rts.NotNil(req, "New Request should have been created")

	reqURL, err := url.Parse("http://example.com/some/path")
	rts.NoError(err, "the URL should have been parsed")
	srv := tests.NewMockServer(rts.T()).
		WithRequestURL(reqURL).
		WithRequestMethod("GET").
		WithRequestHeaders(http.Header{
			"Content-Type":           []string{"test/html"},
			"Access-Control-Max-Age": []string{"2520"},
		})

	srv.HijackClient(cli.HTTPClient)

	resp, err := cli.Do(req, nil)
	rts.NoError(err, "The Client should have sent the Request without any error")
	rts.NotNil(resp, "Response should have been returned")
	rts.Equal(http.StatusOK, resp.StatusCode, "Response should be 200 OK")
}

func (rts *RestTestSuite) TestQuery() {
	cli := NewClient("http://example.com/some/path", 1)
	req, err := cli.MakeNewRequest("GET", nil, map[string]string{
		"Query1": "Value1",
		"Query2": "12345",
		"Query3": "[{?}]",
	}, nil, nil)
	rts.NoError(err, "New Request should have been created without any error")
	rts.NotNil(req, "New Request should have been created")

	reqURL, err := url.Parse("http://example.com/some/path")
	rts.NoError(err, "the URL should have been parsed")
	srv := tests.NewMockServer(rts.T()).
		WithRequestURL(reqURL).
		WithRequestMethod("GET").
		WithRequestQueryParams(url.Values{
			"Query1": []string{"Value1"},
			"Query2": []string{"12345"},
			"Query3": []string{"[{?}]"},
		})

	srv.HijackClient(cli.HTTPClient)

	resp, err := cli.Do(req, nil)
	rts.NoError(err, "The Client should have sent the Request without any error")
	rts.NotNil(resp, "Response should have been returned")
	rts.Equal(http.StatusOK, resp.StatusCode, "Response should be 200 OK")
}

// TestRest runs the whole test suite
func TestRest(t *testing.T) {
	suite.Run(t, new(RestTestSuite))
}
