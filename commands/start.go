package commands

import (
	"errors"
	"jtams/playertrackerbot/tracker"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Start command starts the tracker and also sets the BattleMetrics ID if provided
func StartCommand() *discordgo.ApplicationCommand {
	cmd := &discordgo.ApplicationCommand{
		Name:        "start",
		Description: "Start the tracker",
		Options: []*discordgo.ApplicationCommandOption{

			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "battle_metrics_id",
				Description: "The BattleMetrics ID of the server to track. Paste the ID or the full URL of the BattleMetrics page.",
				Required:    false,
			},
		},
	}

	return cmd
}

// Finds the ID in the URL
func extractID(url string) (string, error) {
	// Searches for number
	re := regexp.MustCompile(`\b\d+\b`)

	match := re.FindString(url)
	if match == "" {
		return "", errors.New("no match found")
	}

	return match, nil
}

func StartHandler(messageTracker *tracker.Messenger, playerTracker *tracker.PlayerTracker) CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		// Check if the tracker is already running
		if playerTracker.IsRunning() {
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Tracker is already running.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}

		// Check for battle_metrics_id
		battleMetricsID := ""
		for _, option := range i.ApplicationCommandData().Options {
			if option.Name == "battle_metrics_id" {
				var err error
				battleMetricsID = option.StringValue()
				battleMetricsID, err = extractID(battleMetricsID)
				if err != nil {
					return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Invalid BattleMetrics ID provided.",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					})
				}
			}
		}

		// Require battle_metrics_id if it's the first time starting the tracker
		if battleMetricsID == "" && playerTracker.BattleMetricsID == "" {
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You must provide a BattleMetrics ID to start the tracker for the first time.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}

		// Set new BattleMetrics ID if provided
		if battleMetricsID != "" {
			playerTracker.BattleMetricsID = battleMetricsID
		}

		tryEditing := func() error {
			if messageTracker.Message != nil {
				content := "Stopped."

				msgEdit := &discordgo.MessageEdit{
					Channel: messageTracker.Message.ChannelID,
					ID:      messageTracker.Message.ID,
					Content: &content,
				}

				_, err := messageTracker.Session.ChannelMessageEditComplex(msgEdit)
				return err
			}

			return errors.New("no message to edit")
		}

		err := tryEditing()
		if err != nil {
			message, _ := s.ChannelMessageSend(i.ChannelID, "Starting...")
			messageTracker.Message = message

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Starting...",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Restarting...",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}

		playerTracker.Channel = make(chan time.Time)
		go messageTracker.StartTracking(playerTracker)

		return err
	}
}

// Used if the bot was shutdown while tracking
func ForceStartup(messageTracker *tracker.Messenger, playerTracker *tracker.PlayerTracker) {
	playerTracker.Channel = make(chan time.Time)
	go messageTracker.StartTracking(playerTracker)
}
