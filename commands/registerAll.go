package commands

import (
	"jtams/playertrackerbot/tracker"
	"log"

	"github.com/bwmarrin/discordgo"
)

// Registers default commands.
// Including: Start, Stop, Group, User
func RegisterAllDefaultCommands(discord *discordgo.Session, registry CommandRegistry, tracker *tracker.PlayerTracker, messageUpdater *tracker.Messenger) {
	// Start Command
	if err := registry.AddCommand(*StartCommand(), StartHandler(messageUpdater, tracker)); err != nil {
		log.Println(err)
	}

	// Stop Command
	if err := registry.AddCommand(*StopCommand(), StopHandler(messageUpdater, tracker)); err != nil {
		log.Println(err)
	}

	// Group Command
	if err := registry.AddCommand(*GroupCommand(), GroupHandler(messageUpdater, tracker, registry)); err != nil {
		log.Println(err)
	}

	// Get list of groups. Used for group options
	groups := []string{}
	for _, group := range tracker.Groups {
		groups = append(groups, group.Name)
	}

	// User Command
	if err := registry.AddCommand(*UserCommand(groups), UserHandler(messageUpdater, tracker)); err != nil {
		log.Println(err)
	}

	// Register comamnds with Discord
	err := registry.Register()
	if err != nil {
		log.Println(err)
	}
}
