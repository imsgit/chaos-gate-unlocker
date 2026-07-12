package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	MutedForeground = color.NRGBA{R: 0xF1, G: 0xEF, B: 0xEF, A: 0x99}

	white       = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	ink         = color.NRGBA{R: 0x15, G: 0x15, B: 0x15, A: 0xFF}
	panel       = color.NRGBA{R: 0x26, G: 0x26, B: 0x26, A: 0xF5}
	inputBg     = color.NRGBA{R: 0x26, G: 0x26, B: 0x26, A: 0xAB}
	overlay     = color.NRGBA{R: 0x18, G: 0x18, B: 0x18, A: 0xFA}
	disabledBg  = color.NRGBA{R: 0x18, G: 0x18, B: 0x18, A: 0xFF}
	disabledFg  = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0x42}
	selection   = color.NRGBA{R: 0x57, G: 0x59, B: 0x5B, A: 0xFF}
	hover       = color.NRGBA{R: 0x57, G: 0x59, B: 0x5B, A: 0xB2}
	pressed     = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xE2}
	scrollTrack = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0x08}
)

type Theme struct{}

func (Theme) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch c {
	case theme.ColorNameForeground, theme.ColorNameHyperlink, theme.ColorNamePrimary:
		return white
	case theme.ColorNamePlaceHolder, theme.ColorNameScrollBar:
		return MutedForeground
	case theme.ColorNameBackground, theme.ColorNameForegroundOnPrimary:
		return ink
	case theme.ColorNameButton, theme.ColorNameShadow:
		return panel
	case theme.ColorNameInputBackground:
		return inputBg
	case theme.ColorNameOverlayBackground, theme.ColorNameMenuBackground:
		return overlay
	case theme.ColorNameDisabledButton:
		return disabledBg
	case theme.ColorNameDisabled:
		return disabledFg
	case theme.ColorNameFocus, theme.ColorNameSelection:
		return selection
	case theme.ColorNameHover:
		return hover
	case theme.ColorNamePressed:
		return pressed
	case theme.ColorNameScrollBarBackground:
		return scrollTrack
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
	case theme.SizeNameDialogRadius, theme.SizeNamePopupRadius:
		return theme.DefaultTheme().Size(theme.SizeNameInputRadius)
	}
	return theme.DefaultTheme().Size(s)
}
