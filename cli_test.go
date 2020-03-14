package poker_test

import (
	"bytes"
	"fmt"
	poker "github.com/charlie-chiu/go_with_test"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

var dummySpyAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

type scheduledAlert struct {
	at    time.Duration
	blind int
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.blind, s.at)
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	s.alerts = append(s.alerts, scheduledAlert{
		at:    duration,
		blind: amount,
	})
}

func TestCLI(t *testing.T) {
	t.Run("it prompts the user to enter the number of players and start", func(t *testing.T) {

		game := &poker.GameSpy{}

		stdOut := &bytes.Buffer{}
		in := strings.NewReader("7\nsomeone wins")
		cli := poker.NewCLI(in, stdOut, game)

		cli.PlayPoker()

		got := stdOut.String()
		want := poker.PlayerPrompt

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		assertGameStartWith(t, game, 7)
	})

	t.Run("start game with 6 players and finish game with \"Chris\" as winner", func(t *testing.T) {
		//arrange
		in := strings.NewReader("6\nChris wins\n")
		game := &poker.GameSpy{}
		cli := poker.NewCLI(in, dummyStdOut, game)

		// act
		cli.PlayPoker()

		// assert
		assertFinishCalledWith(t, game, "Chris")
		assertGameStartWith(t, game, 6)
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		g := &poker.GameSpy{StartCalled: false}
		out := &bytes.Buffer{}
		in := strings.NewReader("someString")

		cli := poker.NewCLI(in, out, g)
		cli.PlayPoker()

		assertMessagesSentToUser(t, out, poker.PlayerPrompt, poker.BadPlayerInputErrorMsg)

		if g.StartCalled {
			t.Errorf("game should not have started")
		}
	})

	t.Run("it prints an error when entered a string not 'someone wins' and game should not finish", func(t *testing.T) {
		stdOut := &bytes.Buffer{}
		stdIn := strings.NewReader("1\nHello, Worlda\n")
		game := &poker.GameSpy{FinishCalled: false}
		cli := poker.NewCLI(stdIn, stdOut, game)

		cli.PlayPoker()
		assertMessagesSentToUser(t, stdOut, poker.PlayerPrompt, poker.BadWinnerInputErrorMsg)

		if game.FinishCalled == true {
			t.Errorf("game should not finished")
		}
	})
}

func TestGame_Start(t *testing.T) {
	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewTexasHoldem(blindAlerter, dummyPlayerStore)

		game.Start(5, ioutil.Discard)

		cases := []scheduledAlert{
			// 5 players, blind increase every 10 minutes
			{0 * time.Minute, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduleAlert(t, got, want)
			})
		}
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewTexasHoldem(blindAlerter, dummyPlayerStore)

		game.Start(7, ioutil.Discard)

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduleAlert(t, got, want)
			})
		}
	})
}

func TestGame_Finish(t *testing.T) {
	//arrange
	store := &poker.StubPlayerStore{}
	game := poker.NewTexasHoldem(dummySpyAlerter, store)
	winner := "Charlie"

	// act
	game.Finish(winner)

	// assert
	poker.AssertPlayerWin(t, store, "Charlie")
}

func assertFinishCalledWith(t *testing.T, game *poker.GameSpy, winner string) {
	t.Helper()
	got := game.FinishedWith
	if got != winner {
		t.Errorf("expected %q wins, got %q", winner, got)
	}
}

func assertGameStartWith(t *testing.T, game *poker.GameSpy, numberOfPlayer int) {
	t.Helper()
	if game.StartedWith != numberOfPlayer {
		t.Errorf("want Start Called with 6 but got %d", game.StartedWith)
	}
}

func assertMessagesSentToUser(t *testing.T, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}

func assertScheduleAlert(t *testing.T, got, want scheduledAlert) {
	blindGot := got.blind
	if blindGot != want.blind {
		t.Errorf("got amount %d, want %d", blindGot, want.blind)
	}

	gotScheduledTime := got.at
	if gotScheduledTime != want.at {
		t.Errorf("got scheduled time of %v, want %v", gotScheduledTime, want.at)
	}
}
