/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"discordbot/pkg/util"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Debug bool
var Config util.Config

// rootCmd represents the run command
var rootCmd = &cobra.Command{
	Use:   "discordbot",
	Short: "A Discord bot",
	Long:  `Discord bot with Database and Job Scheduler support.`,
}

func Execute() {
	err := runCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(util.LoadConfig)

	rootCmd.PersistentFlags().BoolVarP(&Config.Debug, "debug", "d", false, "Display debugging output in the console. (default: false)")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}
