package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTemperature(t *testing.T) {
	t.Run("should run successfully", func(t *testing.T) {
		r, err := http.NewRequest("GET", "/temperature?cep=01001000", nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		getTemperature(w, r)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code 200, got %d", resp.StatusCode)
		}

		println(w.Body.String())
	})

	t.Run("should return error when cep is invalid", func(t *testing.T) {
		r, err := http.NewRequest("GET", "/temperature?cep=123", nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		getTemperature(w, r)

		resp := w.Result()
		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("expected status code 422, got %d", resp.StatusCode)
		}

		println(w.Body.String())
	})

	t.Run("should return error when cep is not found", func(t *testing.T) {
		r, err := http.NewRequest("GET", "/temperature?cep=99999999", nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		getTemperature(w, r)

		resp := w.Result()
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status code 404, got %d", resp.StatusCode)
		}

		println(w.Body.String())
	})
}
