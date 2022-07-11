/*
Copyright © 2022 Harsh Varagiya

*/
package cmd

import (
	"fmt"
	"net"
	"ticTacToe/server"
	"ticTacToe/structures"
	"ticTacToe/structures/tictactoe"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var grpcPort = 8888

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A simple tic-tac-toe game server",
	Run: func(cmd *cobra.Command, args []string) {
		var gameServer server.GameServer
		gameServer.Games = make(map[string]*structures.Game)
		srv := grpc.NewServer()
		tictactoe.RegisterTicTacToeGameServer(srv, gameServer)
		log.Debugf("attempting to listen on port: %v", grpcPort)
		listener, err := net.Listen("tcp", fmt.Sprintf(":%v", grpcPort))
		if err != nil {
			log.Fatal("error listening to port. error = %v", err)
		}
		log.Infof("service running on port: %v", grpcPort)
		log.Fatal(srv.Serve(listener))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
