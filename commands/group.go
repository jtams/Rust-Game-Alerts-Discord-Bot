package commands

import (
	"jtams/playertrackerbot/bot"
	"jtams/playertrackerbot/tracker"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func GroupCommand() *discordgo.ApplicationCommand {
	cmd := &discordgo.ApplicationCommand{
		Name:        "group",
		Description: "manage groups",
		Options: []*discordgo.ApplicationCommandOption{

			{
				Name:        "add",
				Description: "Add group",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "name",
						Description: "Name of the group to add",
						Required:    true,
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
						Name:        "name",
						Description: "Name of the group to remove",
						Required:    true,
					},
				},
			},
		},
	}

	return cmd
}

func GroupHandler(messageTracker *tracker.Messenger, playerTracker *tracker.PlayerTracker, registry bot.CommandRegistry) bot.CommandHandler {
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

		switch options[0].Name {
		case "add":
			options = options[0].Options
			groupNameRaw := findOptionByName("name", options)
			groupName := ""
			if groupNameRaw != nil {
				groupName = groupNameRaw.StringValue()
			}
			groupName = strings.ToLower(groupName)
			res = addGroup(playerTracker, groupName, messageTracker, registry)
			break
		case "remove":
			options = options[0].Options
			groupNameRaw := findOptionByName("name", options)
			groupName := ""
			if groupNameRaw != nil {
				groupName = groupNameRaw.StringValue()
			}
			groupName = strings.ToLower(groupName)
			res = removeGroup(playerTracker, groupName, messageTracker, registry)
		}

		_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: res,
		})

		if err != nil {
			log.Println(err)
		}

		return nil
	}
}

func updateCommand(registry bot.CommandRegistry, messageUpdater *tracker.Messenger, playerTracker *tracker.PlayerTracker) error {
	groups := []string{}
	for _, group := range playerTracker.Groups {
		groups = append(groups, group.Name)
	}
	if err := registry.UpdateCommand(*UserCommand(groups), UserHandler(messageUpdater, playerTracker)); err != nil {
		return err
	}

	err := registry.Register()
	if err != nil {
		return err
	}

	return nil
}

func addGroup(playerTracker *tracker.PlayerTracker, groupName string, messageUpdater *tracker.Messenger, registry bot.CommandRegistry) string {
	if groupName == "" {
		return "Invalid group name"
	}

	if err := playerTracker.AddNewGroup(groupName); err != nil {
		return "Failed to add group"
	}

	err := updateCommand(registry, messageUpdater, playerTracker)
	if err != nil {
		log.Println(err)
		return "Failed to update command"
	}

	return "Group added"
}

func removeGroup(playerTracker *tracker.PlayerTracker, groupName string, messageUpdater *tracker.Messenger, registry bot.CommandRegistry) string {
	if groupName == "" {
		return "Invalid group name"
	}

	if playerTracker.RemoveGroup(groupName) {
		return "Group removed"
	}

	err := updateCommand(registry, messageUpdater, playerTracker)
	if err != nil {
		log.Println(err)
		return "Failed to update command"
	}

	return "Group not found"
}
