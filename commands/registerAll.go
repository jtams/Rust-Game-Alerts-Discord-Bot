package commands

import (
	"jtams/playertrackerbot/tracker"

	"github.com/bwmarrin/discordgo"
)

// Registers default commands.
// Including: Start, Stop, Group, User
func RegisterAllDefaultCommands(discord *discordgo.Session, registry CommandRegistry, tracker *tracker.PlayerTracker, messageUpdater *tracker.Messenger) {
	// Start Command
	if err := registry.AddCommand(*StartCommand(), StartHandler(messageUpdater, tracker)); err != nil {
		logger.Error("Failed to register command", "name", "start", "error", err)
	}

	// Stop Command
	if err := registry.AddCommand(*StopCommand(), StopHandler(messageUpdater, tracker)); err != nil {
		logger.Error("Failed to register command", "name", "stop", "error", err)
	}

	// Info Command
	if err := registry.AddCommand(*InfoCommand(), InfoHandler(messageUpdater, tracker)); err != nil {
		logger.Error("Failed to register command", "name", "info", "error", err)
	}

	// Get list of groups. Used for group options
	groups := []string{}
	for _, group := range tracker.Groups {
		groups = append(groups, group.Name)
	}

	// Group Command
	if err := registry.AddCommand(*GroupCommand(groups), GroupHandler(messageUpdater, tracker, registry)); err != nil {
		logger.Error("Failed to register command", "name", "group", "error", err)
	}

	// User Command
	if err := registry.AddCommand(*UserCommand(groups), UserHandler(messageUpdater, tracker)); err != nil {
		logger.Error("Failed to register command", "name", "user", "error", err)
	}

	// Register comamnds with Discord
	err := registry.Register()
	if err != nil {
		logger.Error("Failed to register commands", "error", err)
	}
}
