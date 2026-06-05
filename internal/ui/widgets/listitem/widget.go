package listitem

import (
	"chaos-gate-unlocker/internal/features"
	"chaos-gate-unlocker/internal/objects"
	"chaos-gate-unlocker/internal/ui"
	"chaos-gate-unlocker/internal/ui/widgets/tooltip"

	"image/color"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	criticalColor = color.RGBA{R: 255, G: 50, B: 50, A: 255}
	moderateColor = color.RGBA{R: 255, G: 165, B: 0, A: 255}
	lightColor    = color.RGBA{R: 255, G: 255, B: 127, A: 255}
)

type Widget struct {
	widget.BaseWidget
	tooltip.WidgetExtend

	hoverBg   *canvas.Rectangle
	iconClass *canvas.Image
	classBox  *fyne.Container
	imgLvl    *canvas.Image

	textName   *canvas.Text
	textLvl    *canvas.Text
	textStatus *canvas.Text
}

func New() fyne.CanvasObject {
	i := &Widget{
		hoverBg:    canvas.NewRectangle(color.Transparent),
		iconClass:  ui.NewIcon(fyne.NewSize(46, 46)),
		imgLvl:     ui.NewIcon(fyne.NewSize(30, 30)),
		textName:   canvas.NewText("", color.White),
		textLvl:    canvas.NewText("", fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameBackground, 0)),
		textStatus: canvas.NewText("", color.White),
	}

	i.imgLvl.Image = ui.DecodeIcon(ui.GetWidgetUnitLevelIcon())
	i.classBox = container.NewStack(i.iconClass)

	i.textName.TextStyle = fyne.TextStyle{Bold: true}
	i.textLvl.TextStyle = fyne.TextStyle{Bold: true}
	i.textStatus.TextSize = 12

	i.ExtendBaseWidget(i)
	return i
}

func (i *Widget) ExtendBaseWidget(wid fyne.Widget) {
	i.ExtendToolTipWidget(wid)
	i.BaseWidget.ExtendBaseWidget(wid)
}

func (i *Widget) MinSize() fyne.Size {
	i.ExtendBaseWidget(i)
	return fyne.NewSize(0, 54)
}

func (i *Widget) MouseIn(e *desktop.MouseEvent) {
	if !tooltip.OverlayShown(i) {
		i.WidgetExtend.MouseIn(e)
	}
	i.hoverBg.FillColor = i.Theme().Color(theme.ColorNameHover, fyne.CurrentApp().Settings().ThemeVariant())
	i.hoverBg.Refresh()
}

func (i *Widget) MouseMoved(e *desktop.MouseEvent) {
	i.WidgetExtend.MouseMoved(e)
}

func (i *Widget) MouseOut() {
	i.WidgetExtend.MouseOut()
	i.hoverBg.FillColor = color.Transparent
	i.hoverBg.Refresh()
}

func (i *Widget) CreateRenderer() fyne.WidgetRenderer {
	i.hoverBg.CornerRadius = i.Theme().Size(theme.SizeNameSelectionRadius)

	classContainer := container.NewPadded(i.classBox)

	lvlContainer := container.NewPadded(container.NewCenter(
		i.imgLvl,
		i.textLvl,
	))

	nameContainer := container.NewCenter(container.NewVBox(
		i.textName,
		i.textStatus,
	))

	return widget.NewSimpleRenderer(
		container.NewStack(i.hoverBg,
			container.NewBorder(nil, nil, container.NewHBox(
				classContainer,
				lvlContainer,
				nameContainer), nil,
			)))
}

func (i *Widget) Bind(val interface{}) {
	var name, class, lvl string
	var healthStatus int
	var noPilot, underRepair, sideMission bool

	switch object := val.(type) {
	case *objects.KnightState:
		name = override(object.GivenName, object.GivenNameOverride) + " " +
			override(features.Surnames[object.SurnameIndex], object.SurnameOverride)
		class, lvl = parseClassLvl(object.CurrentLevelData.Key)
		healthStatus = object.HealthState.Status
		sideMission = object.CurrentSideMission.MissionID != ""
	case *objects.DreadnoughtState:
		name = override(object.GivenName, object.GivenNameOverride) + " " + features.Surnames[object.SurnameIndex]
		class, lvl = parseClassLvl(object.CurrentLevelData.Key)
		healthStatus = object.HealthState.Status
		noPilot = !object.HasPilot
		underRepair = object.HealthState.RecoveryTimeLeft > 0
		sideMission = object.CurrentSideMission.MissionID != ""
	case *objects.AssassinState:
		name = override(object.GivenName, object.GivenNameOverride) + " " + features.AssassinSurnames[object.SurnameIndex]
		class, lvl = parseClassLvl(object.CurrentLevelData.Key)
		healthStatus = object.HealthState.Status
		sideMission = object.CurrentSideMission.MissionID != ""
	}

	i.iconClass.Resource = nil
	i.iconClass.Image = ui.DecodeIcon(ui.GetIconByName(class))
	i.iconClass.Refresh()
	i.classBox.Refresh()
	i.SetToolTip(splitOnCapital(class))

	i.textName.Text = name
	i.textLvl.Text = lvl
	i.textStatus.Color = ui.MutedForeground
	i.textStatus.Text = "Battle ready"

	switch healthStatus {
	case 3:
		i.textStatus.Color = criticalColor
		if class == features.DreadnoughtClass {
			i.textStatus.Text = "Unavailable - Critical damage"
			if noPilot {
				i.textStatus.Text = "Unavailable - No pilot"
			}
		} else {
			i.textStatus.Text = "Unavailable - Critical wound"
		}
	case 2:
		i.textStatus.Color = moderateColor
		if class == features.DreadnoughtClass {
			i.textStatus.Text += " - Damage"
			if underRepair {
				i.textStatus.Text = "Unavailable - Under repair"
			}
			if noPilot {
				i.textStatus.Text = "Unavailable - No pilot"
				i.textStatus.Color = criticalColor
			}
		} else {
			i.textStatus.Text += " - Wound"
		}
	case 1:
		i.textStatus.Color = lightColor
		if class == features.DreadnoughtClass {
			i.textStatus.Text += " - Light damage"
			if underRepair {
				i.textStatus.Text = "Unavailable - Under repair"
				i.textStatus.Color = moderateColor
			}
			if noPilot {
				i.textStatus.Text = "Unavailable - No pilot"
				i.textStatus.Color = criticalColor
			}
		} else {
			i.textStatus.Text += " - Light wound"
		}
	}

	if sideMission {
		i.textStatus.Text = "Unavailable - On mission"
		i.textStatus.Color = moderateColor
	}

	i.textName.Refresh()
	i.textLvl.Refresh()
	i.textStatus.Refresh()
}

func override(base, over string) string {
	if over != "" {
		return over
	}
	return base
}

func parseClassLvl(s string) (string, string) {
	splits := strings.Split(s, "_")
	if len(splits) == 2 {
		return splits[0], splits[1]
	}
	return "", ""
}

var re = regexp.MustCompile("([A-Z][a-z]*)")

func splitOnCapital(s string) string {
	if s == features.TechmarineClass {
		return "Techmarine"
	}
	return strings.Join(re.FindAllString(s, -1), " ")
}
