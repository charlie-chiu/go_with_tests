package poker

import (
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	t.Run("record Chris win from user input", func(t *testing.T) {
		//arrange
		in := strings.NewReader("Chris wins\n")
		playerStore := &StubPlayerStore{}
		cli := &CLI{playerStore, in}

		// act
		cli.PlayPoker()

		// assert
		assertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record Charlie win from user input", func(t *testing.T) {
		//arrange
		in := strings.NewReader("Charlie wins\n")
		playerStore := &StubPlayerStore{}
		cli := &CLI{playerStore, in}

		// act
		cli.PlayPoker()

		// assert
		assertPlayerWin(t, playerStore, "Charlie")
	})

}

func assertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
	t.Helper()

	got := len(store.winCalls)
	want := 1
	if got != want {
		t.Fatalf("got %d calls to RecordWin want %d", got, want)
	}

	if store.winCalls[0] != winner {
		t.Errorf("didn't store correct winner, got %s want %s", store.winCalls[0], winner)
	}
}
