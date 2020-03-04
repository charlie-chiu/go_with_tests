package poker

import (
	"net/http/httptest"
	"reflect"
	"testing"
)

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

func AssertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	contentType := response.Result().Header.Get("content-type")
	if contentType != want {
		t.Errorf("response did not have content-type of %s, got %v", want, contentType)
	}
}

func AssertLeague(t *testing.T, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func AssertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got status %d, want %d", got, want)
	}
}

func AssertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("didn't expect an error but got one, %v", err)
	}
}

func AssertScoreEquals(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Got %d want %d", got, want)
	}
}
