package command

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type Command struct {
	ApplicationCommand discordgo.ApplicationCommand
	Handler            func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func getCommands() (commands []*Command) {
	// add your commands here
	commands = append(commands, &Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "basic-command",
			Description: "Basic command",
		},
		Handler: basicCommand,
	})

	return
}

func RegisterCommands(s *discordgo.Session) {
	var guilds []string
	commands := getCommands()

	logger := zap.L().Sugar()
	defer logger.Sync()
	logger.Infoln("Adding commands...")
	for _, g := range s.State.Guilds {
		logger.Debugf("Add for guild: %v", g.ID)
		guilds = append(guilds, g.ID)
		for _, v := range commands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, g.ID, &v.ApplicationCommand)
			if err != nil {
				logger.Errorf("Cannot create '%v' command: %v", v.ApplicationCommand.Name, err)
			}
		}
	}
	logger.Infof("Registered commands on %v guild(s).", len(guilds))
}

func AddCommandHandlers(s *discordgo.Session) {
	commandHandlers := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	commands := getCommands()

	for _, v := range commands {
		commandHandlers[v.ApplicationCommand.Name] = v.Handler
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		logger := zap.L().Sugar()
		defer logger.Sync()
		logger.Infof("Received interaction: %v", i.ApplicationCommandData().Name)

		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func CleanUpCommands(s *discordgo.Session) {
	logger := zap.L().Sugar()
	defer logger.Sync()
	logger.Infof("Cleaning up commands...")

	for _, g := range s.State.Guilds {
		logger.Debugf("Cleaning up guild: %v", g.ID)
		DeleteAllCommands(s, g.ID)
	}
}

func DeleteCommand(s *discordgo.Session, guildID string, command string) {
	logger := zap.L().Sugar()
	defer logger.Sync()
	logger.Debugf("Deleting command '%v' on guild '%v'", command, guildID)

	err := s.ApplicationCommandDelete(s.State.User.ID, guildID, command)
	if err != nil {
		logger.Errorf("Cannot delete '%v' command: %v", command, err)
	}
}

func DeleteAllCommands(s *discordgo.Session, guildID string) {
	logger := zap.L().Sugar()
	defer logger.Sync()
	logger.Debugf("Deleting all commands on guild '%v'", guildID)

	commands, err := s.ApplicationCommands(s.State.User.ID, guildID)
	if err != nil {
		logger.Errorf("Cannot get commands: %v", err)
		return
	}

	for _, v := range commands {
		DeleteCommand(s, guildID, v.ID)
	}
}
