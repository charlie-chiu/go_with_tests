package mocking

import (
	"bytes"
	"reflect"
	"testing"
)

// this spy struct implements both io.Writer and Sleeper
type CountdownOperationsSpy struct {
	Calls []string
}

const sleep = "sleep"
const write = "write"

func (s *CountdownOperationsSpy) Sleep() {
	s.Calls = append(s.Calls, sleep)
}
func (s *CountdownOperationsSpy) Write(p []byte) (n int, err error) {
	s.Calls = append(s.Calls, write)
	return
}

func TestCountdown(t *testing.T) {
	t.Run("print 3 to Go!", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		spySleeper := &CountdownOperationsSpy{}

		Countdown(buffer, spySleeper)

		got := buffer.String()
		want := "3\n2\n1\nGo!"
		// backtick creating a string, too
		//	want := `3
		//2
		//1
		//Go!`

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

		sleepCount := len(spySleeper.Calls)
		if sleepCount != 4 {
			t.Errorf("not enough calls to sleeper, want 4 got %d", sleepCount)
		}
	})

	t.Run("test sleep before every print", func(t *testing.T) {
		spySleeperPrinter := &CountdownOperationsSpy{}
		Countdown(spySleeperPrinter, spySleeperPrinter)

		want := []string{
			sleep, write,
			sleep, write,
			sleep, write,
			sleep, write,
		}

		got := spySleeperPrinter.Calls
		if !reflect.DeepEqual(want, got) {
			t.Errorf("wanted calls %v got %v", want, got)
		}
	})
}
