package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// Create new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// Logs may appear twice because a single file write
				// can trigger multiple filesystem events (like WRITE
				// and RENAME), especially when saving via text editors.
				fmt.Println("Event:", event)

				if event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("Created:", event.Name)
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("Modified:", event.Name)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Println("Removed:", event.Name)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					fmt.Println("Renamed:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("Error:", err)
			}
		}
	}()

	// Add a directory to watch
	if watcher.Add(""); err != nil {
		log.Fatal(err)
	}

	// Block forever
	select {}
}
