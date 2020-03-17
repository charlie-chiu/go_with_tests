package poker

import (
	"testing"
)

type StubPlayerStore struct {
	Score    map[string]int
	WinCalls []string
	League   League
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	return s.Score[name]
}
func (s *StubPlayerStore) RecordWin(name string) {
	s.WinCalls = append(s.WinCalls, name)
}
func (s *StubPlayerStore) GetLeague() League {
	return s.League
}

func AssertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
	t.Helper()

	got := len(store.WinCalls)
	want := 1
	if got != want {
		t.Fatalf("got %d calls to RecordWin want %d", got, want)
	}

	if store.WinCalls[0] != winner {
		t.Errorf("didn't store correct winner, got %s want %s", store.WinCalls[0], winner)
	}
}
