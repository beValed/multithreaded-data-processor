package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetData_Success(t *testing.T) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	}

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handlerFunc(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.Equal(t, "Success", rr.Body.String())
}

func TestGetData_NotFound(t *testing.T) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handlerFunc(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestPostData_Success(t *testing.T) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Data created successfully"))
	}

	req, err := http.NewRequest("POST", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handlerFunc(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	assert.Equal(t, "Data created successfully", rr.Body.String())
}
