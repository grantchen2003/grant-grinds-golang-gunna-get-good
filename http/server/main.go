package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// w http.ResponseWriter is an interface that holds a type
// (the concrete type that implements the interface) and
// value (an instance of that concrete type, which can be a
// pointer or a value), so it's passed by value, while
// r *http.Request is a struct passed as a pointer to avoid
// copying large data and allow modifications. Most
// http.ResponseWriter implementations use pointer receivers,
// so w usually holds a pointer, though interfaces can also
// hold non-pointer values.
func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving %s\n", strings.Split(r.RemoteAddr, ":")[0])

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// process request body
	fmt.Println(string(body))

	// Write a response
	// fmt.Fprintf writes a formatted string to a specified
	// destination, such as an http.ResponseWriter, allowing dynamic
	// content to be sent as an HTTP response.
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Hello, World! You've requested: %s\n", r.URL.Path)
}

func main() {
	http.HandleFunc("/", handler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server", err)
	}
}
