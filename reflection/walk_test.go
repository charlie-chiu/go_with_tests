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
}
