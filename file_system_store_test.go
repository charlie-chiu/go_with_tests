package poker_test

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

	poker "github.com/charlie-chiu/go_with_test"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("/league from a reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},	
		{"Name": "Charlie", "Wins": 999}]`)
		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		if err != nil {
			log.Fatalf("problem creating file system player store, %v", err)
		}
		got := store.GetLeague()

		want := poker.League{
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

		store, err := poker.NewFileSystemPlayerStore(database)
		if err != nil {
			log.Fatalf("problem creating File system player store, %v", err)
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

		store, err := poker.NewFileSystemPlayerStore(database)
		if err != nil {
			log.Fatalf("problem creating File system player store, %v", err)
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

		store, err := poker.NewFileSystemPlayerStore(database)
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

		_, err := poker.NewFileSystemPlayerStore(database)

		AssertNoError(t, err)
	})
	t.Run("league sorted from highest to lowest", func(t *testing.T) {
		file, cleanDatabase := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},	
		{"Name": "Charlie", "Wins": 999}]`)
		defer cleanDatabase()
		store, _ := poker.NewFileSystemPlayerStore(file)
		got := store.GetLeague()
		want := poker.League{
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

func AssertLeague(t *testing.T, got, want poker.League) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func AssertScoreEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Got %d want %d", got, want)
	}
}

func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("didn't expect an error but got one, %v", err)
	}
}
