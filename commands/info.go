package commands

import (
	"encoding/json"
	"fmt"
	"jtams/playertrackerbot/tracker"
	"time"

	"github.com/bwmarrin/discordgo"
)

func InfoCommand() *discordgo.ApplicationCommand {
	cmd := &discordgo.ApplicationCommand{
		Name:        "info",
		Description: "Get information about the tracker",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "user",
				Description: "Get information about a user",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "username",
						Description: "Username of the user to get information about.",
						Required:    true,
					},
				},
			},
			{
				Name:        "group",
				Description: "Get information about a group",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "group",
						Description: "Name of the group to get information about.",
						Required:    true,
					},
				},
			},
		},
	}

	return cmd
}

func InfoHandler(messageTracker *tracker.Messenger, playerTracker *tracker.PlayerTracker) CommandHandler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		options := i.ApplicationCommandData().Options

		findOptionByName := func(name string, opts []*discordgo.ApplicationCommandInteractionDataOption) *discordgo.ApplicationCommandInteractionDataOption {
			for _, opt := range opts {
				if opt.Name == name {
					return opt
				}
			}

			return nil
		}

		var res string
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})

		if len(options) == 0 {
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Failed, no options provided",
			})
		}

		defer func() {
			if recover() != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Failed",
				})
			}
		}()

		switch options[0].Name {
		case "user":
			options = options[0].Options
			username := findOptionByName("username", options).StringValue()
			res = getPlayerInfo(playerTracker, username)
			break
		case "group":
			options = options[0].Options
			group := findOptionByName("group", options).StringValue()
			res = getGroupInfo(playerTracker, group)
			break
		}

		_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: res,
		})

		if err != nil {
			_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Failed",
			})
		}

		if err != nil {
			logger.Error("Failed to send followup message", "error", err)
		}

		// Immediately update the player tracker
		playerTracker.Channel <- time.Now()

		return nil
	}
}

type UserWithoutTimes struct {
	ID        string             `json:"id"`
	Usernames []string           `json:"usernames"`
	Group     string             `json:"group"`
	Status    tracker.UserStatus `json:"status"`
}

func getPlayerInfo(playerTracker *tracker.PlayerTracker, username string) string {
	player := playerTracker.GetUserByUsername(username)
	if player == nil {
		return fmt.Sprintf("Player %s not found", username)
	}

	playerCopy := UserWithoutTimes{
		ID:        player.ID,
		Usernames: player.Usernames,
		Group:     player.Group,
		Status:    player.Status,
	}

	playerStr, err := json.MarshalIndent(playerCopy, "", " ")
	if err != nil {
		return "Failed to get player info"
	}

	return "https://www.battlemetrics.com/players/" + playerCopy.ID + "\n```json\n" + string(playerStr) + "```"
}

type UserShrunk struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type GroupShrunk struct {
	Name     string       `json:"name"`
	Users    []UserShrunk `json:"users"`
	Notes    string       `json:"notes"`
	Location string       `json:"location"`
}

func getGroupInfo(playerTracker *tracker.PlayerTracker, groupName string) string {
	group := playerTracker.GetGroupByName(groupName)
	if group == nil {
		return fmt.Sprintf("Group %s not found", groupName)
	}

	strippedUserList := []UserShrunk{}
	for _, user := range group.Users {
		strippedUserList = append(strippedUserList, UserShrunk{
			ID:       user.ID,
			Username: user.GetUsername(),
		})
	}

	groupShrunk := GroupShrunk{
		Name:     group.Name,
		Users:    strippedUserList,
		Notes:    group.Notes,
		Location: group.Location,
	}

	groupStr, err := json.MarshalIndent(groupShrunk, "", " ")
	if err != nil {
		return "Failed to get group info"
	}

	return "```json\n" + string(groupStr) + "```"

}
