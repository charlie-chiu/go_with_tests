package clockface

import (
	"bytes"
	"encoding/xml"
	"math"
	"testing"
	"time"
)

func TestSVGWriterHourHand(t *testing.T) {
	cases := []struct {
		time time.Time
		line Line
	}{
		{simpleTime(6, 0, 0), Line{X1: 150, Y1: 150, X2: 150, Y2: 200}},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			b := bytes.Buffer{}
			SVGWriter(&b, c.time)
			svg := SVG{}
			xml.Unmarshal(b.Bytes(), &svg)

			if !containsLine(c.line, svg.Line) {
				t.Errorf("expected to find the hour hand line %+v, in the SVG lines %+v", c.line, svg.Line)
			}
		})
	}
}

func TestHourHandPoint(t *testing.T) {
	cases := []struct {
		time  time.Time
		point Point
	}{
		{simpleTime(6, 0, 0), Point{0, -1}},
		{simpleTime(12, 0, 0), Point{0, 1}},
		{simpleTime(21, 0, 0), Point{-1, 0}},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := hourHandPoint(c.time)
			if !roughlyEqualPoint(c.point, got) {
				t.Errorf("wanted %v point got %v", c.point, got)
			}
		})
	}
}

func TestHourHandInRadians(t *testing.T) {
	cases := []struct {
		t     time.Time
		angle float64
	}{
		{simpleTime(0, 0, 0), 0},
		{simpleTime(6, 0, 0), math.Pi},
		{simpleTime(21, 0, 0), math.Pi * 1.5},
		{simpleTime(0, 1, 30), math.Pi / ((6 * 60 * 60) / 90)},
	}

	for _, c := range cases {
		t.Run(testName(c.t), func(t *testing.T) {
			got := hourHandInRadians(c.t)
			if !roughlyEqualFloat64(got, c.angle) {
				t.Errorf("Wanted %v radians, got %v", c.angle, got)
			}
		})
	}
}

func TestSVGWriterMinuteHand(t *testing.T) {
	cases := []struct {
		time time.Time
		line Line
	}{
		{simpleTime(0, 0, 0), Line{X1: 150, Y1: 150, X2: 150, Y2: 70}},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			b := bytes.Buffer{}
			SVGWriter(&b, c.time)
			svg := SVG{}
			xml.Unmarshal(b.Bytes(), &svg)

			if !containsLine(c.line, svg.Line) {
				t.Errorf("expected to find the minute hand line %+v, in the SVG lines %+v", c.line, svg.Line)
			}
		})
	}
}

func TestMinuteHandPoint(t *testing.T) {
	cases := []struct {
		time  time.Time
		point Point
	}{
		{simpleTime(0, 0, 0), Point{0, 1}},
		{simpleTime(0, 30, 0), Point{0, -1}},
		{simpleTime(0, 45, 0), Point{-1, 0}},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := minuteHandPoint(c.time)
			if !roughlyEqualPoint(c.point, got) {
				t.Errorf("wanted %v point got %v", c.point, got)
			}
		})
	}
}

func TestMinuteHandInRadians(t *testing.T) {
	cases := []struct {
		t     time.Time
		angle float64
	}{
		{simpleTime(0, 0, 0), 0},
		{simpleTime(0, 30, 0), math.Pi},
		{simpleTime(0, 45, 0), (math.Pi / 2) * 3},
		{simpleTime(0, 41, 0), (math.Pi / 30) * 41},
		// minute hand should move a little bit every second!
		{simpleTime(0, 0, 7), 7 * (math.Pi / (30 * 60))},
	}

	for _, c := range cases {
		t.Run(testName(c.t), func(t *testing.T) {
			got := minutesInRadians(c.t)
			if !roughlyEqualFloat64(got, c.angle) {
				t.Errorf("Wanted %v radians, got %v", c.angle, got)
			}
		})
	}
}

func TestSVGWriterSecondHand(t *testing.T) {
	cases := []struct {
		time time.Time
		line Line
	}{
		{
			simpleTime(0, 0, 0),
			Line{X1: clockCentreX, Y1: clockCentreY, X2: clockCentreX, Y2: clockCentreY - secondHandLength},
		},
		{
			simpleTime(0, 0, 15),
			Line{X1: clockCentreX, Y1: clockCentreY, X2: clockCentreX + secondHandLength, Y2: clockCentreY},
		},
		{
			simpleTime(0, 0, 30),
			Line{X1: clockCentreX, Y1: clockCentreY, X2: clockCentreX, Y2: clockCentreY + secondHandLength},
		},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			b := bytes.Buffer{}
			SVGWriter(&b, c.time)
			svg := SVG{}
			xml.Unmarshal(b.Bytes(), &svg)

			if !containsLine(c.line, svg.Line) {
				t.Errorf("expected to find the second hand line %+v, in the SVG lines %+v", c.line, svg.Line)
			}
		})
	}
}

func TestSecondHandPoint(t *testing.T) {
	cases := []struct {
		time  time.Time
		point Point
	}{
		{simpleTime(0, 0, 30), Point{0, -1}},
		{simpleTime(0, 0, 45), Point{-1, 0}},
	}

	for _, c := range cases {
		t.Run(testName(c.time), func(t *testing.T) {
			got := secondHandPoint(c.time)
			if !roughlyEqualPoint(c.point, got) {
				t.Errorf("wanted %v point got %v", c.point, got)
			}
		})
	}
}

func TestSecondHandInRadians(t *testing.T) {
	cases := []struct {
		t     time.Time
		angle float64
	}{
		{simpleTime(0, 0, 30), math.Pi},
		{simpleTime(0, 0, 0), 0},
		{simpleTime(0, 0, 45), (math.Pi / 2) * 3},
		{simpleTime(0, 0, 7), (math.Pi / secondsInHalfClock) * 7},
	}

	for _, c := range cases {
		t.Run(testName(c.t), func(t *testing.T) {
			got := secondsInRadians(c.t)
			if !roughlyEqualFloat64(got, c.angle) {
				t.Errorf("Wanted %v radians, got %v", c.angle, got)
			}
		})
	}
}

// helper

func containsLine(needle Line, haystack []Line) bool {
	for _, line := range haystack {
		if line == needle {
			return true
		}
	}

	return false
}

func simpleTime(hours, minutes, seconds int) time.Time {
	return time.Date(312, time.January, 1, hours, minutes, seconds, 0, time.UTC)
}

func testName(time time.Time) string {
	return time.Format("15:04:05")
}

func roughlyEqualFloat64(a, b float64) bool {
	const equalityThreshold = 1e-7
	return math.Abs(a-b) < equalityThreshold
}

func roughlyEqualPoint(a, b Point) bool {
	return roughlyEqualFloat64(a.X, b.X) && roughlyEqualFloat64(a.Y, b.Y)
}
