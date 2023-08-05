package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	Token   string
	OwnerId string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	envToken, _ := os.LookupEnv("TOKEN")
	flag.StringVar(&Token, "t", envToken, "Bot Token")
	envOwnerId, _ := os.LookupEnv("OWNER_ID")
	flag.StringVar(&OwnerId, "id", envOwnerId, "Bot owner's ID")
	flag.Parse()

	if len(Token) == 0 || len(OwnerId) == 0 {
		log.Fatal("ERR: Missing env token or owner id")
	}
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if mentionsOwnerOrBot(s.State.User.ID, m) {
		log.Println(m.Content)
	}
}

func mentionsOwnerOrBot(botId string, m *discordgo.MessageCreate) bool {
	for _, mentionedUser := range m.Mentions {
		if mentionedUser.ID == OwnerId || mentionedUser.ID == botId {
			return true
		}
	}
	return false
}
