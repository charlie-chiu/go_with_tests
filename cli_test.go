package poker_test

import (
	"bytes"
	"fmt"
	poker "github.com/charlie-chiu/go_with_test"
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

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, blind int) {
	s.alerts = append(s.alerts, scheduledAlert{
		at:    duration,
		blind: blind,
	})
}

func TestCLI(t *testing.T) {
	t.Run("record Chris win from user input", func(t *testing.T) {
		//arrange
		in := strings.NewReader("7\nChris wins\n")
		playerStore := &poker.StubPlayerStore{}
		game := poker.NewGame(dummySpyAlerter, playerStore)
		cli := poker.NewCLI(in, dummyStdOut, game)

		// act
		cli.PlayPoker()

		// assert
		poker.AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record Charlie win from user input", func(t *testing.T) {
		//arrange
		in := strings.NewReader("5\nCharlie wins\n")
		playerStore := &poker.StubPlayerStore{}
		game := poker.NewGame(dummySpyAlerter, playerStore)
		cli := poker.NewCLI(in, dummyStdOut, game)

		// act
		cli.PlayPoker()

		// assert
		poker.AssertPlayerWin(t, playerStore, "Charlie")
	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, playerStore)

		in := strings.NewReader("5\n")
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

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

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		playerStore := &poker.StubPlayerStore{}
		game := poker.NewGame(blindAlerter, playerStore)

		stdOut := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		cli := poker.NewCLI(in, stdOut, game)

		cli.PlayPoker()

		got := stdOut.String()
		want := poker.PlayerPrompt

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

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

func TestGame_Start(t *testing.T) {
	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		game.Start(5)

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
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		game.Start(7)

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
	game := poker.NewGame(dummySpyAlerter, store)
	winner := "Charlie"

	// act
	game.Finish(winner)

	// assert
	poker.AssertPlayerWin(t, store, "Charlie")
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
