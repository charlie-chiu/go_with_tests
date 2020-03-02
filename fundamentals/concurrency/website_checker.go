package concurrency

type WebsiteChecker func(string) bool
type result struct {
	url         string
	checkResult bool
}

// first try...
func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)

	for _, url := range urls {
		results[url] = wc(url)
	}

	return results
}

// use goroutine, about 100x faster!
func CheckWebsites2(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)
	resultChannel := make(chan result)

	for _, url := range urls {

		go func(u string) {
			// send statement with <- operator
			resultChannel <- result{u, wc(u)}
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		// receive expression uses the <- operator, too
		result := <-resultChannel
		results[result.url] = result.checkResult
	}

	return results
}
