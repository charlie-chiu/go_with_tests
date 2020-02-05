package concurrency

import (
	"reflect"
	"testing"
	"time"
)

func mockWebsiteChecker(url string) bool {
	//which one is better?
	//return !(url == "waat://hello.world")
	if url == "waat://hello.world" {
		return false
	}

	return true
}

func slowWebsiteChecker(url string) bool {
	time.Sleep(20 * time.Millisecond)
	return true
}

func TestCheckWebsites(t *testing.T) {
	websites := []string{
		"http://google.com",
		"https://medium.com/@charliecc",
		"waat://hello.world",
	}

	wants := map[string]bool{
		"http://google.com":             true,
		"https://medium.com/@charliecc": true,
		"waat://hello.world":            false,
		//"waat://hello.world":            true,
	}

	got := CheckWebsites(mockWebsiteChecker, websites)

	if !reflect.DeepEqual(got, wants) {
		t.Fatalf("Wanted %v, got %v", wants, got)
	}
}

func BenchmarkCheckWebsites(b *testing.B) {
	urls := make([]string, 100)
	for i := 0; i < len(urls); i++ {
		urls[i] = "some url"
	}


	for i := 0; i < b.N; i++ {
		//CheckWebsites(mockWebsiteChecker, urls)
		CheckWebsites(slowWebsiteChecker, urls)
	}
}
