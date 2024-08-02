package commands

import (
	"github.com/bwmarrin/discordgo"
	"jtams/playertrackerbot/bot"
)

func PingCommand() *discordgo.ApplicationCommand {
	cmd := &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Ping Pong!",
	}

	return cmd
}

func PingHandler() bot.CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Pong!",
			},
		})

	}
}
