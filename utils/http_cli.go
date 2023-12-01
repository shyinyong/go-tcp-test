package main

import (
	"fmt"
	"net/http"
	"sync"
)

type CallbackFunc func(url string, success bool)

func testHTTPGetAsync(url string, callback CallbackFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		callback(url, false)
		return
	}
	defer resp.Body.Close()

	success := resp.StatusCode == http.StatusOK
	callback(url, success)
}

func main() {
	urls := []string{"https://google.com", "https://facebook.com", "https://twitter.com", "https://bing.com"}

	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go testHTTPGetAsync(url, func(url string, success bool) {
			fmt.Printf("URL: %s, Success: %t\n", url, success)
		}, &wg)
	}

	// 等待所有 Goroutines 完成
	wg.Wait()
}
