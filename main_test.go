package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList["moscow"])},
	}

	for _, v := range requests {
		reqStr := fmt.Sprintf("/cafe?city=moscow&count=%d", v.count)
		req := httptest.NewRequest("GET", reqStr, nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)

		respLen := 0
		if response.Body.String() != "" {
			respLen = len(strings.Split(strings.TrimSpace(response.Body.String()), ","))
		}

		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, v.want, respLen)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		searchWord string
		count      int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		reqStr := fmt.Sprintf("/cafe?city=moscow&search=%s", v.searchWord)
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", reqStr, nil)

		handler.ServeHTTP(response, req)

		respSlice := strings.Split(strings.TrimSpace(response.Body.String()), ",")

		for _, s := range respSlice {
			strings.Contains(strings.ToLower(s), strings.ToLower(v.searchWord))
		}

		respLen := 0
		if response.Body.String() != "" {
			respLen = len(strings.Split(strings.TrimSpace(response.Body.String()), ","))
		}

		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, v.count, respLen)
	}

}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}
