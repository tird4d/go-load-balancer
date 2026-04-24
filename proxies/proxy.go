package proxies

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var backendURL = "http://localhost:2000"

var client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	},
}

var hopByHope = map[string]bool{
	"Connection":          true,
    "Keep-Alive":          true,
    "Transfer-Encoding":   true,
    "Upgrade":             true,
    "Proxy-Authorization": true,
    "Proxy-Authenticate":  true,
    "Te":                  true,
    "Trailers":            true,
	"X-Dangerous-Header":	true,
}



func proxy(w http.ResponseWriter, req *http.Request) {

	serverReq, err := http.NewRequest(req.Method, backendURL+req.RequestURI, nil)
	if err != nil {
		fmt.Println("url error")
	}
	headers := map[string][]string{}

	for i, v := range req.Header {
		fmt.Println(i, strings.Join(v, ","))
		if hopByHope[i] {
			continue
		}

		headers[i] = v
	}
	headers["X-Forwarded-For"] = []string{req.RemoteAddr}

	serverReq.Header = headers
	serverRes, err := client.Do(serverReq)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "server not respond error", http.StatusInternalServerError)
		return
	}

	for i, v := range serverRes.Header {
		fmt.Println(i, strings.Join(v, ","))

		w.Header().Set(i, strings.Join(v, ","))

	}
	w.Header().Set("X-Proxy", "Cached")
	w.WriteHeader(serverRes.StatusCode)
	io.Copy(w, serverRes.Body)
}

func Http(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Http serve on port: 2001")

	proxyMux := http.NewServeMux()
	proxyMux.HandleFunc("/", proxy)

	log.Fatal(http.ListenAndServe(":2001", proxyMux))

}
