package clockface

import (
	"math"
	"time"
)

const secondHandLength = 90
const clockCentreX = 150
const clockCentreY = 150

// a point represents a two dimensional Cartesian coordinate
type Point struct {
	X, Y float64
}

func SecondHand(t time.Time) (p Point) {
	p = secondHandPoint(t)

	// scale
	p = Point{X: p.X * secondHandLength, Y: p.Y * secondHandLength}
	// flip, origin of SVG in the top left hand corner
	p = Point{X: p.X, Y: -p.Y}
	// translate
	p = Point{X: p.X + clockCentreX, Y: p.Y + clockCentreY}

	return
}

func secondsInRadians(t time.Time) float64 {
	return math.Pi / (30 / float64(t.Second()))
}

func secondHandPoint(t time.Time) Point {
	angle := secondsInRadians(t)
	x := math.Sin(angle)
	y := math.Cos(angle)

	return Point{x, y}
}
