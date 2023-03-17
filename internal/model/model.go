package model

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type wiki interface {
	GetTablesMatrix(ctx context.Context, page string, lang string, cleanRef bool, tables ...int) ([][][]string, error)
}

const (
	pageIndex           = 0
	langIndex           = 1
	cleanRefIndex       = 2
	maxColumnWidthIndex = 3
)

type input struct {
	inputs         []textinput.Model
	focus          int
	maxColumnWidth int
}

type Model struct {
	wiki     wiki
	input    input
	inputErr error
	tables   []*table
	index    int
	height   int
	width    int
	mode     string
}

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	redStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("001"))
	noStyle      = lipgloss.NewStyle()

	focusedButton = focusedStyle.Copy().Render("Submit")
	blurredButton = blurredStyle.Render("Submit")
)

func NewModel(wiki wiki) *Model {
	var inputs []textinput.Model

	page := textinput.New()
	page.Placeholder = "Arhaan_Khan"
	page.PromptStyle = focusedStyle
	page.TextStyle = focusedStyle
	page.Focus()
	inputs = append(inputs, page)

	lang := textinput.New()
	lang.Placeholder = "en"
	inputs = append(inputs, lang)

	cleanRef := textinput.New()
	cleanRef.Placeholder = "true"
	inputs = append(inputs, cleanRef)

	maxColumnWidth := textinput.New()
	inputs = append(inputs, maxColumnWidth)

	return &Model{
		mode: "input",
		wiki: wiki,
		input: input{
			inputs: inputs,
		},
		index: 0,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case "input":
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.height = msg.Height
			m.width = msg.Width
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "tab", "enter", "down":
				key := msg.String()

				if key == "enter" && m.input.focus == len(m.input.inputs) {
					data, err := m.readInput(context.Background())
					if err != nil {
						m.inputErr = err
						return m, nil
					}
					m.inputErr = nil

					m.setTables(data)

					m.mode = "table"
					m.index = 0
					return m, nil
				}

				m.input.focus++
				if m.input.focus > len(m.input.inputs) {
					m.input.focus = 0
				}

				return m, m.setInputFocus()
			case "up":
				m.input.focus--
				if m.input.focus < 0 {
					m.input.focus = len(m.input.inputs)
				}
				return m, m.setInputFocus()
			default:
				cmds := make([]tea.Cmd, len(m.input.inputs))
				for i := range m.input.inputs {
					m.input.inputs[i], cmds[i] = m.input.inputs[i].Update(msg)
				}
				return m, tea.Batch(cmds...)

			}
		}
	case "table":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "ctrl+n":
				m.mode = "input"
				m.input.focus = 0
				return m, m.setInputFocus()
			case "tab":
				m.index++
				if m.index >= len(m.tables) {
					m.index = 0
				}
			case "shift+tab":
				m.index--
				if m.index < 0 {
					m.index = len(m.tables) - 1
				}
			case "enter", "down", "j":
				m.tables[m.index].moveDown(1)
			case "up", "k":
				m.tables[m.index].moveUp(1)
			case "g":
				m.tables[m.index].goToTop()
			case "G":
				m.tables[m.index].goToBottom()
			case "ctrl+d":
				m.tables[m.index].remove()
			case "ctrl+k":
				m.tables[m.index].switchCursorMode()
			case "ctrl+r":
				m.tables[m.index].reset(m.height - 1)
			case "ctrl+t":
				if m.index == 0 {
					m.tables = m.tables[1:]
					return m, nil
				}
				if m.index == len(m.tables)-1 {
					m.tables = m.tables[:len(m.tables)-1]
					m.index--
					return m, nil
				}
				m.tables = append(m.tables[:m.index], m.tables[m.index+1:]...)
				return m, nil
			}
		case tea.WindowSizeMsg:
			m.height = msg.Height
			m.width = msg.Width
		}
		return m, nil
	}
	return m, nil
}

