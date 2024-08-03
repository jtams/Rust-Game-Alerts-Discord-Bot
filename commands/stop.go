package commands

import (
	"jtams/playertrackerbot/bot"
	"jtams/playertrackerbot/tracker"

	"github.com/bwmarrin/discordgo"
)

func StopCommand() *discordgo.ApplicationCommand {
	cmd := &discordgo.ApplicationCommand{
		Name:        "stop",
		Description: "Stops the tracker",
	}

	return cmd
}

func StopHandler(messageTracker *tracker.Messenger, playerTracker *tracker.PlayerTracker) bot.CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		if messageTracker.Message != nil {
			content := "Stopped."

			msgEdit := &discordgo.MessageEdit{
				Channel: messageTracker.Message.ChannelID,
				ID:      messageTracker.Message.ID,
				Content: &content,
			}

			messageTracker.Session.ChannelMessageEditComplex(msgEdit)
		}

		playerTracker.Stop()

		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Stopped",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

	}
}
