package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVisitHandler(t *testing.T) {
	repo := NewInMemoryRepo()
	handler := VisitHandler(repo)

	req := httptest.NewRequest("GET", "/visit", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
}
