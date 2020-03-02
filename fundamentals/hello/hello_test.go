package main

import (
	"testing"
)

func TestHello(t *testing.T) {
	assertCorrectMessage := func(t *testing.T, got, want string) {
		// t.Helper()
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}

	t.Run("empty string defaults to 'World'", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello, World"

		assertCorrectMessage(t, got, want)
	})

	// t.Run("a test to fail", func(t *testing.T) {
	// 	got := Hello("")
	// 	want := "YOU SHALL NOT PASS"

	// 	assertCorrectMessage(t, got, want)
	// })

	t.Run("saying hello to people", func(t *testing.T) {
		got := Hello("Charlie", "")
		want := "Hello, Charlie"

		assertCorrectMessage(t, got, want)
	})

	t.Run("in Spanish", func(t *testing.T) {
		got := Hello("Elodie", "Spanish")
		Want := "Hola, Elodie"
		assertCorrectMessage(t, got, Want)
	})

	t.Run("in French", func(t *testing.T) {
		got := Hello("Charlie", "French")
		Want := "Bonjour, Charlie"
		assertCorrectMessage(t, got, Want)
	})
}
