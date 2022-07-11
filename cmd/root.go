/*
Copyright © 2022 Harsh Varagiya

*/
package cmd

import (
	"os"
	"ticTacToe/client"
	"ticTacToe/structures/tictactoe"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

var cfgFile = "tic-tac-toe-client.yaml"
var config client.ServerConfiguration

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ticTacToe",
	Short: "A simple 2 player Tic-Tac-Toe game ;)",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PreRunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile("tic-tac-toe.yaml")
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			return err
		}
		log.Infof("Player '%v' connecting to server '%v'", config.PlayerName, config.ServerHost)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(config.ServerHost, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect to game server. error = %v", err)
		}
		gameClient := tictactoe.NewTicTacToeGameClient(conn)
		client.PlayGame(&config, &gameClient)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is tic-tac-toe.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
