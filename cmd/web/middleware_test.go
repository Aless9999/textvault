package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"textvault/internal/assert"
	"testing"
)

func TestCommonHeaders(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	commonHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equals(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	expectedValue = "origin-when-cross-origin"
	assert.Equals(t, rs.Header.Get("Referrer-Policy"), expectedValue)

	expectedValue = "nosniff"

	assert.Equals(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)

	expectedValue = "deny"
	assert.Equals(t, rs.Header.Get("X-Frame-Options"), expectedValue)

	expectedValue = "0"

	assert.Equals(t, rs.Header.Get("X-XSS-Protection"), expectedValue)

	expectedValue = "Go"

	assert.Equals(t, rs.Header.Get("Server"), expectedValue)

	assert.Equals(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)

	assert.Equals(t, string(body), "OK")

}
