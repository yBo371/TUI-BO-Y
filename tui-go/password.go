package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os/exec"
	"runtime"
	"strings"
)

type passwordModel struct {
	password string
	length   int
	message  string

	useUpper   bool
	useLower   bool
	useDigits  bool
	useSymbols bool
}

func newPasswordModel() passwordModel {
	return passwordModel{
		length:     20,
		message:    "按 Enter 生成密码",
		useUpper:   true,
		useLower:   true,
		useDigits:  true,
		useSymbols: true,
	}
}

func (p passwordModel) update(key string) passwordModel {
	switch key {
	case "enter":
		p = p.refresh("已重新生成密码")

	case "+", "=":
		if p.length < 64 {
			p.length++
		}
		p = p.refresh(fmt.Sprintf("密码长度：%d", p.length))

	case "-", "_":
		if p.length > 8 {
			p.length--
		}
		p = p.refresh(fmt.Sprintf("密码长度：%d", p.length))

	case "1":
		if p.useUpper && p.enabledCharsetCount() == 1 {
			p.message = "至少保留一种字符类型"
			return p
		}
		p.useUpper = !p.useUpper
		p = p.refresh("已切换大写英文字母")

	case "2":
		if p.useLower && p.enabledCharsetCount() == 1 {
			p.message = "至少保留一种字符类型"
			return p
		}
		p.useLower = !p.useLower
		p = p.refresh("已切换小写英文字母")

	case "3":
		if p.useDigits && p.enabledCharsetCount() == 1 {
			p.message = "至少保留一种字符类型"
			return p
		}
		p.useDigits = !p.useDigits
		p = p.refresh("已切换数字")

	case "4":
		if p.useSymbols && p.enabledCharsetCount() == 1 {
			p.message = "至少保留一种字符类型"
			return p
		}
		p.useSymbols = !p.useSymbols
		p = p.refresh("已切换特殊符号")

	case "5":
		if p.password == "" {
			p.message = "当前还没有密码，先按 Enter 生成"
			return p
		}

		if err := copyToClipboard(p.password); err != nil {
			p.message = "复制失败：" + err.Error()
		} else {
			p.message = "密码已复制到剪贴板"
		}
	}

	return p
}

func (p passwordModel) view(width int) string {
	var lines []string

	lines = append(lines, titleStyle.Render("密码生成器"))
	lines = append(lines, subtitleStyle.Render("使用 crypto/rand 生成随机强密码"))
	lines = append(lines, "")

	options := []string{
		fmt.Sprintf("长度：%d", p.length),
		"1 大写英文字母：" + renderSwitch(p.useUpper),
		"2 小写英文字母：" + renderSwitch(p.useLower),
		"3 数字：" + renderSwitch(p.useDigits),
		"4 特殊符号：" + renderSwitch(p.useSymbols),
	}

	for _, option := range options {
		lines = append(lines, normalStyle.Render(option))
	}

	lines = append(lines, "")

	if p.password == "" {
		lines = append(lines, warnStyle.Render("按 Enter 生成密码"))
	} else {
		lines = append(lines, passwordStyle.Render(p.password))
	}

	lines = append(lines, "")
	lines = append(lines, messageStyle.Render("状态："+p.message))
	lines = append(lines, "")
	lines = append(lines, helpStyle.Render("Enter 生成 · + 增加长度 · - 减少长度 · 1/2/3/4 开关类型 · 5 复制 · Esc 返回"))

	content := strings.Join(lines, "\n")

	panelWidth := 70
	if width > 90 {
		panelWidth = 78
	}

	return panelStyle.Width(panelWidth).Render(content)
}

func (p passwordModel) refresh(msg string) passwordModel {
	password, err := generatePassword(
		p.length,
		p.useUpper,
		p.useLower,
		p.useDigits,
		p.useSymbols,
	)

	if err != nil {
		p.message = "密码生成失败：" + err.Error()
		return p
	}

	p.password = password
	p.message = msg
	return p
}

func (p passwordModel) enabledCharsetCount() int {
	count := 0

	if p.useUpper {
		count++
	}
	if p.useLower {
		count++
	}
	if p.useDigits {
		count++
	}
	if p.useSymbols {
		count++
	}

	return count
}

func renderSwitch(enabled bool) string {
	if enabled {
		return onStyle.Render("开启")
	}

	return offStyle.Render("关闭")
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
