package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

// Provide the bot with a token from the command line
// Example: go run main.go -t <bot_token>
func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// Create a counter to keep track of the number of "gm" messages sent by a specific user.
// This is a variable that will be incremented every time a specific user sends a "gm" message.
// We will also use this to keep track of the number of "gm" messages sent by a specific user in a row.
var gmCounter = make(map[string]int)

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is saying "gm" reply with "Good Morning!"
	if m.Content == "gm" {
		// If the user has said gm before 24 hours have passed then don't increment their counter.
		if m.Author.ID != "" && gmCounter[m.Author.ID] < 1 {
			gmCounter[m.Author.ID]++
			// Reply with the number of "gm" messages the user has sent
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Good Morning! 🌞 You have sent %d Good Morning messages.", gmCounter[m.Author.ID]))
		}
	}
}
