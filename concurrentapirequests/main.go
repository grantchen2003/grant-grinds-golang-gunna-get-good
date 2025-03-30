package main

import (
	"fmt"
	"sync"
	"time"
)

func generateUserName() string {
	return "Alice"
}

func fetchUserPost(userName string) string {
	fmt.Printf("Fetching post for %s\n", userName)
	time.Sleep(5 * time.Second) // simulate api request

	return fmt.Sprintf("This is a post from %s", userName)
}

func fetchUserFriends(userName string) []string {
	fmt.Printf("Fetching friends for %s\n", userName)
	time.Sleep(4 * time.Second) // simulate api request

	return []string{"Bob", "Carlos", "Doug"}

}

func makeApiRequest(apiUrl string) string {
	fmt.Printf("Fetching api: %s\n", apiUrl)
	time.Sleep(2 * time.Second) // simulate api request
	return fmt.Sprintf("Dummy data from %s", apiUrl)
}

func main() {
	startTime := time.Now()

	userName := generateUserName()

	userPostChannel := make(chan string)
	userFriendsChannel := make(chan []string)
	apiChannel := make(chan string)

	go func() {
		userPostChannel <- fetchUserPost(userName)
		// We close userPostChannel because only one value
		// will be sent to this channel – fetchUserPost(userName)
		// returns a single string, so there’s no need to keep
		// the channel open. Closing the channel signals completion,
		// ensuring the main function doesn't block indefinitely
		// when it reads from the channel (<-userPostChannel),
		// expecting more data.
		close(userPostChannel)
	}()

	go func() {
		userFriendsChannel <- fetchUserFriends(userName)
		close(userFriendsChannel)
	}()

	wg := &sync.WaitGroup{}

	urls := []string{
		"http://example.1",
		"http://example.2",
		"http://example.3",
	}

	wg.Add(len(urls))
	for _, url := range urls {
		go func(url string) {
			defer wg.Done()
			apiChannel <- makeApiRequest(url)
		}(url)
	}

	// wg.Wait() must be called in a separate goroutine to
	// prevent a deadlock. Writing to a channel is blocking
	// until the receiver reads from it. So if wg.Wait() is
	// called directly, the main goroutine waits for all
	// API requests to finish, but the API goroutines are
	// blocked, waiting for the main goroutine to read from
	// the channel. This causes a circular dependency since
	// the main goroutine can't start receiving data until
	// wg.Wait() completes, but the API goroutines can't
	// send data until the receiver is ready which is after
	// wg.Wait() completes. By calling wg.Wait() asynchronously,
	// we allow the main goroutine to start receiving data from
	// the channel while the API requests are still running,
	// avoiding the deadlock.
	go func() {
		wg.Wait()
		// we close here knowing all the goroutines are done
		// writing to the channel. We close so the for loop below
		// that reads from the channel knows when to break
		close(apiChannel)
	}()

	for data := range apiChannel {
		fmt.Println(data)
	}

	fmt.Println("User Post:", <-userPostChannel)
	fmt.Println("User Friends:", <-userFriendsChannel)

	// prints 5 seconds
	fmt.Printf("Total time taken: %v\n", time.Since(startTime))
}
