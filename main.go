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
	Token           string
	OwnerId         string
	GuildID         string
	ChannelID       string
	MessageTemplate string
)

type UIHandler interface {
	ShowMessage(m *discordgo.MessageCreate)
	HandleInput(m *MessageUtils)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	loadEnv(&Token, "TOKEN", "t", "Bot Token")
	loadEnv(&OwnerId, "OWNER_ID", "id", "Bot owner's ID")
	loadEnv(&GuildID, "GUILD_ID", "gid", "The ID of the discord server where this bot will act")
	loadEnv(&ChannelID, "CHANNEL_ID", "cid", "The ID of the channel where the messages will be sent")
	loadEnv(&MessageTemplate, "MESSAGE_TEMPLATE", "tmpl", "A template string with a single %s where the message will be inserted")
	flag.Parse()

	if len(Token) == 0 || len(OwnerId) == 0 {
		log.Fatal("ERR: Missing env token or owner id")
	}
}

func loadEnv(target *string, key string, short string, use string) {
	envValue, _ := os.LookupEnv(key)
	flag.StringVar(target, key, envValue, use)
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("ERR: couldn't create discord session. cause:", err)
		return
	}

	// Init UI dependency
	uiHandler := BasicUiHandler{}
	messagUtils := CreateMessageUtils(GuildID, ChannelID, MessageTemplate, dg)

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		messageCreate(uiHandler, s, m)
	})

	dg.Identify.Intents = discordgo.IntentGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("ERR: couldn't open session. cause: ", err)
		return
	}

	go handleInput(uiHandler, dg, &messagUtils)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func handleInput(uiHandler UIHandler, s *discordgo.Session, m *MessageUtils) {
	uiHandler.HandleInput(m)
}

func messageCreate(uiHandler UIHandler, s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if mentionsOwnerOrBot(s.State.User.ID, m) {
		uiHandler.ShowMessage(m)
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
