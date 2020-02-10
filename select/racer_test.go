package racer

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRacer(t *testing.T) {
	t.Run("test Racer return faster http request", func(t *testing.T) {
		// using real http requests
		//slowURL := "http://www.facebook.com"
		//slowURL := "http://127.0.0.1"
		//fastURL := "http://www.google.com"

		// replace with mock HTTP server
		slowServer := makeDelayedServer(20 * time.Millisecond)
		slowURL := slowServer.URL
		defer slowServer.Close()
		fastServer := makeDelayedServer(0 * time.Millisecond)
		fastURL := fastServer.URL
		defer fastServer.Close()
		// httptest.NewServer will find an open port to listen on

		want := fastURL
		got, err := Racer(slowURL, fastURL)

		if err != nil {
			t.Fatalf("did not expect an error but got one %v", err)
		}

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("return an error if a server doesn't respond within 10s", func(t *testing.T) {
		//serverA :=makeDelayedServer(11 * time.Second)
		//defer serverA.Close()
		//serverB := makeDelayedServer(12 * time.Second)
		//defer serverB.Close()

		//we doesn't really care about which server is faster.
		s := makeDelayedServer(25 * time.Millisecond)
		defer s.Close()

		timeout := 10 * time.Millisecond
		_, err := ConfigurableRacer(s.URL, s.URL, timeout)

		if err == nil {
			t.Errorf("expected an error, but didn't get one")
		}

	})
}

func makeDelayedServer(d time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(d)
		w.WriteHeader(http.StatusOK)
	}))
}
