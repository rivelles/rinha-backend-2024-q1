package api

import (
	"fmt"
	"net/http"
)

func hello(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("Received a request on hello!")
}

func run() {
	http.HandleFunc("/hello", hello)

	http.ListenAndServe(":9091", nil)
}
