package commands

import (
	"jtams/playertrackerbot/tracker"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// User command manages users.
// Includes: Add, Remove, Move
func UserCommand(groups []string) *discordgo.ApplicationCommand {
	groupChoices := make([]*discordgo.ApplicationCommandOptionChoice, len(groups))
	for i, group := range groups {
		groupChoices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  group,
			Value: group,
		}
	}

	cmd := &discordgo.ApplicationCommand{
		Name:        "users",
		Description: "manage users",
		Options: []*discordgo.ApplicationCommandOption{

			{
				Name:        "add",
				Description: "Add user(s)",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "username",
						Description: "Username(s) of the user(s) to add. Seperate multiple with commas",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "group",
						Description: "Group to add the user to",
						Required:    false,
						Choices:     groupChoices,
					},
				},
			},

			{
				Name:        "add-by-id",
				Description: "Add user by BattleMetric Player ID",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "id",
						Description: "ID of the user to add.",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "group",
						Description: "Group to add the user to",
						Required:    false,
						Choices:     groupChoices,
					},
				},
			},

			{
				Name:        "remove",
				Description: "Remove user(s)",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "username",
						Description: "Username(s) of the user(s) to remove",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "group",
						Description: "Group to remove the user from, if not provided will remove from all groups",
						Required:    false,
						Choices:     groupChoices,
					},
				},
			},
			{
				Name:        "move",
				Description: "Move user(s)",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "username",
						Description: "Username(s) of the user(s) to move",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "group",
						Description: "Group to move the user to",
						Required:    false,
						Choices:     groupChoices,
					},
				},
			},
		},
	}

	return cmd
}

func UserHandler(messageTracker *tracker.Messenger, playerTracker *tracker.PlayerTracker) CommandHandler {
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

		defer func() {
			if recover() != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Failed.",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
			}
		}()

		var res string

		switch options[0].Name {
		case "add":
			options = options[0].Options
			username := findOptionByName("username", options).StringValue()
			groupNameRaw := findOptionByName("group", options)
			groupName := ""
			if groupNameRaw != nil {
				groupName = groupNameRaw.StringValue()
			}
			groupName = strings.ToLower(groupName)
			res = addUser(playerTracker, username, groupName)
			break
		case "add-by-id":
			options = options[0].Options
			id := findOptionByName("id", options).StringValue()
			groupNameRaw := findOptionByName("group", options)
			groupName := ""
			if groupNameRaw != nil {
				groupName = groupNameRaw.StringValue()
			}
			groupName = strings.ToLower(groupName)
			res = addUserByID(playerTracker, id, groupName)
			break
		case "remove":
			options = options[0].Options
			username := findOptionByName("username", options).StringValue()
			groupNameRaw := findOptionByName("group", options)
			groupName := ""
			if groupNameRaw != nil {
				groupName = groupNameRaw.StringValue()
			}
			groupName = strings.ToLower(groupName)
			res = removeUser(playerTracker, username, groupName)
		case "move":
			options = options[0].Options
			username := findOptionByName("username", options).StringValue()
			groupNameRaw := findOptionByName("group", options)
			groupName := ""
			if groupNameRaw != nil {
				groupName = groupNameRaw.StringValue()
			}
			groupName = strings.ToLower(groupName)
			res = moveUser(playerTracker, username, groupName)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: res,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		// Immediately update the player tracker
		playerTracker.Channel <- time.Now()

		return nil
	}
}

func addUser(playerTracker *tracker.PlayerTracker, username string, groupName string) string {
	if groupName == "" {
		groupName = "others"
	}

	failed := []string{}

	usernames := strings.Split(username, ",")
	for _, username := range usernames {
		username = strings.TrimSpace(username)
		err := playerTracker.AddUserToGroup(username, groupName)
		if err != nil {
			failed = append(failed, username)
			return "Error adding user to group"
		}
	}

	if len(failed) == 0 {
		return "User(s) added"
	}

	joined := strings.Join(failed, ", ")
	return "Failed to add user(s): " + joined
}

func addUserByID(playerTracker *tracker.PlayerTracker, id string, groupName string) string {
	if groupName == "" {
		groupName = "others"
	}

	failed := []string{}

	id = strings.TrimSpace(id)
	err := playerTracker.AddUserToGroupByID(id, groupName)
	if err != nil {
		failed = append(failed, id)
		return "Error adding user to group"
	}

	if len(failed) == 0 {
		return "User added to list. Please wait for the tracker to update."
	}

	joined := strings.Join(failed, ", ")
	return "Failed to add user " + joined
}

func removeUser(playerTracker *tracker.PlayerTracker, username string, groupName string) string {
	users := strings.Split(username, ",")

	failed := []string{}

	for _, user := range users {
		user = strings.TrimSpace(user)
		if groupName == "" {
			if !playerTracker.RemoveUserByUsername(user) {
				failed = append(failed, user)
			}
		} else {
			if !playerTracker.RemoveUserByUsernameAndGroup(user, groupName) {
				failed = append(failed, user)
			}
		}
	}

	if len(failed) == 0 {
		return "User(s) removed"
	}

	joined := strings.Join(failed, ", ")
	return "Failed to remove user(s): " + joined
}

func moveUser(playerTracker *tracker.PlayerTracker, username string, groupName string) string {
	users := strings.Split(username, ",")

	failed := []string{}

	for _, user := range users {
		user = strings.TrimSpace(user)
		if !playerTracker.MoveUserToGroup(user, groupName) {
			failed = append(failed, user)
		}
	}

	if len(failed) == 0 {
		return "User(s) moved"
	}

	joined := strings.Join(failed, ", ")
	return "Failed to move user(s): " + joined
}
