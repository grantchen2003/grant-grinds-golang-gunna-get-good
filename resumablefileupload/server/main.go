package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func initiateUploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving %s\n", strings.Split(r.RemoteAddr, ":")[0])
	uploadId := "123"
	fmt.Fprint(w, uploadId)
}

func main() {
	http.HandleFunc("/upload/initiate", initiateUploadHandler)

	log.Println("Starting server on port :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
