package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type page int

const (
	pageHome page = iota
	pagePassword
)

type model struct {
	choices  []string
	cursor   int
	message  string
	page     page
	width    int
	height   int
	password passwordModel
}

func initialModel() model {
	return model{
		choices: []string{
			"密码生成器",
			"退出",
		},
		message:  "欢迎使用 my-tui",
		page:     pageHome,
		width:    80,
		height:   24,
		password: newPasswordModel(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		key := msg.String()

		if key == "ctrl+c" || key == "q" {
			return m, tea.Quit
		}

		if m.page == pagePassword {
			if key == "esc" {
				m.page = pageHome
				m.message = "已返回主面板"
				return m, nil
			}

			m.password = m.password.update(key)
			return m, nil
		}

		if m.page == pageHome {
			switch key {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}

			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}

			case "enter":
				switch m.cursor {
				case 0:
					m.page = pagePassword
					m.password = m.password.refresh("已进入密码生成器")

				case 1:
					return m, tea.Quit
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.page {
	case pagePassword:
		return appStyle.Render(m.password.view(m.width))
	default:
		return appStyle.Render(m.homeView())
	}
}

func (m model) homeView() string {
	var lines []string

	lines = append(lines, titleStyle.Render("MY-TUI"))
	lines = append(lines, subtitleStyle.Render("一个用 Go 写的终端工具箱"))
	lines = append(lines, "")

	for i, choice := range m.choices {
		cursor := "  "
		text := normalStyle.Render(choice)

		if m.cursor == i {
			cursor = "❯ "
			text = selectedStyle.Render(choice)
		}

		lines = append(lines, fmt.Sprintf("%s%s", cursor, text))
	}

	lines = append(lines, "")
	lines = append(lines, messageStyle.Render("状态："+m.message))
	lines = append(lines, "")
	lines = append(lines, helpStyle.Render("↑/↓ 或 k/j 选择 · Enter 确认 · q 退出"))

	content := strings.Join(lines, "\n")

	width := 70
	if m.width > 90 {
		width = 78
	}

	return panelStyle.Width(width).Render(content)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "程序出错：%v\n", err)
		os.Exit(1)
	}
}
