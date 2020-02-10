package racer

import (
	"fmt"
	"net/http"
	"time"
)

const tenSecondTimeout = 10 * time.Second

func Racer(url1, url2 string) (winner string, err error) {
	return ConfigurableRacer(url1, url2, tenSecondTimeout)
}

func ConfigurableRacer(url1, url2 string, timeout time.Duration) (winner string, err error) {
	select {
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for %s and %s", url1, url2)
	case <-ping(url1):
		return url1, nil
	case <-ping(url2):
		return url2, nil
	}
}

func ping(url string) chan struct{} {
	ch := make(chan struct{})
	go func() {
		http.Get(url)
		close(ch)
	}()
	return ch
}

// step 1
func RacerV1(url1, url2 string) (winner string) {
	duration1 := measureResponseTime(url1)
	duration2 := measureResponseTime(url2)

	if duration1 > duration2 {
		return url2
	}

	return url1
}

func measureResponseTime(url string) time.Duration {
	start := time.Now()
	http.Get(url)

	return time.Since(start)
}
