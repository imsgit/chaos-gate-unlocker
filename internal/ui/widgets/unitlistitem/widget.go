package unitlistitem

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
	imgLvl    *canvas.Image

	textName   *canvas.Text
	textLvl    *canvas.Text
	textStatus *canvas.Text
}

func New() fyne.CanvasObject {
	i := &Widget{
		hoverBg:    canvas.NewRectangle(color.Transparent),
		iconClass:  ui.NewIconImage(fyne.NewSize(46, 46)),
		imgLvl:     ui.NewIconImage(fyne.NewSize(30, 30)),
		textName:   canvas.NewText("", color.White),
		textLvl:    canvas.NewText("", fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameBackground, 0)),
		textStatus: canvas.NewText("", color.White),
	}

	i.imgLvl.Image = ui.DecodeMasked(ui.WidgetUnitLevelIcon())

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
	return fyne.NewSize(0, 54)
}

func (i *Widget) MouseIn(e *desktop.MouseEvent) {
	i.MouseInUnlessOverlay(e)
	i.hoverBg.FillColor = i.Theme().Color(theme.ColorNameHover, fyne.CurrentApp().Settings().ThemeVariant())
	i.hoverBg.Refresh()
}

func (i *Widget) MouseOut() {
	i.WidgetExtend.MouseOut()
	i.hoverBg.FillColor = color.Transparent
	i.hoverBg.Refresh()
}

func (i *Widget) CreateRenderer() fyne.WidgetRenderer {
	i.hoverBg.CornerRadius = i.Theme().Size(theme.SizeNameSelectionRadius)

	classContainer := container.NewPadded(i.iconClass)

	lvlContainer := container.NewPadded(container.NewCenter(
		i.imgLvl,
		i.textLvl,
	))

	nameContainer := container.NewCenter(container.NewVBox(
		i.textName,
		i.textStatus,
	))

	return widget.NewSimpleRenderer(
		container.NewStack(i.hoverBg, container.NewHBox(
			classContainer,
			lvlContainer,
			nameContainer,
		)))
}

func (i *Widget) Bind(val any) {
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

	iconClass := ui.DecodeMasked(ui.IconByName(class))
	i.SetToolTip(splitOnCapital(class))

	status := "Battle ready"
	var statusColor color.Color = ui.MutedForeground

	isDread := class == features.DreadnoughtClass
	switch healthStatus {
	case 3:
		statusColor = criticalColor
		if isDread {
			status = "Unavailable - Critical damage"
			if noPilot {
				status = "Unavailable - No pilot"
			}
		} else {
			status = "Unavailable - Critical wound"
		}
	case 1, 2:
		statusColor = lightColor
		damage, wound := " - Damage", " - Wound"
		if healthStatus == 1 {
			damage, wound = " - Light damage", " - Light wound"
		}
		if isDread {
			status += damage
			if underRepair {
				status = "Unavailable - Under repair"
			}
			if noPilot {
				status = "Unavailable - No pilot"
				statusColor = criticalColor
			}
		} else {
			status += wound
		}
	}

	if sideMission {
		status = "Unavailable - On mission"
		statusColor = moderateColor
	}

	dim := 0.0
	if strings.HasPrefix(status, "Unavailable") {
		dim = 0.5
	}
	if i.iconClass.Image != iconClass || i.iconClass.Translucency != dim {
		i.iconClass.Image = iconClass
		i.iconClass.Translucency = dim
		i.iconClass.Refresh()
	}
	if i.imgLvl.Translucency != dim {
		i.imgLvl.Translucency = dim
		i.imgLvl.Refresh()
	}

	setText(i.textName, name, i.textName.Color)
	setText(i.textLvl, lvl, i.textLvl.Color)
	setText(i.textStatus, status, statusColor)
}

func setText(t *canvas.Text, text string, c color.Color) {
	if t.Text == text && t.Color == c {
		return
	}
	t.Text, t.Color = text, c
	t.Refresh()
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
