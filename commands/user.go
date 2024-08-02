package commands

import (
	"jtams/playertrackerbot/bot"
	"jtams/playertrackerbot/tracker"

	"github.com/bwmarrin/discordgo"
)

func UserCommand() *discordgo.ApplicationCommand {
	cmd := &discordgo.ApplicationCommand{
		Name:        "users",
		Description: "manage users",
		Options: []*discordgo.ApplicationCommandOption{

			{
				Name:        "add",
				Description: "Add user(s)",
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "username",
						Description: "Username of the user to add",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "group",
						Description: "Group to add the user to",
						Required:    false,
					},
				},
			},

			{
				Name:        "remove",
				Description: "Remove user(s)",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "username",
						Description: "Username of the user to remove",
						Required:    true,
					},
				},
			},
		},
	}

	return cmd
}

func UserHandler(messageTracker *tracker.Messenger, playerTracker *tracker.PlayerTracker) bot.CommandHandler {
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

		switch options[0].Name {
		case "add":
			options = options[0].Options
			username := findOptionByName("username", options).StringValue()
			groupName := findOptionByName("group", options).StringValue()
			res = addUser(playerTracker, username, groupName)
			break
		case "remove":
			options = options[0].Options
			username := findOptionByName("username", options).StringValue()
			res = removeUser(playerTracker, username)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: res,
			},
		})

		return nil
	}
}

func addUser(playerTracker *tracker.PlayerTracker, username string, groupName string) string {
	if groupName == "" {
		groupName = "Others"
	}

	err := playerTracker.AddUserToGroup(username, groupName)
	if err != nil {
		return "Error adding user to group"
	}

	return "User added to group"
}

func removeUser(playerTracker *tracker.PlayerTracker, username string) string {
	if playerTracker.RemoveUserByUsername(username) {
		return "User removed"
	} else {
		return "User not found"
	}
}
