package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestServeHome contains sub-tests for the serveHome handler.
func TestServeHome(t *testing.T) {

	// Sub-test for the "happy path": a valid GET request.
	t.Run("it returns 200 OK for GET requests", func(t *testing.T) {
		// 1. Create a new HTTP request for our handler.
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		// 2. Create a ResponseRecorder to record the handler's response.
		rec := httptest.NewRecorder()

		// 3. Call the handler function directly.
		serveHome(rec, req)

		// 4. Check the result.
		res := rec.Result()
		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", res.Status)
		}
	})

	// Sub-test for the "sad path": an invalid POST request.
	t.Run("it returns 405 Method Not Allowed for non-GET requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/", nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		rec := httptest.NewRecorder()
		serveHome(rec, req)
		res := rec.Result()
		if res.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("expected status Method Not Allowed; got %v", res.Status)
		}
	})
}

func serveHome(rec *httptest.ResponseRecorder, req *http.Request) {
	panic("unimplemented")
}
