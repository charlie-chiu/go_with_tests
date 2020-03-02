package main

import (
	"strings"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("/league from a reader", func(t *testing.T) {
		database := strings.NewReader(`[
		{"Name": "Cleo", "Wins": 10},	
		{"Name": "Charlie", "Wins": 999}]`)

		store := FileSystemPlayerStore{database}

		got := store.GetLeague()

		want := []Player{
			{"Cleo", 10},
			{"Charlie", 999},
		}
		assertLeague(t, got, want)

		//read again
		got = store.GetLeague()
		assertLeague(t, got, want)

	})
}
