package poker

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	dummyGame = &GameSpy{}
)

func TestWebSocket(t *testing.T) {
	t.Run("start a game with 3 players and declare Charlie the winner", func(t *testing.T) {
		wantedBlindAlert := "Blind is 100"

		game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}
		playerServer := mustMakePlayerServer(t, &StubPlayerStore{}, game)
		server := httptest.NewServer(playerServer)
		ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

		defer server.Close()
		defer ws.Close()

		numberOfPlayers := 3
		winner := "Charlie"
		writeWSMessage(t, ws, strconv.Itoa(numberOfPlayers))
		writeWSMessage(t, ws, winner)

		const tenMS = 10 * time.Millisecond
		//time.Sleep(tenMS)
		assertGameStartedWith(t, game, numberOfPlayers)
		assertFinishCalledWith(t, game, winner)

		within(t, tenMS, func() { assertWebsocketGotMsg(t, ws, wantedBlindAlert) })
	})
}

func assertFinishCalledWith(t *testing.T, game *GameSpy, winner string) {
	t.Helper()

	retryTime := 500 * time.Millisecond
	passed := retryUntil(retryTime, func() bool {
		return game.FinishCalledWith == winner
	})

	if !passed {
		t.Errorf("expected winner is %q but got %q", winner, game.FinishCalledWith)
	}
}

func assertGameStartedWith(t *testing.T, game *GameSpy, numberOfPlayers int) {
	t.Helper()
	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.StartCalledWith == numberOfPlayers
	})

	if !passed {
		t.Errorf("wanted Start called with %d but got %d", numberOfPlayers, game.StartCalledWith)
	}
}

func retryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}

	return false
}

func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, want string) {
	t.Helper()
	_, gotBlindAlert, _ := ws.ReadMessage()
	if string(gotBlindAlert) != want {
		t.Errorf("got blind alert %q, want %q", string(gotBlindAlert), want)
	}
}

func TestGetGame(t *testing.T) {
	t.Run("GET /game return 200", func(t *testing.T) {
		server := mustMakePlayerServer(t, &StubPlayerStore{}, dummyGame)

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusOK)
	})
}

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{score: map[string]int{
		"Pepper": 20,
		"Floyd":  10,
	}}
	server := mustMakePlayerServer(t, &store, dummyGame)

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertResponseBody(t, response.Body.String(), "20")
		AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertResponseBody(t, response.Body.String(), "10")
		AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Charlie")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestPOST(t *testing.T) {
	store := StubPlayerStore{
		score: map[string]int{},
	}
	server := mustMakePlayerServer(t, &store, dummyGame)

	t.Run("record win player on POST", func(t *testing.T) {
		player := "Pepper"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner, got %q want %q", store.winCalls[0], player)
		}
	})
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	// run integration test by using InMemoryPlayerStore instead of StubPlayerStore
	//store := InMemoryPlayerStore{store: map[string]int{}}

	// run test by using FileSystemPlayerStore
	database, cleanDB := createTempFile(t, "[]")
	defer cleanDB()
	//store := FileSystemPlayerStore{database}
	store, err := NewFileSystemPlayerStore(database)
	if err != nil {
		log.Fatalf("problem creating file system player store, %v", err)
	}

	server := mustMakePlayerServer(t, store, dummyGame)
	player := "charlie"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		AssertStatus(t, response.Code, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		AssertStatus(t, response.Code, http.StatusOK)
		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{"charlie", 3},
		}
		AssertLeague(t, got, want)
	})

}

func TestLeague(t *testing.T) {
	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := League{
			{"Cleo", 32},
			{"Charlie", 20},
			{"Toby", 14},
		}
		store := StubPlayerStore{nil, nil, wantedLeague}
		server := mustMakePlayerServer(t, &store, dummyGame)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)

		AssertStatus(t, response.Code, http.StatusOK)
		AssertLeague(t, got, wantedLeague)
		AssertContentType(t, response, "application/json")
	})
}

func getLeagueFromResponse(t *testing.T, body io.Reader) (league []Player) {
	t.Helper()
	league, err := NewLeague(body)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}

	return league
}

func newLeagueRequest() *http.Request {
	r, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return r
}

func newPostWinRequest(player string) *http.Request {
	r, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", player), nil)
	return r
}

func newGetScoreRequest(player string) *http.Request {
	r, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", player), nil)

	return r
}

func newGameRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return request
}

func mustMakePlayerServer(t *testing.T, store PlayerStore, game Game) *PlayerServer {
	server, err := NewPlayerServer(store, game)
	if err != nil {
		t.Fatalf("problem creating player server %v", err)
	}
	return server
}

func writeWSMessage(t *testing.T, conn *websocket.Conn, message string) {
	t.Helper()
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		t.Fatalf("could not send message over ws connetcion %v", err)
	}
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	t.Helper()
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}
	return ws
}

func within(t *testing.T, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out")
	case <-done:
	}
}
