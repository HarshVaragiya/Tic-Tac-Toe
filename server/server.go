package server

import (
	"context"
	"encoding/json"
	"fmt"
	"ticTacToe/structures"
	"ticTacToe/structures/tictactoe"

	log "github.com/sirupsen/logrus"
)

type GameServer struct {
	Games map[string]*structures.Game
}

func (this GameServer) StartNewGame(ctx context.Context, request *tictactoe.StartNewGameRequest) (*tictactoe.StartNewGameResponse, error) {
	game := structures.StartNewGame(request.Player, int(request.GridSize))
	this.Games[game.GameId] = game
	log.Infof("player %v started new game with Id %v", request.Player, game.GameId)
	return &tictactoe.StartNewGameResponse{GameId: game.GameId}, nil
}

func (this GameServer) JoinGame(ctx context.Context, request *tictactoe.JoinGameRequest) (*tictactoe.JoinGameResponse, error) {
	game, exists := this.Games[request.GameId]
	if !exists {
		return nil, fmt.Errorf("game '%v' does not exist", request.GameId)
	}
	response := &tictactoe.JoinGameResponse{}
	if err := game.JoinGame(request.Player); err != nil {
		response.Joined = false
		response.Message = err.Error()
	} else {
		response.Joined = true
		response.Message = "joined game"
	}
	log.Printf("player %v joined game %v", request.Player, request.GameId)
	return response, nil
}

func (this GameServer) GetUpdatedGameDetails(ctx context.Context, request *tictactoe.GameRequest) (*tictactoe.GameDetails, error) {
	game, exists := this.Games[request.GameId]
	if !exists {
		return nil, fmt.Errorf("game '%v' does not exist", request.GameId)
	}
	data, err := json.Marshal(game)
	if err != nil {
		log.Error("error converting game structure to JSON")
		return nil, fmt.Errorf("internal game server error")
	}
	response := &tictactoe.GameDetails{GameId: request.GameId, GameDetailsJson: string(data)}
	return response, nil
}

func (this GameServer) PutMarker(ctx context.Context, request *tictactoe.PutMarkerRequest) (*tictactoe.PutMarkerResponse, error) {
	game, exists := this.Games[request.GameId]
	if !exists {
		return nil, fmt.Errorf("game '%v' does not exist", request.GameId)
	}
	response := &tictactoe.PutMarkerResponse{}
	if err := game.PutMarker(request.Player, int(request.Position)); err != nil {
		response.Approved = false
		response.Message = err.Error()
	} else {
		response.Approved = true
	}
	response.Finished = game.Finished
	return response, nil
}

func (this GameServer) GetWinnerDetails(ctx context.Context, request *tictactoe.GameRequest) (*tictactoe.WinnerDetails, error) {
	game, exists := this.Games[request.GameId]
	if !exists {
		return nil, fmt.Errorf("game '%v' does not exist", request.GameId)
	}
	response := &tictactoe.WinnerDetails{}
	if game.WinnerPlayerId != -1 {
		response.Player = game.PlayerMap[game.WinnerPlayerId]
		response.Message = fmt.Sprintf("Player %v '%v' Wins!", game.WinnerPlayerId, response.Player)
	} else {
		response.Player = "None"
		response.Message = "DRAW!"
	}
	log.Printf("game %v result : %v", request.GameId, response.Message)
	return response, nil
}
