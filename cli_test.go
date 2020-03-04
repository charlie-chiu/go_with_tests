package poker_test

import (
	poker "github.com/charlie-chiu/go_with_test"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	t.Run("record Chris win from user input", func(t *testing.T) {
		//arrange
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}
		cli := poker.NewCLI(playerStore, in)

		// act
		cli.PlayPoker()

		// assert
		poker.AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record Charlie win from user input", func(t *testing.T) {
		//arrange
		in := strings.NewReader("Charlie wins\n")
		playerStore := &poker.StubPlayerStore{}
		cli := poker.NewCLI(playerStore, in)

		// act
		cli.PlayPoker()

		// assert
		poker.AssertPlayerWin(t, playerStore, "Charlie")
	})

}
