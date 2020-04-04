package clockface

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"time"
)

type SVG struct {
	XMLName xml.Name `xml:"svg"`
	Xmlns   string   `xml:"xmlns,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	ViewBox string   `xml:"viewBox,attr"`
	Version string   `xml:"version,attr"`
	Circle  Circle   `xml:"circle"`
	Line    []Line   `xml:"line"`
}

type Circle struct {
	Cx float64 `xml:"cx,attr"`
	Cy float64 `xml:"cy,attr"`
	R  float64 `xml:"r,attr"`
}

type Line struct {
	X1 float64 `xml:"x1,attr"`
	Y1 float64 `xml:"y1,attr"`
	X2 float64 `xml:"x2,attr"`
	Y2 float64 `xml:"y2,attr"`
}

//writes an SVG representation of an analogue clock, showing the time t, to the writer w
func SVGWriter(w io.Writer, t time.Time) {
	io.WriteString(w, svgStart)
	io.WriteString(w, bezel)
	clockHandWriter(w, SecondHand(t))
	clockHandWriter(w, MinuteHand(t))
	clockHandWriter(w, HourHand(t))
	io.WriteString(w, svgEnd)
}

func clockHandWriter(w io.Writer, p Point) {
	fmt.Fprintf(w, `<line x1="150" y1="150" x2="%.3f" y2="%.3f" style="fill:none;stroke:#f00;stroke-width:3px;"/>`, p.X, p.Y)
}

const svgStart = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg xmlns="http://www.w3.org/2000/svg"
     width="100%"
     height="100%"
     viewBox="0 0 300 300"
     version="2.0">`

const bezel = `<circle cx="150" cy="150" r="100" style="fill:#fff;stroke:#000;stroke-width:5px;"/>`

const svgEnd = `</svg>`

const (
	secondHandLength = 90
	minuteHandLength = 80
	hourHandLength   = 50
	clockCentreX     = 150
	clockCentreY     = 150
)

const secondsInHalfClock = 30
const secondsInClock = 2 * secondsInHalfClock
const minutesInHalfClock = 30

// a point represents a two dimensional Cartesian coordinate
type Point struct {
	X, Y float64
}

func HourHand(t time.Time) (p Point) {
	return makeHand(hourHandPoint(t), hourHandLength)
}

func hourHandPoint(t time.Time) Point {
	return angleToPoint(hourHandInRadians(t))
}

func hourHandInRadians(t time.Time) float64 {
	return math.Pi/(6/float64(t.Hour()%12)) + (minutesInRadians(t) / 12)
}

func MinuteHand(t time.Time) (p Point) {
	return makeHand(minuteHandPoint(t), minuteHandLength)
}

func minutesInRadians(t time.Time) float64 {
	return math.Pi/(minutesInHalfClock/float64(t.Minute())) + secondsInRadians(t)/secondsInClock
}

func minuteHandPoint(t time.Time) Point {
	return angleToPoint(minutesInRadians(t))
}

func SecondHand(t time.Time) (p Point) {
	return makeHand(secondHandPoint(t), secondHandLength)
}

func secondsInRadians(t time.Time) float64 {
	return math.Pi / (secondsInHalfClock / float64(t.Second()))
}

func secondHandPoint(t time.Time) Point {
	return angleToPoint(secondsInRadians(t))
}

func angleToPoint(angle float64) Point {
	x := math.Sin(angle)
	y := math.Cos(angle)

	return Point{x, y}
}

func makeHand(p Point, length float64) Point {
	// scale
	p = Point{X: p.X * length, Y: p.Y * length}
	// flip, origin of SVG in the top left hand corner
	p = Point{X: p.X, Y: -p.Y}
	// translate
	p = Point{X: p.X + clockCentreX, Y: p.Y + clockCentreY}

	return p
}
