package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChatHandler(t *testing.T) {
	// Create a request to pass to our handler with a sample input
	req, err := http.NewRequest("GET", "/chat?input=Hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(chatHandler)

	// Call the handler with our request and recorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "ChatGPT: "
	if rr.Body.String()[:len(expected)] != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
