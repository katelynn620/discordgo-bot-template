/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	database "discordbot/pkg/db"
	"discordbot/pkg/discord"
	"discordbot/pkg/discord/job"
	"discordbot/pkg/util"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// runCmd represents the base command when called without any subcommands
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		logger := zap.L().Sugar()
		defer logger.Sync()

		debug, _ := cmd.Flags().GetBool("debug")
		logger.Infof("debug=%v", debug)

		token, _ := cmd.Flags().GetString("token")
		if token == "" {
			token = viper.GetString("token")
		}
		if token == "" {
			logger.Panicln("token is required")
			return
		}

		var (
			session *discordgo.Session
			c       *discord.DiscordClient
			err     error
		)

		dbMgr, err := database.InitDatabaseManager()
		if err != nil {
			logger.Panicf("failed to connect database: %v", err)
			return
		}
		dbMgr.Migrate()

		session, err = discordgo.New("Bot " + token)
		if err != nil {
			logger.Errorf("error creating Discord session,", err)
			return
		}

		jobScheduler := job.Initialize(session)

		c, err = discord.Initialize(session, dbMgr, logger, jobScheduler)
		if err != nil {
			logger.Errorf("error creating Discord client,", err)
			return
		}

		err = c.Connect()
		if err != nil {
			logger.Panicf("error connecting Discord,", err)
			return
		}

		logger.Infoln("Bot is now running.  Press CTRL-C to exit.")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc

		logger.Infoln("Bot is quitting.")
		// Cleanly close down the Discord session.
		c.Close()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	cobra.OnInitialize(func() {
		// Load config
		config := util.GetConfig()

		if rootCmd.PersistentFlags().Lookup("debug").Changed {
			config.Debug, _ = rootCmd.PersistentFlags().GetBool("debug")
		}
		if config.Debug {
			config.Log.Level = "debug"
		}

		// Initialize the log
		lc := util.LogConfig{
			Level:      config.Log.Level,
			FileName:   filepath.Join(".", config.Log.Dir, fmt.Sprintf("%v.log", time.Now().Unix())),
			MaxSize:    config.Log.MaxSize,
			MaxBackups: config.Log.MaxBackups,
			MaxAge:     config.Log.MaxAge,
		}
		err := util.InitLogger(lc)
		if err != nil {
			panic(err)
		}
	})

	runCmd.Flags().StringP("token", "t", "", "Discord token")
}
