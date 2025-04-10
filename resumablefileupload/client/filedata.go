package main

import "math/rand"

func getFileData() []byte {
	fileSizeBytes := 10000

	fileData := make([]byte, fileSizeBytes)

	for j := range fileSizeBytes {
		fileData[j] = byte(rand.Intn(256))
	}

	return fileData
}
