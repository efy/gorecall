package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebInfoHandlerNoQuery(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/webinfo", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := mockApi.WebInfoHandler()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected %d got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestWebInfoHandler(t *testing.T) {
	html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Web info testpage</title>
		</head>
		<body>
		</body>
		</html>
	`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	}))
	defer ts.Close()

	req, err := http.NewRequest("GET", "/api/webinfo?url="+ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := mockApi.WebInfoHandler()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, rr.Code)
	}
}
