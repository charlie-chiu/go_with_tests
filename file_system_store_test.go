package main

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("/league from a reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},	
		{"Name": "Charlie", "Wins": 999}]`)
		defer cleanDatabase()

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
	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},	
		{"Name": "Charlie", "Wins": 999}]`)
		defer cleanDatabase()

		store := FileSystemPlayerStore{database}

		got := store.GetPlayerScore("Charlie")
		want := 999

		assertScoreEquals(t, got, want)
	})
	t.Run("store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},	
		{"Name": "Charlie", "Wins": 999}]`)
		defer cleanDatabase()

		store := FileSystemPlayerStore{database}

		store.RecordWin("Cleo")

		got := store.GetPlayerScore("Cleo")
		want := 11
		assertScoreEquals(t, got, want)
	})
	t.Run("store wins for new players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},	
		{"Name": "Charlie", "Wins": 999}]`)
		defer cleanDatabase()

		store := FileSystemPlayerStore{database}

		store.RecordWin("Frog")

		got := store.GetPlayerScore("Frog")
		want := 1
		assertScoreEquals(t, got, want)
	})

}

func assertScoreEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Got %d want %d", got, want)
	}
}

func createTempFile(t *testing.T, initialData string) (io.ReadWriteSeeker, func()) {
	t.Helper()

	tempFile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file, %v", err)
	}

	tempFile.Write([]byte(initialData))

	removeFile := func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}

	return tempFile, removeFile
}
