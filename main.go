package main

import (
	"github.com/charlie-chiu/go_with_test/application"
	"log"
	"net/http"
)

type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return 123345345
}

func main() {
	handler := &application.PlayerServer{
		&InMemoryPlayerStore{},
	}

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("could not listen on port 8080, %v", err)
	}

}
