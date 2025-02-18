package starr

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMockData allows generic testing of http inputs and outputs.
// This is used by the submodule tests.
type TestMockData struct {
	Name            string      // A name for the test.
	ExpectedPath    string      // The path expected in the request ie. /api/v1/thing
	ExpectedRequest string      // The request body (json) expected from the caller.
	ExpectedMethod  string      // The request method (GET/POST) expected from the caller.
	ResponseStatus  int         // This is the status that gets returned the caller.
	ResponseBody    string      // The (json) response body returned to caller.
	WithRequest     interface{} // Caller's request.
	WithResponse    interface{} // Caller's response.
	WithError       error       // Caller's response.
}

const (
	// Error body for 401 response.
	BodyUnauthorized = `{"error": "Unauthorized"}`
	// Error body for 404 response.
	BodyNotFound = `{"message": "NotFound"}`
	// Error body for 405 response.
	BodyMethodNotAllowed = `{"message": "MethodNotAllowed"}`
)

// GetMockServer is used in all the http tests.
func (test *TestMockData) GetMockServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		assert.EqualValues(t, test.ExpectedPath, req.URL.String())
		writer.WriteHeader(test.ResponseStatus)

		assert.EqualValues(t, req.Method, test.ExpectedMethod)

		body, err := ioutil.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.EqualValues(t, test.ExpectedRequest, string(body))

		_, err = writer.Write([]byte(test.ResponseBody))
		assert.NoError(t, err)
	}))
}
