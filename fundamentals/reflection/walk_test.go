package reflection

import (
	"reflect"
	"testing"
)

type Person struct {
	Name    string
	Profile Profile
}

type Profile struct {
	Age  int
	City string
}

func TestWalk(t *testing.T) {

	cases := []struct {
		Name          string
		Input         interface{}
		ExpectedCalls []string
	}{
		{
			"struct with one string field",
			struct{ Name string }{"Charlie"},
			[]string{"Charlie"},
		},
		{
			"struct with two string fields",
			struct {
				Name string
				City string
			}{"Charlie", "Taipei"},
			[]string{"Charlie", "Taipei"},
		},
		{
			"struct with non string fields",
			struct {
				Name string
				Age  int
			}{"Charlie", 7788},
			[]string{"Charlie"},
		},
		{
			"Nested fields",
			Person{
				Name: "Charlie",
				Profile: Profile{
					Age:  7888,
					City: "Taipei",
				},
			},
			[]string{"Charlie", "Taipei"},
		},
		{
			Name: "Pointers to things",
			Input: &Person{
				"Charlie",
				Profile{
					Age:  7788,
					City: "Taipei",
				},
			},
			ExpectedCalls: []string{"Charlie", "Taipei"},
		},
		{
			Name: "Slices",
			Input: []Profile{
				{7788, "foo"},
				{7788, "bar"},
			},
			ExpectedCalls: []string{"foo", "bar"},
		},
		{
			Name: "Arrays",
			Input: [2]Profile{
				{7788, "foo"},
				{7788, "bar"},
			},
			ExpectedCalls: []string{"foo", "bar"},
		},
		{
			// this test sometimes fail because map do not guarantee order.
			Name: "Maps",
			Input: map[int]string{
				3:  "Three",
				10: "ten",
			},
			ExpectedCalls: []string{"Three", "ten"},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			var got []string
			walk(test.Input, func(input string) {
				got = append(got, input)
			})

			if !reflect.DeepEqual(got, test.ExpectedCalls) {
				t.Errorf("got %v, want %v", got, test.ExpectedCalls)
			}
		})
	}

	t.Run("with maps", func(t *testing.T) {
		aMap := map[int]string{
			3:  "Three",
			10: "ten",
		}

		var got []string
		walk(aMap, func(input string) {
			got = append(got, input)
		})

		assertContains(t, got, "Three")
		assertContains(t, got, "ten")
	})
}

func assertContains(t *testing.T, haystack []string, needle string) {
	contains := false
	for _, x := range haystack {
		if x == needle {
			contains = true
		}
	}

	if !contains {
		t.Errorf("expected %+v to containe %q but it didn't", haystack, needle)
	}
}
