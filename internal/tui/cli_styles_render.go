package tui

import (
	"github.com/charmbracelet/lipgloss"
)

//semantic render helpers to use for normal cli mode

var (
	headerStyle        lipgloss.Style
	highPriorityText   lipgloss.Style
	mediumPriorityText lipgloss.Style
	lowPriorityText    lipgloss.Style
	completedText      lipgloss.Style
	pendingText        lipgloss.Style
	mutedText          lipgloss.Style
	normalText         lipgloss.Style
)

func RenderHeader(s string) string {
	return headerStyle.PaddingRight(2).Render(s)
}

func RenderLabel(s string) string {
	return headerStyle.Render(s)
}

func RenderSectionTitle(s string) string {
	return titleStyle.Render(s)
}

func RenderValue(s string) string {
	return normalText.Render(s)
}

func RenderTodoTitle(title string, completed bool) string {
	if completed {
		return completedText.PaddingRight(3).Render(title)
	}
	return normalText.PaddingRight(3).Render(title)
}

func RenderMuted(s string) string {
	return helpStyle.Render(s)
}

func RenderError(s string) string {
	return errorStyle.Render(s)
}

func RenderStatus(completed bool) string {
	if completed {
		return completedText.PaddingLeft(3).PaddingRight(2).Render("DONE")
	}
	return pendingText.PaddingLeft(3).PaddingRight(2).Render("PENDING")
}

func RenderTimestamp(s string) string {
	return mutedText.PaddingLeft(2).Render(s)
}

func RenderPriority(priority string) string {
	switch priority {
	case "High":
		return highPriorityText.PaddingRight(3).Render("High")
	case "Medium":
		return mediumPriorityText.PaddingRight(3).Render("Medium")
	default:
		return lowPriorityText.PaddingRight(6).Render("Low")
	}
}
