package statuslabel

import (
	"chaos-gate-unlocker/internal/ui/widgets/tooltip"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const missionText = "·   MISSION"

type Widget struct {
	widget.BaseWidget

	text    *widget.Label
	mission *missionLabel
	box     *fyne.Container
}

func New() *Widget {
	l := &Widget{text: widget.NewLabel(""), mission: newMissionLabel()}
	l.mission.Hide()
	l.box = container.New(layout.NewCustomPaddedHBoxLayout(0), l.text, l.mission)
	l.ExtendBaseWidget(l)
	return l
}

func (l *Widget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(l.box)
}

func (l *Widget) Set(text, missionToolTip string) {
	l.text.SetText(text)
	l.mission.DismissToolTip()
	l.mission.SetToolTip(missionToolTip)

	if missionToolTip == "" {
		l.mission.Hide()
	} else {
		l.mission.Show()
	}
	l.Refresh()

	if missionToolTip != "" {
		l.mission.PopUpToolTip()
	}
}

type missionLabel struct {
	*widget.RichText
	tooltip.WidgetExtend
}

func newMissionLabel() *missionLabel {
	m := &missionLabel{RichText: widget.NewRichText(&widget.TextSegment{
		Text:  missionText,
		Style: widget.RichTextStyle{Inline: true, ColorName: theme.ColorNameError},
	})}
	m.ExtendBaseWidget(m)
	m.ExtendToolTipWidget(m)
	return m
}
