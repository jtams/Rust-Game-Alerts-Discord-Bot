package tracker

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Messenger is used to update the message with the tracker data
type Messenger struct {
	Session         *discordgo.Session
	Message         *discordgo.Message
	ChannelID       string
	content         string
	MessageOverflow int
}

// Creates a new Messenger with default settings
func NewMessageUpdater(session *discordgo.Session) *Messenger {
	return &Messenger{
		Session: session,
		Message: nil,
	}
}

// Starts tracking the player tracker and updates the message with the tracker data
// This is the entry point for starting the tracker. This will automatically start the
// PlayerTracker's update loop.
func (updater *Messenger) StartTracking(tracker *PlayerTracker) {
	// Start the tracker
	go tracker.Start()

	// Clears the content to force an update on the first run (used when tracking is stopped then stopped)
	updater.content = ""

	// Wait for tracker update
	for range tracker.Channel {
		if !tracker.IsRunning() {
			// Save to register that the tracker was stopped manually
			SaveTrackerData(os.Getenv("SAVE_FILE"), tracker, updater)
			logger.Info("Messenger shutting down.")
			return
		}

		// If the message is nil, create a new message.
		// If the /start command doesn't create a message for some reason.
		if updater.Message == nil {
			createNewMessage(updater, updater.Session, updater.ChannelID, "Tracker is starting up...")
			return
		}

		updater.ChannelID = updater.Message.ChannelID

		content := fmt.Sprintf("https://www.battlemetrics.com/servers/rust/%s\n```diff\n%s\n%d/%d Online\n\n", tracker.BattleMetricsID, tracker.ServerName, tracker.Online[0], tracker.Online[1])
		empty := true

		// Create list of groups with their players and their status
		for _, group := range tracker.Groups {
			if len(group.Users) == 0 {
				continue
			}
			empty = false

			playerList := ""

			for _, player := range group.Users {
				symbol := "-"
				if player.Status == StatusOnline {
					symbol = "+"
				}

				name := player.GetUsername()
				if player.Status == StatusUntracked || player.Status == StatusUnknown {
					name += " (waiting for player to come online)"
				}

				if player.Status == StatusOffline {
					name += fmt.Sprintf(" (last seen %s)", timeAgo(player.GetLastOnline()))
				}

				playerList += fmt.Sprintf("%s %s\n", symbol, name)
			}

			dividerCount := 20 - len(group.Name)/2
			divider1 := strings.Repeat("═", dividerCount)
			var divider2 string
			if len(group.Name)%2 == 1 {
				divider2 = strings.Repeat("═", dividerCount-1)
			} else {
				divider2 = strings.Repeat("═", dividerCount)
			}
			content += fmt.Sprintf("%s %s %s\n", divider1, strings.ToUpper(group.Name), divider2)
			content += playerList + "\n"
		}

		if empty {
			content += "- No users are being tracked. Run /users add [username] [group] to add a user to tracking."
		}

		content += "\n```"

		// No need to update, nothing changed
		if updater.content == content {
			continue
		}

		updater.content = content

		msgEdit := &discordgo.MessageEdit{
			Channel: updater.ChannelID,
			ID:      updater.Message.ID,
			Content: &content,
		}

		// If users have send a few messages since the last update, create a new message
		// so that the tracker is always visible in the channel.
		if updater.MessageOverflow > 4 {
			updater.Session.ChannelMessageDelete(updater.ChannelID, updater.Message.ID)
			message, _ := updater.Session.ChannelMessageSend(updater.ChannelID, content)
			updater.Message = message
			updater.MessageOverflow = 0
		} else {
			_, err := updater.Session.ChannelMessageEditComplex(msgEdit)
			// If the message was deleted, create a new message
			if err != nil {
				createNewMessage(updater, updater.Session, updater.ChannelID, content)
			}
		}

		// Save
		if err := SaveTrackerData(os.Getenv("SAVE_FILE"), tracker, updater); err != nil {
			logger.Error("Failed to save tracker data: ", err)
		}
	}
}

// Creates a new message in the channel with the content
func createNewMessage(messenger *Messenger, session *discordgo.Session, channelID string, content string) error {
	message, err := session.ChannelMessageSend(channelID, content)
	if err != nil {
		return err
	}

	messenger.Message = message

	return nil
}

// Converts a time to a string that represents how long ago it was
// Includes, seconds, minutes, hours, and days
func timeAgo(t time.Time) string {
	duration := time.Since(t)
	seconds := int(duration.Seconds())
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24

	if days > 0 {
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
	if hours > 0 {
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	if minutes > 0 {
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	if seconds == 1 {
		return "1 second ago"
	}
	return fmt.Sprintf("%d seconds ago", seconds)
}
