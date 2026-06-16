package main

import "github.com/charmbracelet/lipgloss"

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