func (m *Model) View() string {
	switch m.mode {
	case "input":
		return m.ViewInput()
	case "table":
		return m.tables[m.index].model.View()
	default:
		return ""
	}
}

func (m *Model) ViewInput() string {
	var b strings.Builder
	b.WriteString("Comma-separated Wikipedia page titles\n")
	b.WriteString(fmt.Sprintf("%s\n", m.input.inputs[pageIndex].View()))
	b.WriteString("\n")
	b.WriteString("Comma-separated language codes of the pages\n")
	b.WriteString(fmt.Sprintf("%s\n", m.input.inputs[langIndex].View()))
	b.WriteString("\n")
	b.WriteString("Remove the reference link texts (true or false)\n")
	b.WriteString(fmt.Sprintf("%s\n", m.input.inputs[cleanRefIndex].View()))
	b.WriteString("\n")
	b.WriteString("Maximum width of columns (leave empty for no maximum)\n")
	b.WriteString(fmt.Sprintf("%s\n", m.input.inputs[maxColumnWidthIndex].View()))

	button := blurredButton
	if m.input.focus == len(m.input.inputs) {
		button = focusedButton
	}
	b.WriteString(fmt.Sprintf("\n\n%s\n\n", button))

	if m.inputErr != nil {
		b.WriteString(redStyle.Render(m.inputErr.Error()))
	}

	return lipgloss.NewStyle().Width(m.width).Height(m.height).Align(lipgloss.Center, lipgloss.Center).Render(b.String())
}

func (m *Model) setTables(data [][][]string) {
	var tables []*table
	for _, table := range data {
		tables = append(tables, newTable(table, m.height-1, m.input.maxColumnWidth))
	}
	m.tables = tables
}

func (m *Model) setInputFocus() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.input.inputs))
	for i := 0; i < len(m.input.inputs); i++ {
		if i == m.input.focus {
			cmds[i] = m.input.inputs[i].Focus()
			m.input.inputs[i].PromptStyle = focusedStyle
			m.input.inputs[i].TextStyle = focusedStyle
			continue
		}

		m.input.inputs[i].Blur()
		m.input.inputs[i].PromptStyle = noStyle
		m.input.inputs[i].TextStyle = noStyle
	}

	return tea.Batch(cmds...)
}

func (m *Model) readInput(ctx context.Context) ([][][]string, error) {
	var err error

	page := m.input.inputs[pageIndex].Value()
	if page == "" {
		return nil, fmt.Errorf("invalid value: page must be set")
	}
	pages := strings.Split(page, ",")

	lang := m.input.inputs[langIndex].Value()
	if lang == "" {
		return nil, fmt.Errorf("invalid value: language code must be set")
	}
	langs := strings.Split(lang, ",")

	if len(pages) != len(langs) {
		return nil, fmt.Errorf("invalid value: number of pages and languages codes are not equal")
	}

	v := m.input.inputs[cleanRefIndex].Value()
	cleanRef, err := strconv.ParseBool(v)
	if err != nil {
		return nil, fmt.Errorf("invalid value %v: must be true or false", v)
	}

	v = m.input.inputs[maxColumnWidthIndex].Value()
	if v == "" {
		m.input.maxColumnWidth = 0
	} else {
		m.input.maxColumnWidth, err = strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value %v: must be a valid number", v)
		}
		if m.input.maxColumnWidth <= 0 {
			m.input.maxColumnWidth = 0
		}
	}

	var tables [][][]string
	for i := range pages {
		data, err := m.wiki.GetTablesMatrix(ctx, pages[i], langs[i], cleanRef)
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			return nil, fmt.Errorf("no tables on page %s", page)
		}

		for _, table := range data {
			fillRowData(table)
			tables = append(tables, table)
		}
	}

	return tables, nil
}

func fillRowData(data [][]string) {
	colLen := len(data[0])
	for i := 1; i < len(data); i++ {
		if rowLen := len(data[i]); rowLen < colLen {
			for j := 0; j < colLen-rowLen; j++ {
				data[i] = append(data[i], "")
			}
		}
	}
}
