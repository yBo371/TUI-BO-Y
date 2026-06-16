package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices []string
	cursor  int
	message string
}

func initialModel() model {
	return model{
		choices: []string{
			"查看状态",
			"启动服务",
			"停止服务",
			"重启服务",
			"查看日志",
			"退出",
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

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
				m.message = "当前状态：运行中"
			case 1:
				m.message = "正在启动服务..."
			case 2:
				m.message = "正在停止服务..."
			case 3:
				m.message = "正在重启服务..."
			case 4:
				m.message = "这里以后显示日志"
			case 5:
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "\n  my-tui 管理面板\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("  %s %s\n", cursor, choice)
	}

	s += "\n  ↑/↓ 或 k/j 选择，Enter 确认，q 退出\n"

	if m.message != "" {
		s += "\n  " + m.message + "\n"
	}

	return s
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "程序出错：%v\n", err)
		os.Exit(1)
	}
}
