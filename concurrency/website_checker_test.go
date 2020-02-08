package concurrency

/*
	run test with...
	go test
	go test -bench . (same as go test -bench=.)
	go test -race
		using race detector help debug
*/

import (
	"reflect"
	"testing"
	"time"
)

func mockWebsiteChecker(url string) bool {
	if url == "waat://hello.world" {
		return false
	}
	return true
}

func TestCheckWebsites(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.charliecc.com",
		"waat://hello.world",
	}

	want := map[string]bool{
		"http://google.com":         true,
		"http://blog.charliecc.com": true,
		"waat://hello.world":        false,
	}

	got := CheckWebsites(mockWebsiteChecker, websites)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Wanted %v, got %v", want, got)
	}
}

func slowStubWebsiteChecker(_ string) bool {
	time.Sleep(20 * time.Millisecond)
	return true
}

func BenchmarkCheckWebsites(b *testing.B) {
	websites := make([]string, 100)
	for i := 0; i < len(websites); i++ {
		websites[i] = "a url"
	}

	b.Run("without goroutine", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			CheckWebsites(slowStubWebsiteChecker, websites)
		}
	})

	b.Run("with goroutine", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			CheckWebsites2(slowStubWebsiteChecker, websites)
		}
	})

}
