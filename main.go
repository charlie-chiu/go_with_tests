package main

import (
	"os"
)
import "github.com/charlie-chiu/go_with_test/mocking"

func main() {
	sleeper := &mocking.DefaultSleeper{}
	mocking.Countdown(os.Stdout, sleeper)
}
