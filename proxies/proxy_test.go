package proxies

import (
	"encoding/json"
	"fmt"

	"net/http"
	"net/http/httptest"
	"testing"
)


type response struct {
	Message string `json:"message"`
	Error   int    `json:"error"`
}

func TestProxy(t *testing.T) {

	handler := func(w http.ResponseWriter, r *http.Request){
		homeResponse := response {
		Message: "Hello from world",
		Error:   0,
	}

		jdata, err := json.Marshal(homeResponse)
		if err != nil {
			fmt.Printf("json conversion error %s", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Custom-Server-Header", "secret")
		w.Write(jdata)
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	backendURL = ts.URL // point the proxy at the test server

	req := httptest.NewRequest("GET", "/home", nil)
	req.Header.Add("X-SafeHeader", "true")
	req.Header.Add("X-Dangerous-Header", "danger")
	w := httptest.NewRecorder()
	proxy(w, req)

	result := w.Result()

	if result.StatusCode != 200 {
		t.Errorf("status code should be 200 but is: %d", result.StatusCode)
	}

	if result.Header.Get("X-Proxy") != "Cached" {
		t.Errorf("header X-Proxy not found")
	}



}
