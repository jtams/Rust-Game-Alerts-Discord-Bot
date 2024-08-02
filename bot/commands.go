package bot

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler func(*discordgo.Session, *discordgo.InteractionCreate) error

type CommandRegistry interface {
	// Adds command to registry
	AddCommand(discordgo.ApplicationCommand, CommandHandler) error

	// Removes command from registry
	RemoveCommand(commandName string) error

	// Registers command with Discord
	Register() error

	// Unregisters command with Discord
	Unregister() error

	// Run a commmands handler
	Run(*discordgo.InteractionCreate) error
}

type CommandData struct {
	name       string
	command    *discordgo.ApplicationCommand
	handler    CommandHandler
	registered bool
}

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

func (r *DefaultCommandRegistry) Register() error {
	failed := []string{}

	for _, command := range r.commands {
		if !command.registered {
			c, err := r.session.ApplicationCommandCreate(r.session.State.User.ID, r.guildID, command.command)
			if err != nil {
				log.Print(err)
				failed = append(failed, command.name)
			} else {
				log.Print(c.Name, " registered. ID: ", c.ID)
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

func (r *DefaultCommandRegistry) Unregister() error {

	failed := []string{}

	for _, command := range r.commands {
		if command.registered {
			err := r.session.ApplicationCommandDelete(r.session.State.User.ID, r.guildID, command.command.ID)
			if err != nil {
				failed = append(failed, command.name)
			} else {
				command.registered = false
			}
		}
	}

	if len(failed) == 0 {
		return nil
	}

	joined := strings.Join(failed, ", ")
	return errors.New(fmt.Sprintf("failed to unregister [%s]", joined))
}

func (r *DefaultCommandRegistry) Run(i *discordgo.InteractionCreate) error {
	cmd, found := r.commands[i.ApplicationCommandData().Name]
	if cmd == nil || !found {
		return errors.New("failed to find command")
	}

	go cmd.handler(r.session, i)
	return nil
}
