package client

import (
	"context"
	"encoding/json"
	"fmt"
	"ticTacToe/structures"
	"ticTacToe/structures/tictactoe"

	log "github.com/sirupsen/logrus"
)

func NewGameClient(rpcClient tictactoe.TicTacToeGameClient, playerName string) *GameClient {
	return &GameClient{
		RpcClient:  rpcClient,
		PlayerName: playerName,
		Game:       &structures.Game{},
	}
}

type GameClient struct {
	RpcClient  tictactoe.TicTacToeGameClient
	GameId     string
	Game       *structures.Game
	PlayerName string
	GameOver   bool
}

func (this *GameClient) StartNewGame(ctx context.Context) error {
	request := &tictactoe.StartNewGameRequest{
		Player:   this.PlayerName,
		GridSize: 3,
	}
	resp, err := this.RpcClient.StartNewGame(ctx, request)
	if err != nil {
		log.Error("error starting new game")
		log.Error(err)
		return err
	}
	this.GameId = resp.GameId
	log.Printf("started new game with Id: %v", this.GameId)
	return nil
}

func (this *GameClient) JoinGame(ctx context.Context, gameId string) error {
	request := &tictactoe.JoinGameRequest{
		GameId: gameId,
		Player: this.PlayerName,
	}
	resp, err := this.RpcClient.JoinGame(ctx, request)
	if err != nil {
		log.Error("error joining game")
		log.Error(err)
		return err
	}
	if !resp.Joined {
		log.Error(resp.Message)
		return fmt.Errorf(resp.Message)
	}
	log.Printf("Server: %v", resp.Message)
	this.GameId = gameId
	return nil
}

func (this *GameClient) GetGameDetails(ctx context.Context) error {
	request := &tictactoe.GameRequest{GameId: this.GameId}
	gameDetails, err := this.RpcClient.GetUpdatedGameDetails(ctx, request)
	if err != nil {
		log.Error("error getting game details")
		log.Error(err)
		return err
	}
	err = json.Unmarshal([]byte(gameDetails.GameDetailsJson), this.Game)
	if err != nil {
		log.Error("error parsing server response")
		log.Error(err)
		return err
	}
	this.GameOver = this.Game.Finished
	return nil
}

func (this *GameClient) PutMarker(ctx context.Context, position int) (bool, bool, error) {
	request := &tictactoe.PutMarkerRequest{
		Position: int32(position),
		GameId:   this.GameId,
		Player:   this.PlayerName,
	}
	resp, err := this.RpcClient.PutMarker(ctx, request)
	if err != nil {
		log.Error("error putting marker")
		log.Error(err)
		return false, false, err
	}
	this.GameOver = resp.Finished
	if !resp.Approved {
		log.Error(resp.Message)
	}
	return resp.Approved, resp.Finished, nil
}

func (this *GameClient) GetWinnerDetails(ctx context.Context) (string, string, error) {
	request := &tictactoe.GameRequest{GameId: this.GameId}
	resp, err := this.RpcClient.GetWinnerDetails(ctx, request)
	if err != nil {
		log.Error("error getting winner details")
		log.Error(err)
		return "", "", err
	}
	return resp.Player, resp.Message, nil
}
