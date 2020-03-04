package poker

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("/league from a reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},	
		{"Name": "Charlie", "Wins": 999}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}
		got := store.GetLeague()

		want := []Player{
			{"Charlie", 999},
			{"Cleo", 10},
		}
		AssertLeague(t, got, want)

		//read again
		got = store.GetLeague()
		AssertLeague(t, got, want)
	})
	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},	
		{"Name": "Charlie", "Wins": 999}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}
		got := store.GetPlayerScore("Charlie")
		want := 999

		AssertScoreEquals(t, got, want)
	})
	t.Run("store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},	
		{"Name": "Charlie", "Wins": 999}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}

		store.RecordWin("Cleo")

		got := store.GetPlayerScore("Cleo")
		want := 11
		AssertScoreEquals(t, got, want)
	})
	t.Run("store wins for new players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},	
		{"Name": "Charlie", "Wins": 999}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}

		store.RecordWin("Frog")

		got := store.GetPlayerScore("Frog")
		want := 1
		AssertScoreEquals(t, got, want)
	})
	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)

		AssertNoError(t, err)
	})
	t.Run("league sorted from highest to lowest", func(t *testing.T) {
		file, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},	
		{"Name": "Charlie", "Wins": 999}]`)
		defer cleanDatabase()
		store, _ := NewFileSystemPlayerStore(file)
		got := store.GetLeague()
		want := League{
			{"Charlie", 999},
			{"Cleo", 10},
		}

		AssertLeague(t, got, want)
	})

}

func createTempFile(t *testing.T, initialData string) (*os.File, func()) {
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
