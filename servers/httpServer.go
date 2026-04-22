package servers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type response struct {
	Message string `json:"message"`
	Error   int    `json:"error"`
}

func home(w http.ResponseWriter, req *http.Request) {
	homeResponse := response{
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

func Http(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Http serve on port: 2000")

	http.HandleFunc("/home", home)
	log.Fatal(http.ListenAndServe(":2000", nil))

}
