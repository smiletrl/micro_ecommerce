package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

// APITestCase represents the data needed to describe an API test case.
type APITestCase struct {
	Name         string
	Method       string
	URL          string
	Context      map[string]interface{}
	Body         string
	WantStatus   int
	WantResponse string
}

// Endpoint tests an HTTP endpoint using the given AuthAPITestCase spec.
func Endpoint(t *testing.T, e *echo.Echo, tc APITestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		req, _ := http.NewRequest(tc.Method, tc.URL, bytes.NewBufferString(tc.Body))

		res := httptest.NewRecorder()
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}

		e.ServeHTTP(res, req)

		assert.Equal(t, tc.WantStatus, res.Code, "status mismatch")
		pattern := strings.Trim(tc.WantResponse, "*")
		// compare string firstly
		if pattern != tc.WantResponse {
			assert.Contains(t, res.Body.String(), pattern, "response string mismatch")
		} else {
			// compare json
			assert.JSONEq(t, tc.WantResponse, res.Body.String(), "response json mismatch")
		}
	})
}
