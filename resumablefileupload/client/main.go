package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fileData := getFileData()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter command (upload, download, pause, resume, abort, exit): ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		switch input {
		case "upload":
			fmt.Println("Command received: upload")
			uploadFile(fileData)

		case "download":
			fmt.Println("Command received: download")

		case "pause":
			fmt.Println("Command received: pause")

		case "resume":
			fmt.Println("Command received: resume")

		case "abort":
			fmt.Println("Command received: abort")

		case "exit":
			return

		default:
			fmt.Println("Unknown command. Please enter start, stop, or exit.")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
