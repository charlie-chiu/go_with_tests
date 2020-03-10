package poker

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

const PlayerPrompt = "Please enter the number of players: "

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

type CLI struct {
	store   PlayerStore
	in      *bufio.Scanner
	out     io.Writer
	alerter BlindAlerter
}

//just like constructor in php?
func NewCLI(store PlayerStore, in io.Reader, out io.Writer, alerter BlindAlerter) *CLI {
	return &CLI{
		store:   store,
		in:      bufio.NewScanner(in),
		out:     out,
		alerter: alerter,
	}
}

func (cli *CLI) PlayPoker() {
	cli.scheduleBlindAlerts()
	userInput := cli.readLine()
	fmt.Fprint(cli.out, PlayerPrompt)
	cli.store.RecordWin(extractWinner(userInput))
}

func (cli *CLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + 10*time.Minute
		//blindTime = blindTime + 10*time.Second
	}
}

func extractWinner(userInput string) string {
	return strings.TrimSuffix(userInput, " wins")
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
