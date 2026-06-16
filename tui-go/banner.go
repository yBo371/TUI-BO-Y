package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const bannerBOY = `
██████╗  ██████╗      ██╗   ██╗
██╔══██╗██╔═══██╗     ╚██╗ ██╔╝
██████╔╝██║   ██║█████╗╚████╔╝ 
██╔══██╗██║   ██║╚════╝ ╚██╔╝  
██████╔╝╚██████╔╝        ██║   
╚═════╝  ╚═════╝         ╚═╝   
`

var bannerStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFD866"))

func renderBanner(terminalWidth int) string {
	banner := strings.TrimRight(bannerBOY, "\n")

	// ANSI Shadow 很宽。
	// 如果终端太窄，就不要强行显示大字，否则会自动换行变碎。
	bannerWidth := lipgloss.Width(banner)

	// 这里多加 16，是为了给外层 padding、border 留空间。
	if terminalWidth > 0 && terminalWidth < bannerWidth+16 {
		return titleStyle.Render("BO-Y")
	}

	return bannerStyle.Render(banner)
}
