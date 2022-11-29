package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to my webiste</h1>")

	fmt.Printf("<h1>Enter your name: </h1>")
	var name string
	fmt.Scan(&name)

	fmt.Printf("<h2>Hello, %v</h2>", name)
}

func main() {
	http.HandleFunc("/", handlerFunc)
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", nil)
}
