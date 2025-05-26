package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Theme struct{}

func (Theme) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch c {
	case theme.ColorNameForeground, theme.ColorNameHyperlink:
		return color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0x15, G: 0x15, B: 0x15, A: 0xFF}
	case theme.ColorNameInputBackground, theme.ColorNameButton, theme.ColorNameShadow:
		return color.NRGBA{R: 0x26, G: 0x26, B: 0x26, A: 0xF5}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 0x18, G: 0x18, B: 0x18, A: 0xFF}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0x42}
	case theme.ColorNameFocus, theme.ColorNameHover, theme.ColorNameSelection:
		return color.NRGBA{R: 0x57, G: 0x59, B: 0x5B, A: 0xFF}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xE2}
	case theme.ColorNamePrimary, theme.ColorNameScrollBar:
		return color.NRGBA{R: 0xF1, G: 0xEF, B: 0xEF, A: 0x99}
	case theme.ColorNameOverlayBackground, theme.ColorNameMenuBackground:
		return color.NRGBA{R: 0x18, G: 0x18, B: 0x18, A: 0xFA}
	}

	return theme.DefaultTheme().Color(c, v)
}

func (Theme) Font(s fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(s)
}

func (Theme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (Theme) Size(s fyne.ThemeSizeName) float32 {
	switch s {
	case theme.SizeNameText:
		return 13.0
	case theme.SizeNameScrollBarSmall:
		return 4.0
	}
	return theme.DefaultTheme().Size(s)
}
