package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"resumablefileuploadclient/utils"
)

func initializeUploadSession() string {
	resp, err := http.Post("http://localhost:8080/upload/initiate", "text/plain", nil)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: Received status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	uploadId := string(body)

	return uploadId
}

func uploadFileChunk(fileChunk []byte, uploadId string) error {
	resp, err := http.Post(fmt.Sprintf("http://localhost:8080/upload/%s/chunk", uploadId), "text/plain", bytes.NewBuffer(fileChunk))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: Received status code %d", resp.StatusCode)
	}

	return nil
}

func uploadFile(fileData []byte) {
	uploadId := initializeUploadSession()

	fileChunks := utils.Chunkify(fileData, 100)
	for _, fc := range fileChunks {
		go uploadFileChunk(fc, uploadId)
	}
}
