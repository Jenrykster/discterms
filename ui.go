package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	ResetColor = "\033[0m"
	BlueColor  = "\033[34m"
)

type BasicUiHandler struct{}

func showPrompt() {
	fmt.Print("> ")
}

func (BasicUiHandler) HandleInput(m *MessageUtils) {
	fmt.Println("You can send messages by pressing ENTER")
	showPrompt()
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		showPrompt()
		text := scanner.Text()
		m.SendMessage(text)
	}
}

func (BasicUiHandler) ShowMessage(m *discordgo.MessageCreate) {
	clearLastLine()
	log.Println(cleanMessage(m))
	showPrompt()
}

func clearLastLine() {
	fmt.Print("\r")     // Move the cursor to the beginning of the line
	fmt.Print("\033[K") // Clear the line from the cursor position to the end
}

func cleanMessage(m *discordgo.MessageCreate) string {
	messageResult := m.Content
	for _, mention := range m.Mentions {
		mentionId := fmt.Sprintf("<@%s>", mention.ID)
		cleanMention := fmt.Sprintf("%s@%s%s", BlueColor, mention.Username, ResetColor)
		messageResult = strings.ReplaceAll(messageResult, mentionId, cleanMention)
	}
	return messageResult
}
