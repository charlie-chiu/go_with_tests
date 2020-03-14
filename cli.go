package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const PlayerPrompt = "Please enter the number of players: "
const BadPlayerInputErrorMsg = "Bad value received for number of players, please try again with a number"
const BadWinnerInputErrorMsg = "invalid input, expected '{Name} wins'"

type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game Game
}

//just like constructor in php?
func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

func (cli *CLI) PlayPoker() {
	fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayerInput := cli.readLine()
	numberOfPlayers, err := strconv.Atoi(numberOfPlayerInput)
	if err != nil {
		fmt.Fprint(cli.out, BadPlayerInputErrorMsg)
		return
	}
	cli.game.Start(numberOfPlayers, cli.out)

	winnerInput := cli.readLine()
	winner, err := extractWinner(winnerInput)
	if err != nil {
		fmt.Fprint(cli.out, BadWinnerInputErrorMsg)
		return
	}
	cli.game.Finish(winner)
}

func extractWinner(userInput string) (string, error) {
	suffix := " wins"
	if !strings.Contains(userInput, suffix) {
		return "", fmt.Errorf("oops")
	}
	return strings.TrimSuffix(userInput, " wins"), nil
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
