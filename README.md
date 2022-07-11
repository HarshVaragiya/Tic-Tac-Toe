# Tic Tac Toe

a simple project to try to build a 2 player "tic-tac-toe" game :) 

## Why ?
- just for fun... 
- to be able to play a game of tic-tac-toe with my friends remotely

## Architecture 
- game server is responsible for coordinating the game and is the source of truth
- clients connect to game server to play game remotely 

![assets/diagram.png](assets/diagram.png)


## Server 
```
make cli 
./bin/tictactoe server
```
![assets/server.png](assets/server.png)

## Clients

- Player 1
```
bin/tictactoe 
```
![assets/player-1.png](assets/player-1.png)

- Player 2
```
bin/tictactoe 
```
![assets/player-2.png](assets/player-2.png)

- When game is over

![assets/game-over.png](assets/game-over.png)


