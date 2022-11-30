package main

import (
	"fmt"
	"net/http"

	"github.com/rollbar/rollbar-go"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to my webiste</h1>")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", nil)
	rollbar.SetToken("db0eeaf619584284873f3ba1fbc012e1")
	rollbar.SetEnvironment("production")
	rollbar.SetServerRoot("/")
}
