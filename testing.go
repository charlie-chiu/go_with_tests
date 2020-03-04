package poker

import "testing"

type StubPlayerStore struct {
	score    map[string]int
	winCalls []string
	league   League
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	return s.score[name]
}
func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}
func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

func AssertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
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
