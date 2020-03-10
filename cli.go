package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const PlayerPrompt = "Please enter the number of players: "

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game *Game
}

//just like constructor in php?
func NewCLI(in io.Reader, out io.Writer, game *Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayerInput := cli.readLine()
	numberOfPlayers, _ := strconv.Atoi(numberOfPlayerInput)
	cli.game.Start(numberOfPlayers)

	winnerInput := cli.readLine()
	winner := extractWinner(winnerInput)
	cli.game.Finish(winner)
}

func extractWinner(userInput string) string {
	return strings.TrimSuffix(userInput, " wins")
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
