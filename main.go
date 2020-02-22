package main

import (
	"github.com/charlie-chiu/go_with_test/application"
	"log"
	"net/http"
)

type InMemoryPlayerStore struct {
	store map[string]int
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.store[name]++
}

func main() {
	handler := &application.PlayerServer{
		&InMemoryPlayerStore{store: map[string]int{}},
	}

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("could not listen on port 8080, %v", err)
	}

}
