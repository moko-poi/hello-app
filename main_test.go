package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()

	cases := []struct {
		method     string
		target     string
		body       io.Reader
		wantStatus int
		wantBody   string
	}{
		{method: "GET", target: "/hello", wantStatus: http.StatusOK, wantBody: "Hello, world!"},
	}

	for _, c := range cases {
		req := httptest.NewRequest(c.method, c.target, c.body)
		w := httptest.NewRecorder()
		helloHandler(w, req)

		resp := w.Result()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		err = resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != c.wantStatus {
			t.Errorf("The respons code sholud be %d, but was %d.", c.wantStatus, resp.StatusCode)
		}
		if string(body) != c.wantBody {
			t.Errorf("The respons body should be %s but was %s.", c.wantBody, string(body))
		}
	}
}
