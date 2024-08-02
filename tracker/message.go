package tracker

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Messenger struct {
	Session   *discordgo.Session
	Message   *discordgo.Message
	ChannelID string
	content   string
}

func NewMessageUpdater(session *discordgo.Session) *Messenger {
	return &Messenger{
		Session: session,
		Message: nil,
	}
}

func (updater *Messenger) StartTracking(tracker *PlayerTracker) {
	go tracker.Start()

	for range tracker.Channel {
		if !tracker.IsRunning() {
			// Save to register that the tracker was stopped manually
			SaveTrackerData(os.Getenv("SAVE_FILE"), tracker, updater)
			log.Println("Messenger shutting down.")
			return
		}

		if updater.Message == nil {
			log.Println("Message is nil, cannot update")
			return
		}

		updater.ChannelID = updater.Message.ChannelID

		content := fmt.Sprintf("https://www.battlemetrics.com/servers/rust/%s\n```diff\n%s\n%d/%d Online\n\n", tracker.BattleMetricsID, tracker.ServerName, tracker.Online[0], tracker.Online[1])
		empty := true

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

			content += fmt.Sprintf("══════════════ %s ══════════════\n", strings.ToUpper(group.Name))
			content += playerList
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
			Channel: updater.Message.ChannelID,
			ID:      updater.Message.ID,
			Content: &content,
		}

		updater.Session.ChannelMessageEditComplex(msgEdit)

		if err := SaveTrackerData(os.Getenv("SAVE_FILE"), tracker, updater); err != nil {
			log.Println("Failed to save file", err)
		}
	}
}

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

// https://www.battlemetrics.com/servers/rust/18566638
// ```diff
// [US West] Facepunch 2 (online)
// 18/250 Online
//
// ══════════════ SQUAD ══════════════
// - realm (untracked)
// - thelowerrealm (untracked)
// - void (untracked)
//
// ```
