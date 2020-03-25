package pkg_test

import (
	"context"
	pkg "github.com/charlie-chiu/go_with_test/fundamentals/context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type SpyStore struct {
	response  string
	cancelled bool
}

func (s *SpyStore) Cancel() {
	s.cancelled = true
}

func (s *SpyStore) Fetch() string {
	return s.response
}

func (s *SpyStore) assertWasCancelled(t *testing.T) {
	t.Helper()
	if !s.cancelled {
		t.Errorf("store was not told to cancel")
	}
}

func (s *SpyStore) assertWasNotCancelled(t *testing.T) {
	t.Helper()
	if s.cancelled {
		t.Errorf("it should not have cancelled the store")
	}
}

func TestHandler(t *testing.T) {
	var data = "hello, world"

	t.Run("happy path", func(t *testing.T) {
		store := &SpyStore{response: data}
		svr := pkg.Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		got := response.Body.String()
		if got != data {

			t.Errorf("got %q, want %q", got, data)
		}

		store.assertWasNotCancelled(t)
	})

	t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
		store := &SpyStore{
			response:  data,
			cancelled: false,
		}
		svr := pkg.Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)

		cancellingCtx, cancelFunc := context.WithCancel(request.Context())
		time.AfterFunc(5*time.Millisecond, cancelFunc)
		time.Sleep(7 * time.Millisecond)
		request = request.WithContext(cancellingCtx)

		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		store.assertWasCancelled(t)
	})
}
