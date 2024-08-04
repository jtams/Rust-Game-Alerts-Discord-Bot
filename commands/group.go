package commands

import (
	"jtams/playertrackerbot/tracker"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Group command is used to add and remove groups
func GroupCommand(groups []string) *discordgo.ApplicationCommand {
	groupChoices := make([]*discordgo.ApplicationCommandOptionChoice, len(groups))
	for i, group := range groups {
		groupChoices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  group,
			Value: group,
		}
	}

	cmd := &discordgo.ApplicationCommand{
		Name:        "groups",
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
				Description: "Remove group",
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

			{
				Name:        "notes",
				Description: "Set notes for a group",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "group",
						Description: "Group to set notes for",
						Choices:     groupChoices,
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "notes",
						Description: "Notes to set",
						Required:    false,
					},
				},
			},

			{
				Name:        "location",
				Description: "Sets location for a group",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "group",
						Description: "Group to set notes for",
						Choices:     groupChoices,
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "location",
						Description: "Location of the group",
						Required:    false,
					},
				},
			},
		},
	}

	return cmd
}

func GroupHandler(messageTracker *tracker.Messenger, playerTracker *tracker.PlayerTracker, registry CommandRegistry) CommandHandler {
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
		case "notes":
			options = options[0].Options
			var notes string
			noteOption := findOptionByName("notes", options)
			if noteOption == nil {
				notes = ""
			} else {
				notes = noteOption.StringValue()
			}
			groupNameRaw := findOptionByName("group", options)
			groupName := ""
			if groupNameRaw != nil {
				groupName = groupNameRaw.StringValue()
			}
			groupName = strings.ToLower(groupName)
			res = setGroupNotes(playerTracker, notes, groupName)
		case "location":
			options = options[0].Options
			var location string
			noteOption := findOptionByName("location", options)
			if noteOption == nil {
				location = ""
			} else {
				location = noteOption.StringValue()
			}
			groupNameRaw := findOptionByName("group", options)
			groupName := ""
			if groupNameRaw != nil {
				groupName = groupNameRaw.StringValue()
			}
			groupName = strings.ToLower(groupName)
			res = setGroupLocation(playerTracker, location, groupName)
		}

		_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: res,
		})

		if err != nil {
			logger.Error("Failed to send followup message", "error", err)
		}

		playerTracker.Channel <- time.Now()

		return nil
	}
}

func updateCommand(registry CommandRegistry, messageUpdater *tracker.Messenger, playerTracker *tracker.PlayerTracker) error {
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

func addGroup(playerTracker *tracker.PlayerTracker, groupName string, messageUpdater *tracker.Messenger, registry CommandRegistry) string {
	if groupName == "" {
		return "Invalid group name"
	}

	if err := playerTracker.AddNewGroup(groupName); err != nil {
		return "Failed to add group"
	}

	err := updateCommand(registry, messageUpdater, playerTracker)
	if err != nil {
		logger.Error("Failed to update command", "error", err)
		return "Failed to update command"
	}

	return "Group added"
}

func removeGroup(playerTracker *tracker.PlayerTracker, groupName string, messageUpdater *tracker.Messenger, registry CommandRegistry) string {
	if groupName == "" {
		return "Invalid group name"
	}

	if playerTracker.RemoveGroup(groupName) {
		return "Group removed"
	}

	err := updateCommand(registry, messageUpdater, playerTracker)
	if err != nil {
		logger.Error("Failed to update command", "error", err)
		return "Failed to update command"
	}

	return "Group not found"
}

func setGroupNotes(playerTracker *tracker.PlayerTracker, notes string, groupName string) string {
	if groupName == "" {
		groupName = "others"
	}

	group := playerTracker.GetGroupByName(groupName)
	if group == nil {
		return "Group not found"
	}

	group.Notes = notes
	return "Notes set"
}

func setGroupLocation(playerTracker *tracker.PlayerTracker, location string, groupName string) string {
	if groupName == "" {
		groupName = "others"
	}

	group := playerTracker.GetGroupByName(groupName)
	if group == nil {
		return "Group not found"
	}

	group.Location = location
	return "Location set"
}
