package client

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"ticTacToe/structures/tictactoe"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	pollingTime = time.Millisecond * 200
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func PlayGame(config *ServerConfiguration, gameClient *tictactoe.TicTacToeGameClient) {

	cli := NewGameClient(*gameClient, config.PlayerName)
	reader := bufio.NewReader(os.Stdin)

	ctx := context.Background()

	fmt.Printf("[0] Start New Game \n[1] Join Existing Game\nChoice: ")
	choice := GetIntFromStdin(reader)

	switch choice {
	case 0:
		cli.StartNewGame(ctx)
		fmt.Printf("Waiting for player 2 to join ...\n")
	case 1:
		fmt.Printf("Enter the game Id: ")
		cli.JoinGame(ctx, GetStringFromStdin(reader))
	}

	for {
		cli.GetGameDetails(ctx)
		if cli.Game.ReadyToPlay {
			log.Printf("game is now ready to play!")
			break
		}
		time.Sleep(pollingTime)
	}

	previousGameText := ""
	for {
		cli.GetGameDetails(ctx)
		gameText := cli.Game.DisplayGame()

		if gameText != previousGameText {
			// some update happened!
			fmt.Printf("%v\n", strings.Repeat("=", 60))
			fmt.Printf("Game: \n\n%v", gameText)
			previousGameText = gameText
		}

		if cli.GameOver {
			fmt.Printf("\n%vGAME OVER!%v\n", colorRed, colorReset)
			_, message, _ := cli.GetWinnerDetails(ctx)
			fmt.Printf("%vServer: %v %v", colorGreen, message, colorReset)
			break
		}

		if strings.Contains(gameText, config.PlayerName) {
			fmt.Printf("\nYour Move: ")
			pos := GetIntFromStdin(reader)
			_, _, _ = cli.PutMarker(ctx, pos)
		}
		time.Sleep(pollingTime)
	}
}

func GetStringFromStdin(reader *bufio.Reader) string {
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Error("error reading from stdin")
			continue
		}
		return strings.TrimSuffix(str, "\n")
	}
}

func GetIntFromStdin(reader *bufio.Reader) int {
	for {
		number, err := strconv.Atoi(GetStringFromStdin(reader))
		if err != nil {
			log.Error("error converting input to int. please retry")
			continue
		}
		return number
	}
}
