package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"jtams/playertrackerbot/commands"
	"jtams/playertrackerbot/tracker"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Check for required environment variables
	// BOT_TOKEN can be found at https://discord.com/developers/applications at it can be set in the .env file or in the environment
	// GUILD_ID can be found by right-clicking on the server icon in Discord and selecting "Copy ID"
	// SAVE_FILE is the file where the tracker data will be saved (optional, defaults to data/save_data.json)
	// LOG_LEVEL is the level of logging to use (optional, defaults to errors only)
	if os.Getenv("BOT_TOKEN") == "" || os.Getenv("GUILD_ID") == "" {
		panic("BOT_TOKEN and GUILD_ID required in environment. Please set them in the .env file or in the environment.")
	}

	// Optional save file
	if os.Getenv("SAVE_FILE") == "" {
		os.Setenv("SAVE_FILE", "save_data.json")
	}

	// discordgo.Logger
	if os.Getenv("LOG_LEVEL") == "" {
		os.Setenv("LOG_LEVEL", "error")
	}

	// Create Bot
	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	// Logging
	getDiscordGoLogLevel(discord, os.Getenv("LOG_LEVEL"))
	logLevel := getSlogLevel(os.Getenv("LOG_LEVEL"))
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	loggerHandler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(loggerHandler)
	slog.SetDefault(logger)
	setDiscordGoLogger(logger)

	logger.Info("Log Level", "level", logLevel.String())

	// Handles log in event
	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Info("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Bot Login
	err = discord.Open()
	if err != nil {
		panic("Cannot open the session")
	}
	defer discord.Close()

	// Create Command Registry for handling commands
	commandRegistry := commands.NewDefaultCommandRegistry(discord, os.Getenv("GUILD_ID"))

	// Initalize trakcer and messenger
	var playerTracker *tracker.PlayerTracker
	var messageUpdater *tracker.Messenger

	// Attempt to load from save file
	err, messageUpdater, playerTracker = tracker.LoadTrackerData(os.Getenv("SAVE_FILE"), discord)
	if messageUpdater == nil {
		messageUpdater = tracker.NewMessageUpdater(discord)
	}
	if playerTracker == nil {
		playerTracker = tracker.NewPlayerTracker()
	}

	// If the tracker crashed or was stopped while running, resume.
	if playerTracker.Running {
		logger.Info("Forcing startup")
		commands.ForceStartup(messageUpdater, playerTracker)
	}

	// Register all default commands
	commands.RegisterAllDefaultCommands(discord, commandRegistry, playerTracker, messageUpdater)

	// Listens for interactions and runs the appropriate command handler
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if err := commandRegistry.Run(i); err != nil {
			logger.Error(err.Error())
		}
	})

	// Handles messages sent by users.
	// Specifically, it counts every message sent so that the messenger can resend the tracker
	// if theres been too many messages sent since the tracker message was sent. This ensures
	// the tracker is always visible in the channel.
	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Check if the message is from the bot
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.ChannelID == messageUpdater.Message.ChannelID {
			messageUpdater.MessageOverflow++
		}
	})

	// Wait here until CTRL-C or other term signal is received.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	logger.Info("Press Ctrl+C to exit")
	<-stop
	logger.Info("Shutting Down...")

	// When bot stops, remove registered commands
	if err = commandRegistry.Unregister(); err != nil {
		logger.Error(err.Error())
	} else {
		logger.Info("All commands unregistered")
	}

	logger.Info("Goodbye!")
}

func getSlogLevel(logLevel string) slog.Level {
	logLevel = strings.ToLower(logLevel)

	switch logLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelError
	}
}

func setDiscordGoLogger(logger *slog.Logger) {
	// var Logger func(msgL, caller int, format string, a ...interface{})
	discordgo.Logger = func(msgL, caller int, format string, a ...interface{}) {
		msg := fmt.Sprintf(format, a...)

		if strings.Contains(msg, "RawData:json.RawMessag") {
			return
		}

		if strings.Contains(msg, "Heartbeat") {
			return
		}

		switch msgL {
		case discordgo.LogDebug:
			logger.Debug(msg, "type", "discordgo")
			break
		case discordgo.LogInformational:
			logger.Info(msg, "type", "discordgo")
			break
		case discordgo.LogWarning:
			logger.Warn(msg, "type", "discordgo")
			break
		case discordgo.LogError:
			logger.Error(msg, "type", "discordgo")
			break
		default:
			logger.Info(msg, "type", "discordgo")
		}
	}
}

// Sets the logging level for the discord session
func getDiscordGoLogLevel(discord *discordgo.Session, logLevel string) {
	logLevel = strings.ToLower(logLevel)

	switch logLevel {
	case "debug":
		discord.LogLevel = discordgo.LogDebug
	case "info":
		discord.LogLevel = discordgo.LogInformational
	case "warn":
		discord.LogLevel = discordgo.LogWarning
	case "error":
		discord.LogLevel = discordgo.LogError
	default:
		discord.LogLevel = discordgo.LogWarning
	}
}
