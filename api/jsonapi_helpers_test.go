package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/OpenBazaar/openbazaar-go/test"
	"github.com/stretchr/testify/assert"

	manet "gx/ipfs/QmPpRcbNUXauP3zWZ1NJMLWpe4QnmEHrd2ba2D3yqWznw7/go-multiaddr-net"
	ma "gx/ipfs/QmYzDkkgAEmrcNzFCiYo6L1dTX4EAG1gZkbtdbd9trL4vd/go-multiaddr"
)

// testURIRoot is the root http URI to hit for testing
const testURIRoot = "http://127.0.0.1:9191"

// anyResponseJSON is a sentinel denoting any valid JSON response body is valid
const anyResponseJSON = "__anyresponsebodyJSON__"

// testHTTPClient is the http client to use for tests
var testHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

// newTestGateway starts a new API gateway listening on the default test interface
func newTestGateway() (*Gateway, error) {
	// Create a test node, cookie, and config
	node, err := test.NewNode()
	if err != nil {
		return nil, err
	}

	apiConfig, err := test.NewAPIConfig()
	if err != nil {
		return nil, err
	}

	// Create an address to bind the API to
	addr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/9191")
	if err != nil {
		return nil, err
	}

	listener, err := manet.Listen(addr)
	if err != nil {
		return nil, err
	}

	return NewGateway(node, *test.GetAuthCookie(), listener.NetListener(), *apiConfig)
}

// apiTest is a test case to be run against the api blackbox
type apiTest struct {
	method      string
	path        string
	requestBody string

	expectedResponseCode int
	expectedResponseBody string
}

// apiTests is a slice of apiTest
type apiTests []apiTest

// request issues an http request directly to the blackbox handler
func request(method string, path string, body string) (*http.Response, error) {
	// Create a JSON request to the given endpoint
	req, err := http.NewRequest(method, testURIRoot+path, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}

	// Set headers/auth/cookie
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth("test", "test")
	req.AddCookie(test.GetAuthCookie())

	// Make the request
	resp, err := testHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func runAPITests(t *testing.T, tests apiTests) {
	// Create test repo
	repository, err := test.NewRepository()
	if err != nil {
		t.Fatal(err)
	}

	// Reset repo state
	repository.Reset()
	if err != nil {
		t.Fatal(err)
	}

	// Run each test in serial
	for _, jsonAPITest := range tests {
		runAPITest(t, jsonAPITest)
	}
}

// runTest executes the given test against the blackbox
func runAPITest(t *testing.T, test apiTest) {
	// Make the request
	resp, err := request(test.method, test.path, test.requestBody)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure correct status code
	isEqual := assert.Equal(t, test.expectedResponseCode, resp.StatusCode)
	if !isEqual {
		return
	}

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	// Parse response as JSON
	var responseJSON interface{}
	err = json.Unmarshal(respBody, &responseJSON)
	if err != nil {
		t.Fatal(err)
	}

	// Unless explicity saying any JSON is expected check for equality
	if test.expectedResponseBody != anyResponseJSON {
		var expectedJSON interface{}
		err = json.Unmarshal([]byte(test.expectedResponseBody), &expectedJSON)
		if err != nil {
			t.Fatal(err)
		}

		isEqual = assert.True(t, reflect.DeepEqual(responseJSON, expectedJSON))
		if !isEqual {
			fmt.Println("expected:", test.expectedResponseBody)
			fmt.Println("actual:", string(respBody))
		}
	}
}
