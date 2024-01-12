package discord

import (
	database "discordbot/pkg/db"
	"discordbot/pkg/discord/command"
	"discordbot/pkg/discord/job"
	"discordbot/pkg/model"
	"discordbot/pkg/repo"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type DiscordClient struct {
	session      *discordgo.Session
	DBMgr        *database.DatabaseManager
	logger       *zap.SugaredLogger
	jobScheduler *job.JobScheduler
}

func Initialize(session *discordgo.Session, dbmgr *database.DatabaseManager, logger *zap.SugaredLogger, jobScheduler *job.JobScheduler) (*DiscordClient, error) {
	return &DiscordClient{
		session:      session,
		DBMgr:        dbmgr,
		logger:       logger,
		jobScheduler: jobScheduler,
	}, nil
}

func (c *DiscordClient) Connect() (err error) {
	defer c.logger.Sync()

	// Start job scheduler
	c.jobScheduler.Start()

	// Register event handlers
	c.session.AddHandler(c.onReady)
	c.session.AddHandler(c.messageCreate)

	// ApplicationCommand
	command.AddCommandHandlers(c.session)

	err = c.session.Open()
	if err != nil {
		c.logger.Infoln("error opening connection,", err)
		return
	}

	return
}

func (c *DiscordClient) Close() {
	c.jobScheduler.Stop()
	command.CleanUpCommands(c.session)
	c.session.Close()
}

func (c *DiscordClient) onReady(s *discordgo.Session, r *discordgo.Ready) {
	logger := zap.L().Sugar()
	defer logger.Sync()
	logger.Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)

	// Register ApplicationCommand
	command.RegisterCommands(s)
}

func (c *DiscordClient) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	userRepo := repo.NewUserRepo(c.DBMgr.DB)
	user := userRepo.FindOrCreate(&model.User{
		ID: m.Author.ID,
	})
	c.logger.Debugf("user: %v", user)

	err := userRepo.UpdateById(user.ID, &model.User{
		Username: m.Author.Username,
	})
	if err != nil {
		c.logger.Errorf("error updating user: %v", err)
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
