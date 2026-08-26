package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"snippetbox.macnigor.net/internal/assert"
	"testing"
)

func TestPing(t *testing.T) {
	rr := httptest.NewRecorder() //создаем объект который будет записывать ответ вместо настоящега writer

	r, err := http.NewRequest(http.MethodGet, "/", nil) //создаем запрос фиктивный с методом get

	if err != nil {
		t.Fatal(err)
	}

	ping(rr, r)

	rs := rr.Result() //просмотр результата

	assert.Equals(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)

	assert.Equals(t, string(body), "OK")

}
