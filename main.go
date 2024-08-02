package main

import (
	"log"
	"os"
	"os/signal"

	"jtams/playertrackerbot/bot"
	"jtams/playertrackerbot/commands"
	"jtams/playertrackerbot/tracker"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

func addHandlers(discord *discordgo.Session, registry bot.CommandRegistry, tracker *tracker.PlayerTracker, messageUpdater *tracker.Messenger) {
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if err := registry.Run(i); err != nil {
			log.Println(err)
		}
	})

	if err := registry.AddCommand(*commands.PingCommand(), commands.PingHandler()); err != nil {
		log.Println(err)
	}

	if err := registry.AddCommand(*commands.StartCommand(), commands.StartHandler(messageUpdater, tracker)); err != nil {
		log.Println(err)
	}

	if err := registry.AddCommand(*commands.UserCommand(), commands.UserHandler(messageUpdater, tracker)); err != nil {
		log.Println(err)
	}

	err := registry.Register()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	// Check for required environment variables
	// BOT_TOKEN can be found at https://discord.com/developers/applications at it can be set in the .env file or in the environment
	// GUILD_ID can be found by right-clicking on the server icon in Discord and selecting "Copy ID"
	if os.Getenv("BOT_TOKEN") == "" || os.Getenv("GUILD_ID") == "" {
		panic("BOT_TOKEN and GUILD_ID required in environment. Please set them in the .env file or in the environment.")
	}

	// Optional save file
	if os.Getenv("SAVE_FILE") == "" {
		os.Setenv("SAVE_FILE", "saves/save_data.json")
	}

	// Create Bot
	discord, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}
	discord.LogLevel = discordgo.LogDebug
	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Bot Login
	err = discord.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer discord.Close()

	// Create Command Registry for handling commands
	commandRegistry := bot.NewDefaultCommandRegistry(discord, os.Getenv("GUILD_ID"))

	var playerTracker *tracker.PlayerTracker
	var messageUpdater *tracker.Messenger

	err, messageUpdater, playerTracker = tracker.LoadTrackerData(os.Getenv("SAVE_FILE"), discord)
	if messageUpdater == nil {
		messageUpdater = tracker.NewMessageUpdater(discord)
	}
	if playerTracker == nil {
		playerTracker = tracker.NewPlayerTracker()
	}

	if playerTracker.Running {
		log.Println("Forcing startup")
		commands.ForceStartup(messageUpdater, playerTracker)
	}

	addHandlers(discord, commandRegistry, playerTracker, messageUpdater)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if err = commandRegistry.Unregister(); err != nil {
		log.Println(err)
	}

	log.Println("Shutting Down")
}
