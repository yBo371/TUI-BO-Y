package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type page int

const (
	pageHome page = iota
	pageLogs
	pagePassword
)

type model struct {
	choices        []string
	cursor         int
	message        string
	page           page
	width          int
	height         int
	password       string
	passwordLength int
	useSymbols     bool
}

var (
	appStyle = lipgloss.NewStyle().
			Padding(1, 2)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7DCFFF")).
			Padding(1, 3)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFD866"))

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A6ACCD"))

	cardTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7DCFFF")).
			Bold(true)

	cardValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C3E88D")).
			Bold(true)

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#444A73")).
			Padding(0, 2).
			Width(16).
			Height(4)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD866")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D6DEEB"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C7086"))

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C3E88D")).
			Bold(true)

	warnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F78C6C")).
			Bold(true)

	passwordStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C3E88D")).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#C3E88D")).
			Padding(1, 2).
			Width(58)
)

func initialModel() model {
	return model{
		choices: []string{
			"查看状态",
			"密码生成器",
			"启动服务",
			"停止服务",
			"重启服务",
			"查看日志",
			"修改配置",
			"退出",
		},
		message:        "欢迎使用 my-tui",
		page:           pageHome,
		width:          80,
		height:         24,
		passwordLength: 20,
		useSymbols:     true,
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
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc", "b":
			if m.page != pageHome {
				m.page = pageHome
				m.message = "已返回首页"
				return m, nil
			}
		}

		if m.page == pagePassword {
			switch msg.String() {
			case "g", "enter":
				password, err := generatePassword(m.passwordLength, m.useSymbols)
				if err != nil {
					m.message = "密码生成失败：" + err.Error()
				} else {
					m.password = password
					m.message = "已生成新密码"
				}

			case "+", "=":
				if m.passwordLength < 64 {
					m.passwordLength++
				}
				password, err := generatePassword(m.passwordLength, m.useSymbols)
				if err == nil {
					m.password = password
				}

			case "-", "_":
				if m.passwordLength > 8 {
					m.passwordLength--
				}
				password, err := generatePassword(m.passwordLength, m.useSymbols)
				if err == nil {
					m.password = password
				}

			case "s":
				m.useSymbols = !m.useSymbols
				password, err := generatePassword(m.passwordLength, m.useSymbols)
				if err == nil {
					m.password = password
				}
				if m.useSymbols {
					m.message = "已开启特殊符号"
				} else {
					m.message = "已关闭特殊符号"
				}
			}

			return m, nil
		}

		if m.page == pageHome {
			switch msg.String() {
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
					m.message = "当前状态：服务运行中"

				case 1:
					password, err := generatePassword(m.passwordLength, m.useSymbols)
					if err != nil {
						m.message = "密码生成失败：" + err.Error()
					} else {
						m.password = password
						m.message = "已进入密码生成器"
						m.page = pagePassword
					}

				case 2:
					m.message = "正在启动服务..."

				case 3:
					m.message = "正在停止服务..."

				case 4:
					m.message = "正在重启服务..."

				case 5:
					m.page = pageLogs

				case 6:
					m.message = "配置功能下一步再做"

				case 7:
					return m, tea.Quit
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.page {
	case pageLogs:
		return appStyle.Render(m.logsView())
	case pagePassword:
		return appStyle.Render(m.passwordView())
	default:
		return appStyle.Render(m.homeView())
	}
}

func (m model) homeView() string {
	var lines []string

	lines = append(lines, titleStyle.Render("MY-TUI"))
	lines = append(lines, subtitleStyle.Render("一个用 Go 写的终端交互式管理面板"))
	lines = append(lines, "")

	cards := lipgloss.JoinHorizontal(
		lipgloss.Top,
		renderCard("服务状态", "Running"),
		"  ",
		renderCard("版本", "v0.2.0"),
		"  ",
		renderCard("系统", runtime.GOOS),
	)

	lines = append(lines, cards)
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

func (m model) logsView() string {
	logs := []string{
		"[09:45:01] my-tui started",
		"[09:45:02] loading config...",
		"[09:45:03] service status: running",
		"[09:45:04] checking network...",
		"[09:45:05] everything looks good",
		"[09:45:06] waiting for command...",
	}

	var lines []string

	lines = append(lines, titleStyle.Render("实时日志"))
	lines = append(lines, subtitleStyle.Render("这里以后可以接真实日志文件，比如 app.log"))
	lines = append(lines, "")

	for _, line := range logs {
		lines = append(lines, normalStyle.Render(line))
	}

	lines = append(lines, "")
	lines = append(lines, warnStyle.Render("这是模拟日志，下一步再接真实命令和日志文件"))
	lines = append(lines, "")
	lines = append(lines, helpStyle.Render("Esc / b 返回首页 · q 退出"))

	content := strings.Join(lines, "\n")

	width := 70
	if m.width > 90 {
		width = 78
	}

	return panelStyle.Width(width).Render(content)
}

func (m model) passwordView() string {
	var lines []string

	symbolStatus := "开启"
	if !m.useSymbols {
		symbolStatus = "关闭"
	}

	lines = append(lines, titleStyle.Render("密码生成器"))
	lines = append(lines, subtitleStyle.Render("使用 crypto/rand 生成随机强密码"))
	lines = append(lines, "")

	lines = append(lines, normalStyle.Render(fmt.Sprintf("长度：%d", m.passwordLength)))
	lines = append(lines, normalStyle.Render("特殊符号："+symbolStatus))
	lines = append(lines, "")

	if m.password == "" {
		lines = append(lines, warnStyle.Render("按 g 生成密码"))
	} else {
		lines = append(lines, passwordStyle.Render(m.password))
	}

	lines = append(lines, "")
	lines = append(lines, messageStyle.Render("状态："+m.message))
	lines = append(lines, "")
	lines = append(lines, helpStyle.Render("g / Enter 重新生成 · + 增加长度 · - 减少长度 · s 开关符号 · Esc / b 返回"))

	content := strings.Join(lines, "\n")

	width := 70
	if m.width > 90 {
		width = 78
	}

	return panelStyle.Width(width).Render(content)
}

func renderCard(title string, value string) string {
	content := cardTitleStyle.Render(title) + "\n" + cardValueStyle.Render(value)
	return cardStyle.Render(content)
}

func generatePassword(length int, useSymbols bool) (string, error) {
	if length < 8 {
		length = 8
	}

	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	symbols := "!@#$%^&*()-_=+[]{};:,.?/"

	sets := []string{lower, upper, digits}

	if useSymbols {
		sets = append(sets, symbols)
	}

	allChars := strings.Join(sets, "")

	password := make([]byte, 0, length)

	for _, set := range sets {
		ch, err := randomChar(set)
		if err != nil {
			return "", err
		}
		password = append(password, ch)
	}

	for len(password) < length {
		ch, err := randomChar(allChars)
		if err != nil {
			return "", err
		}
		password = append(password, ch)
	}

	for i := len(password) - 1; i > 0; i-- {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}

		j := int(n.Int64())
		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}

func randomChar(chars string) (byte, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
	if err != nil {
		return 0, err
	}

	return chars[n.Int64()], nil
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "程序出错：%v\n", err)
		os.Exit(1)
	}
}
