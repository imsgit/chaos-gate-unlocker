package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Theme struct{}

func (Theme) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch c {
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	case theme.ColorNameHyperlink:
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0x15, G: 0x15, B: 0x15, A: 0xff}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 38, G: 38, B: 38, A: 0xf5}
	case theme.ColorNameButton:
		return color.NRGBA{R: 38, G: 38, B: 38, A: 0xf5}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 0x18, G: 0x18, B: 0x18, A: 0xff}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x42}
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0x57, G: 0x59, B: 0x5b, A: 0xff}
	case theme.ColorNameHover:
		return color.NRGBA{R: 0x57, G: 0x59, B: 0x5b, A: 0xff}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xe2}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0xf1, G: 0xef, B: 0xef, A: 0x99}
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 0xf1, G: 0xef, B: 0xef, A: 0x99}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 38, G: 38, B: 38, A: 0xf5}
	case theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 0x18, G: 0x18, B: 0x18, A: 0xfa}
	case theme.ColorNameMenuBackground:
		return color.NRGBA{R: 0x18, G: 0x18, B: 0x18, A: 0xfa}
	case theme.ColorNameSelection:
		return color.NRGBA{R: 0x57, G: 0x59, B: 0x5b, A: 0xff}
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
