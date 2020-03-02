package structs_methods_interfaces

import "testing"

type Shape interface {
	Area() float64
	Perimeter() float64
}

func TestArea(t *testing.T) {
	checkArea := func(t *testing.T, s Shape, want float64) {
		t.Helper()
		got := s.Area()

		if got != want {
			t.Errorf("got %g, want %g", got, want)
		}
	}

	t.Run("rectangles", func(t *testing.T) {
		rectangle := Rectangle{
			Width:  12.0,
			Height: 6.0,
		}
		checkArea(t, rectangle, 72.0)
	})

	t.Run("circles", func(t *testing.T) {
		c := Circle{Radius: 10}

		checkArea(t, c, 314.1592653589793)
	})

}

func TestAreaWithTable(t *testing.T) {
	//TableDrivenTests
	areaTests := []struct {
		name  string
		shape Shape
		want  float64
	}{
		{"Rectangle", Rectangle{9, 12}, 108.0},
		{"Circle", Circle{10}, 314.1592653589793},
		{"Triangle", Triangle{10, 10}, 50.0},
	}

	for _, tt := range areaTests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.shape.Area()
			if got != tt.want {
				// can output with chinese
				t.Errorf("shape:%#v 面積錯誤： 得到：%g, 預期：%g", tt.shape, got, tt.want)
			}
		})
	}
}

func TestPerimeter(t *testing.T) {
	//TableDrivenTests
	areaTests := []struct {
		name  string
		shape Shape
		want  float64
	}{
		{"Rectangle", Rectangle{9, 12}, 42.0},
		{"Circle", Circle{10}, 62.83185307179586},
		//don't know how to calculate
		//{"Triangle", Triangle{10, 10}, ???},
	}

	for _, tt := range areaTests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.shape.Perimeter()
			if got != tt.want {
				t.Errorf("got %g, want %g", got, tt.want)
			}
		})
	}
}
