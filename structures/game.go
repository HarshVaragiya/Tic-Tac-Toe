package structures

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorGray   = "\u001b[38;5;245m"
)

type Game struct {
	GameId         string
	PlayerMap      map[int]string
	PositionMap    map[int]int
	GridSize       int
	Finished       bool
	WinnerPlayerId int
	ReadyToPlay    bool
	EventCount     int
}

func StartNewGame(player1 string, GridSize int) *Game {
	game := &Game{}
	game.GameId = uuid.NewString()
	game.PlayerMap = make(map[int]string, 2)
	game.PlayerMap[0] = player1
	game.PositionMap = make(map[int]int, GridSize*GridSize)
	game.GridSize = GridSize
	for i := 0; i < GridSize*GridSize; i++ {
		game.PositionMap[i] = 0
	}
	game.WinnerPlayerId = -1
	return game
}

func (this *Game) IsAvailableForPlaying() error {
	if this.Finished {
		return fmt.Errorf("game is finished, it is not playable")
	}
	if !this.ReadyToPlay {
		return fmt.Errorf("game is not in a ready to play state")
	}
	if this.EventCount == this.GridSize*this.GridSize {
		this.ReadyToPlay = false
		this.Finished = true
		this.CheckIfGameIsFinished()
		return fmt.Errorf("game has ended!")
	}
	return nil
}

func (this *Game) JoinGame(player2 string) error {
	if !this.Finished && !this.ReadyToPlay {
		if _, exists := this.PlayerMap[1]; exists {
			return fmt.Errorf("player 2 has already joined. cannot join a full game")
		}
		this.PlayerMap[1] = player2
		this.ReadyToPlay = true
		return nil
	} else {
		return fmt.Errorf("game not in a playable state ")
	}
}

func (this *Game) GetPlayerMarker(player string) (int, error) {
	if this.PlayerMap[0] == player {
		return -1, nil
	} else if this.PlayerMap[1] == player {
		return 1, nil
	}
	return 0, fmt.Errorf("player not among the game members")
}

func (this *Game) PutMarker(player string, position int) error {
	if err := this.IsAvailableForPlaying(); err != nil {
		return err
	}
	playerMarker, err := this.GetPlayerMarker(player)
	if err != nil {
		return err
	}
	nextPlayerMove := playerMarker
	if playerMarker == -1 {
		nextPlayerMove = 0
	}

	if this.EventCount%2 != nextPlayerMove {
		return fmt.Errorf("this player is not the one supposed to make the next move")
	}

	if value, exists := this.PositionMap[position]; exists && value != 0 {
		return fmt.Errorf("position not empty. cannot place marker there")
	}
	this.PositionMap[position] = playerMarker
	this.CheckIfGameIsFinished()
	this.EventCount += 1
	return nil
}

func (this *Game) CheckIfGameIsFinished() bool {
	// check horizontally
	for i := 0; i < this.GridSize; i++ {
		rowSum := 0
		for j := 0; j < this.GridSize; j++ {
			rowSum += this.PositionMap[(i*this.GridSize)+j]
		}
		if rowSum == this.GridSize {
			this.markWinner(1) // player 2 won!
			return true
		} else if rowSum == -1*this.GridSize {
			this.markWinner(0) // player 1 won!
			return true
		}
	}

	// check vertically
	for i := 0; i < this.GridSize; i++ {
		columnSum := 0
		for j := 0; j < this.GridSize; j++ {
			columnSum += this.PositionMap[i+(j*this.GridSize)]
		}
		if columnSum == this.GridSize {
			this.markWinner(1) // player 2 won!
			return true
		} else if columnSum == -1*this.GridSize {
			this.markWinner(0) // player 1 won!
			return true
		}
	}

	// check diagonally
	diagonalSum1, diagonalSum2 := 0, 0
	for i := 0; i < this.GridSize; i++ {
		diagonalSum1 += this.PositionMap[(i*this.GridSize)+i]
		diagonalSum2 += this.PositionMap[((this.GridSize-i-1)*this.GridSize)+i]
	}
	if diagonalSum1 == this.GridSize || diagonalSum2 == this.GridSize {
		this.markWinner(1) // player 2 won!
		return true
	} else if diagonalSum1 == -1*this.GridSize || diagonalSum2 == -1*this.GridSize {
		this.markWinner(0) // player 1 won!
		return true
	}
	return false
}

func (this *Game) markWinner(WinnerPlayerId int) {
	this.Finished = true
	this.WinnerPlayerId = WinnerPlayerId
	this.ReadyToPlay = false
}

func (this *Game) DisplayGame() string {
	gameText := ""
	for i := 0; i < this.GridSize; i++ {
		stringRepr := ""
		for j := 0; j < this.GridSize; j++ {
			index := (i * this.GridSize) + j
			value := this.PositionMap[index]
			switch value {
			case -1:
				stringRepr += colorYellow + " X " + colorReset
			case 0:
				stringRepr += "   "
			case 1:
				stringRepr += colorCyan + " O " + colorReset
			}
			if j != this.GridSize-1 {
				stringRepr += "|"
			}
		}
		stringRepr += strings.Repeat(" ", 30)
		for j := 0; j < this.GridSize; j++ {
			index := (i * this.GridSize) + j
			value := this.PositionMap[index]
			if value == 0 {
				stringRepr += fmt.Sprintf("%v %v %v", colorGray, index, colorReset)
			} else {
				stringRepr += fmt.Sprintf("   ")
			}
			if j != this.GridSize-1 {
				stringRepr += "|"
			}
		}

		gameText += fmt.Sprintf("%v\n", stringRepr)
		if i != this.GridSize-1 {

			horizontalSeperator := ""
			for k := 0; k < this.GridSize; k++ {
				horizontalSeperator += "---"
				if k != this.GridSize-1 {
					horizontalSeperator += "+"
				}
			}
			horizontalSeperator += strings.Repeat(" ", 30) + horizontalSeperator
			gameText += fmt.Sprintf("%v\n", horizontalSeperator)
		}
	}
	gameText += fmt.Sprintf("\nNext move by player %v: %v \n", this.EventCount%2, this.PlayerMap[this.EventCount%2])
	return gameText
}
