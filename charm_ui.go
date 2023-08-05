package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CharmUiHandler struct {
	teaProgram *tea.Program
}

func CreateCharmUiHandler() CharmUiHandler {
	p := tea.NewProgram(initialModel())

	go p.Run()

	return CharmUiHandler{
		teaProgram: p,
	}
}

func (c *CharmUiHandler) HandleInput(m *MessageUtils) {
	// Input is handled by the bubbleTea loop
	c.teaProgram.Send(m)
}

func (c *CharmUiHandler) ShowMessage(m *discordgo.MessageCreate) {
	c.teaProgram.Send(charmCleanMessage(m))
}

func sendMessage(m *MessageUtils, content string) tea.Cmd {
	return func() tea.Msg {
		m.SendMessage(content)
		return normalMsg(content)
	}
}
func charmCleanMessage(m *discordgo.MessageCreate) string {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	messageResult := m.Content
	for _, mention := range m.Mentions {
		mentionId := fmt.Sprintf("<@%s>", mention.ID)
		cleanMention := fmt.Sprintf("%s@%s%s", BlueColor, mention.Username, ResetColor)
		messageResult = strings.ReplaceAll(messageResult, mentionId, cleanMention)
	}

	author := style.Render(fmt.Sprintf("@%s", m.Author.Username))
	return fmt.Sprintf("%s: %s ", author, messageResult)
}

type (
	errMsg    error
	normalMsg string
)

type model struct {
	viewport     viewport.Model
	messages     []string
	textarea     textarea.Model
	senderStyle  lipgloss.Style
	messageUtils *MessageUtils
	err          error
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message"
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(60)
	ta.SetHeight(3)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(100, 5)
	vp.SetContent(`Welcome to discterms!
	Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea: ta,
		messages: []string{},
		viewport: vp,
		senderStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9")).
			Padding(0, 1),
		err: nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case *MessageUtils:
		newModel := m
		newModel.messageUtils = msg
		return newModel, nil
	case normalMsg:
		m.messages = append(m.messages, m.senderStyle.Render("You: ")+string(msg))
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.textarea.Reset()
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			os.Exit(1)
			return m, tea.Quit
		case tea.KeyEnter:
			return m, sendMessage(m.messageUtils, m.textarea.Value())
		}

	case errMsg:
		m.err = msg
		return m, nil
	case string:
		m.messages = append(m.messages, msg)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.textarea.Reset()
		m.viewport.GotoBottom()
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}
