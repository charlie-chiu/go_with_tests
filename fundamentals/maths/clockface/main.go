package main

import (
	clockface "github.com/charlie-chiu/go_with_test/fundamentals/maths"
	"os"
	"time"
)

func main() {
	t := time.Now()
	clockface.SVGWriter(os.Stdout, t)
}
