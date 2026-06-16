package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/exec"
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
	useUpper       bool
	useLower       bool
	useDigits      bool
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

	onStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C3E88D")).
		Bold(true)

	offStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F78C6C")).
			Bold(true)
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
		useUpper:       true,
		useLower:       true,
		useDigits:      true,
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
		key := msg.String()

		if key == "ctrl+c" || key == "q" {
			return m, tea.Quit
		}

		if m.page == pagePassword {
			return updatePasswordPage(m, key), nil
		}

		if key == "esc" || key == "b" {
			if m.page != pageHome {
				m.page = pageHome
				m.message = "已返回首页"
				return m, nil
			}
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
					m.message = "当前状态：服务运行中"

				case 1:
					m.page = pagePassword
					m = refreshPassword(m, "已进入密码生成器")

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

func updatePasswordPage(m model, key string) model {
	switch key {
	case "esc":
		m.page = pageHome
		m.message = "已返回首页"

	case "enter":
		m = refreshPassword(m, "已重新生成密码")

	case "+", "=":
		if m.passwordLength < 64 {
			m.passwordLength++
		}
		m = refreshPassword(m, fmt.Sprintf("密码长度：%d", m.passwordLength))

	case "-", "_":
		if m.passwordLength > 8 {
			m.passwordLength--
		}
		m = refreshPassword(m, fmt.Sprintf("密码长度：%d", m.passwordLength))

	case "1":
		if m.useUpper && m.enabledCharsetCount() == 1 {
			m.message = "至少保留一种字符类型"
			return m
		}
		m.useUpper = !m.useUpper
		m = refreshPassword(m, "已切换大写英文字母")

	case "2":
		if m.useLower && m.enabledCharsetCount() == 1 {
			m.message = "至少保留一种字符类型"
			return m
		}
		m.useLower = !m.useLower
		m = refreshPassword(m, "已切换小写英文字母")

	case "3":
		if m.useDigits && m.enabledCharsetCount() == 1 {
			m.message = "至少保留一种字符类型"
			return m
		}
		m.useDigits = !m.useDigits
		m = refreshPassword(m, "已切换数字")

	case "4":
		if m.useSymbols && m.enabledCharsetCount() == 1 {
			m.message = "至少保留一种字符类型"
			return m
		}
		m.useSymbols = !m.useSymbols
		m = refreshPassword(m, "已切换特殊符号")

	case "5":
		if m.password == "" {
			m.message = "当前还没有密码，先按 Enter 生成"
			return m
		}

		if err := copyToClipboard(m.password); err != nil {
			m.message = "复制失败：" + err.Error()
		} else {
			m.message = "密码已复制到剪贴板"
		}
	}

	return m
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
		renderCard("版本", "v0.3.0"),
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

	lines = append(lines, titleStyle.Render("密码生成器"))
	lines = append(lines, subtitleStyle.Render("使用 crypto/rand 生成随机强密码"))
	lines = append(lines, "")

	options := []string{
		fmt.Sprintf("长度：%d", m.passwordLength),
		"1 大写英文字母：" + renderSwitch(m.useUpper),
		"2 小写英文字母：" + renderSwitch(m.useLower),
		"3 数字：" + renderSwitch(m.useDigits),
		"4 特殊符号：" + renderSwitch(m.useSymbols),
	}

	for _, option := range options {
		lines = append(lines, normalStyle.Render(option))
	}

	lines = append(lines, "")

	if m.password == "" {
		lines = append(lines, warnStyle.Render("按 Enter 生成密码"))
	} else {
		lines = append(lines, passwordStyle.Render(m.password))
	}

	lines = append(lines, "")
	lines = append(lines, messageStyle.Render("状态："+m.message))
	lines = append(lines, "")
	lines = append(lines, helpStyle.Render("Enter 生成 · + 增加长度 · - 减少长度 · 1/2/3/4 开关类型 · 5 复制 · Esc 返回"))

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

func renderSwitch(enabled bool) string {
	if enabled {
		return onStyle.Render("开启")
	}
	return offStyle.Render("关闭")
}

func refreshPassword(m model, msg string) model {
	password, err := generatePassword(
		m.passwordLength,
		m.useUpper,
		m.useLower,
		m.useDigits,
		m.useSymbols,
	)

	if err != nil {
		m.message = "密码生成失败：" + err.Error()
		return m
	}

	m.password = password
	m.message = msg
	return m
}

func (m model) enabledCharsetCount() int {
	count := 0

	if m.useUpper {
		count++
	}
	if m.useLower {
		count++
	}
	if m.useDigits {
		count++
	}
	if m.useSymbols {
		count++
	}

	return count
}

func generatePassword(length int, useUpper bool, useLower bool, useDigits bool, useSymbols bool) (string, error) {
	if length < 8 {
		length = 8
	}

	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lower := "abcdefghijklmnopqrstuvwxyz"
	digits := "0123456789"
	symbols := "!@#$%^&*()-_=+[]{};:,.?/"

	var sets []string

	if useUpper {
		sets = append(sets, upper)
	}
	if useLower {
		sets = append(sets, lower)
	}
	if useDigits {
		sets = append(sets, digits)
	}
	if useSymbols {
		sets = append(sets, symbols)
	}

	if len(sets) == 0 {
		return "", fmt.Errorf("至少需要开启一种字符类型")
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

func copyToClipboard(text string) error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("powershell", "-NoProfile", "-Command", "Set-Clipboard -Value ([Console]::In.ReadToEnd())")
		cmd.Stdin = strings.NewReader(text)
		return cmd.Run()

	case "darwin":
		cmd := exec.Command("pbcopy")
		cmd.Stdin = strings.NewReader(text)
		return cmd.Run()

	case "linux":
		if _, err := exec.LookPath("wl-copy"); err == nil {
			cmd := exec.Command("wl-copy")
			cmd.Stdin = strings.NewReader(text)
			return cmd.Run()
		}

		if _, err := exec.LookPath("xclip"); err == nil {
			cmd := exec.Command("xclip", "-selection", "clipboard")
			cmd.Stdin = strings.NewReader(text)
			return cmd.Run()
		}

		if _, err := exec.LookPath("xsel"); err == nil {
			cmd := exec.Command("xsel", "--clipboard", "--input")
			cmd.Stdin = strings.NewReader(text)
			return cmd.Run()
		}

		return fmt.Errorf("Linux 需要安装 wl-copy、xclip 或 xsel")

	default:
		return fmt.Errorf("暂不支持当前系统：%s", runtime.GOOS)
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "程序出错：%v\n", err)
		os.Exit(1)
	}
}
