package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"sync"
)

func getResponseSize(url string) (string, int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return url, 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return url, 0, err
	}

	return url, len(body), nil
}

func main() {
	// Parse CLI arguments
	urls := os.Args[1:]

	// Initiate
	var wg sync.WaitGroup

	// Get responses
	ch := make(chan struct {
		url  string
		size int
		err  error
	})

	// Make requests concurrently
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			url, size, err := getResponseSize(url)
			ch <- struct {
				url  string
				size int
				err  error
			}{url, size, err}
		}(url)
	}

	// Finish requests
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Get responses
	var results []struct {
		url  string
		size int
		err  error
	}
	for r := range ch {
		results = append(results, r)
	}

	// Sort responses
	sort.Slice(results, func(i, j int) bool {
		return results[i].size < results[j].size
	})

	// Output results
	for _, r := range results {
		if r.err == nil {
			fmt.Printf("%s %d\n", r.url, r.size)
		} else {
			fmt.Printf("%s error: %s\n", r.url, r.err)
		}
	}
}
