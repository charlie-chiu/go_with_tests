package maps

import "testing"

func TestSearch(t *testing.T) {
	dictionary := Dictionary{
		"test": "this is just a test",
	}

	t.Run("known word", func(t *testing.T) {
		got, _ := dictionary.Search("test")
		want := "this is just a test"

		assertEqual(t, got, want)
	})

	t.Run("unknown word", func(t *testing.T) {
		_, got := dictionary.Search("unknown")

		assertError(t, got, ErrNotFound)
	})
}

func TestAdd(t *testing.T) {
	d := Dictionary{}

	t.Run("new word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		_ = d.Add(word, definition)

		got, err := d.Search(word)
		if err != nil {
			t.Fatal("should find added word:", err)
		}

		assertEqual(t, got, definition)
	})

	t.Run("existing word", func(t *testing.T) {
		word := "test"
		definition := "something else"
		err := d.Add(word, definition)

		assertError(t, err, ErrWordExists)

		got, _ := d.Search(word)
		assertEqual(t, got, "this is just a test")
	})
}

func assertEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func assertError(t *testing.T, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}
