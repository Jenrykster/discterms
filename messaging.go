package main

import (
	"fmt"
	"log"
	"regexp"

	"strings"

	"github.com/bwmarrin/discordgo"
)

type MessageUtils struct {
	messageTemplate  string
	members          []*discordgo.Member
	GuildID          string
	DefaultChannelID string
	Session          *discordgo.Session
}

func CreateMessageUtils(guildId string, channelId string, messageTemplate string, session *discordgo.Session) MessageUtils {
	members, err := session.GuildMembers(guildId, "", 20)
	if err != nil {
		fmt.Println("ERR: couldn't get users. cause: ", err)
	}
	return MessageUtils{
		members:          members,
		GuildID:          guildId,
		DefaultChannelID: channelId,
		Session:          session,
		messageTemplate:  messageTemplate,
	}
}

func (m *MessageUtils) getUserWithUsername(username string) *discordgo.Member {
	cleanUsername := strings.Replace(username, "@", "", 1)

	for _, member := range m.members {
		if member.User.Username == cleanUsername {
			return member
		}
	}
	return nil
}

func (m *MessageUtils) replaceUsernamesWithId(content string) string {
	currContent := content
	regexPattern := `@\w+`
	reg := regexp.MustCompile(regexPattern)
	matches := reg.FindAllString(content, -1)

	for _, match := range matches {
		usernameOwner := m.getUserWithUsername(match)
		if usernameOwner != nil {
			mention := fmt.Sprintf("<@%s>", usernameOwner.User.ID)
			currContent = strings.ReplaceAll(currContent, match, mention)
		}
	}

	return currContent
}

func (m *MessageUtils) SendMessage(content string) {
	content = m.replaceUsernamesWithId(content)
	_, err := m.Session.ChannelMessageSend(m.DefaultChannelID, fmt.Sprintf(m.messageTemplate, content))
	if err != nil {
		log.Println("ERR: ", err)
	}
}
