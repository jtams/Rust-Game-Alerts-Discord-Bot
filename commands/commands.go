package commands

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var logger *slog.Logger = slog.Default()

// Function that gets called when associated command it triggered
type CommandHandler func(*discordgo.Session, *discordgo.InteractionCreate) error

// Manages bot commands
type CommandRegistry interface {
	// Adds command to registry
	AddCommand(discordgo.ApplicationCommand, CommandHandler) error

	// Updates existing command in registry
	UpdateCommand(discordgo.ApplicationCommand, CommandHandler) error

	// Removes command from registry
	RemoveCommand(commandName string) error

	// Registers command with Discord
	Register() error

	// Unregisters command with Discord
	Unregister() error

	// Run a commmands handler
	Run(*discordgo.InteractionCreate) error
}

// Stores data for command
type CommandData struct {
	name       string
	command    *discordgo.ApplicationCommand
	handler    CommandHandler
	registered bool
}

// Default implementation of CommandRegistry
type DefaultCommandRegistry struct {
	commands map[string]*CommandData
	session  *discordgo.Session
	guildID  string
}

// Creates a new Command Registry
func NewDefaultCommandRegistry(session *discordgo.Session, guildID string) *DefaultCommandRegistry {
	return &DefaultCommandRegistry{
		session:  session,
		commands: make(map[string]*CommandData),
		guildID:  guildID,
	}
}

// Adds a command to the registry
func (r *DefaultCommandRegistry) AddCommand(command discordgo.ApplicationCommand, handler CommandHandler) error {
	if _, found := r.commands[command.Name]; found {
		return errors.New("command already registered")
	}

	cmd := CommandData{
		name:       command.Name,
		command:    &command,
		handler:    handler,
		registered: false,
	}

	r.commands[cmd.name] = &cmd
	return nil
}

// Updates an existing command in the registry
func (r *DefaultCommandRegistry) UpdateCommand(command discordgo.ApplicationCommand, handler CommandHandler) error {
	if _, found := r.commands[command.Name]; !found {
		return errors.New("command not found")
	}

	cmd := CommandData{
		name:       command.Name,
		command:    &command,
		handler:    handler,
		registered: false,
	}

	r.commands[cmd.name] = &cmd
	return nil
}

// Removes a command from the registry
func (r *DefaultCommandRegistry) RemoveCommand(commandName string) error {
	cmd, ok := r.commands[commandName]
	if cmd == nil || !ok {
		return errors.New("command not found")
	}

	if cmd.registered {
		err := r.session.ApplicationCommandDelete(r.session.State.User.ID, r.guildID, cmd.command.ID)
		if err != nil {
			return err
		}
	}

	delete(r.commands, commandName)
	return nil
}

// Registers all commands with Discord
func (r *DefaultCommandRegistry) Register() error {
	failed := []string{}

	for _, command := range r.commands {
		if !command.registered {
			c, err := r.session.ApplicationCommandCreate(r.session.State.User.ID, r.guildID, command.command)
			if err != nil {
				logger.Error("Failed to register command", "name", command.name, "error", err)
				failed = append(failed, command.name)
			} else {
				logger.Info("Command registered", "name", c.Name, "id", c.ID)
				command.registered = true
			}
		}
	}

	if len(failed) == 0 {
		return nil
	}

	joined := strings.Join(failed, ", ")
	return errors.New(fmt.Sprintf("failed to register [%s]", joined))
}

// Unregisters all commands with Discord
func (r *DefaultCommandRegistry) Unregister() error {

	failed := []string{}

	registeredCommands, err := r.session.ApplicationCommands(r.session.State.User.ID, r.guildID)
	if err != nil {
		return err
	}
	for _, v := range registeredCommands {
		err := r.session.ApplicationCommandDelete(r.session.State.User.ID, r.guildID, v.ID)
		if err != nil {
			failed = append(failed, v.Name)
		}
	}

	if len(failed) == 0 {
		return nil
	}

	joined := strings.Join(failed, ", ")
	return errors.New(fmt.Sprintf("failed to unregister [%s]", joined))
}

// Runs the command handler for the associated command
func (r *DefaultCommandRegistry) Run(i *discordgo.InteractionCreate) error {
	cmd, found := r.commands[i.ApplicationCommandData().Name]
	if cmd == nil || !found {
		return errors.New("failed to find command")
	}

	logger.Info("Command triggered", "name", cmd.name, "userID", i.Member.User.ID, "username", i.Member.User.Username)
	go cmd.handler(r.session, i)
	return nil
}
