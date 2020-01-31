package integers

import "testing"

import "fmt"

func ExampleAdd() {
	sum := Add(4, 6)
	fmt.Println(sum)
	// output : 10
}

func TestAdder(t *testing.T) {
	sum := Add(2, 2)
	expected := 4

	if sum != expected {
		t.Errorf("expected : %d, got %d", expected, sum)
	}
}
