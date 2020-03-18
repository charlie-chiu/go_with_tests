package poker_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	poker "github.com/charlie-chiu/go_with_test"
	"github.com/gorilla/websocket"
)

var (
	dummyGame = &GameSpy{}
)

func TestWebSocket(t *testing.T) {
	t.Run("start a game with 3 players and declare Charlie the winner", func(t *testing.T) {
		wantedBlindAlert := "Blind is 100"

		game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}
		playerServer := mustMakePlayerServer(t, &poker.StubPlayerStore{}, game)
		server := httptest.NewServer(playerServer)
		ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

		defer server.Close()
		defer ws.Close()

		numberOfPlayers := 3
		winner := "Charlie"
		writeWSMessage(t, ws, strconv.Itoa(numberOfPlayers))
		writeWSMessage(t, ws, winner)

		const tenMS = 10 * time.Millisecond
		assertGameStartedWith(t, game, numberOfPlayers)
		assertFinishCalledWith(t, game, winner)

		within(t, tenMS, func() { assertWebsocketGotMsg(t, ws, wantedBlindAlert) })
	})
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
		server := mustMakePlayerServer(t, &poker.StubPlayerStore{}, dummyGame)

		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
}

func TestGETPlayers(t *testing.T) {
	store := poker.StubPlayerStore{Score: map[string]int{
		"Pepper": 20,
		"Floyd":  10,
	}}
	server := mustMakePlayerServer(t, &store, dummyGame)

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "20")
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "10")
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Charlie")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestPOST(t *testing.T) {
	t.Run("record win player on POST", func(t *testing.T) {
		store := poker.StubPlayerStore{
			Score: map[string]int{},
		}
		server := mustMakePlayerServer(t, &store, dummyGame)
		player := "Pepper"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.WinCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.WinCalls), 1)
		}

		if store.WinCalls[0] != player {
			t.Errorf("did not store correct winner, got %q want %q", store.WinCalls[0], player)
		}
	})
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	// run integration test by using InMemoryPlayerStore instead of StubPlayerStore
	//store := InMemoryPlayerStore{store: map[string]int{}}

	// run test by using FileSystemPlayerStore
	dbFile, cleanDB := createTempFile(t, "[]")
	defer cleanDB()
	store, err := poker.NewFileSystemPlayerStore(dbFile)
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
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertStatus(t, response.Code, http.StatusOK)
		got := getLeagueFromResponse(t, response.Body)
		want := poker.League{
			{"charlie", 3},
		}
		assertLeague(t, got, want)
	})

}

func TestLeague(t *testing.T) {
	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := poker.League{
			{"Cleo", 32},
			{"Charlie", 20},
			{"Toby", 14},
		}
		store := poker.StubPlayerStore{League: wantedLeague}
		server := mustMakePlayerServer(t, &store, dummyGame)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)
		assertContentType(t, response, "application/json")
	})
}

func getLeagueFromResponse(t *testing.T, body io.Reader) (league poker.League) {
	t.Helper()
	league, err := poker.NewLeague(body)
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

func mustMakePlayerServer(t *testing.T, store poker.PlayerStore, game poker.Game) *poker.PlayerServer {
	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		t.Fatalf("problem creating player server %v", err)
	}
	return server
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	contentType := response.Result().Header.Get("content-type")
	if contentType != want {
		t.Errorf("response did not have content-type of %s, got %v", want, contentType)
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got status %d, want %d", got, want)
	}
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
