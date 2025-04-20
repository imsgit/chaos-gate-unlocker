package widgets

import (
	"chaos-gate-unlocker/internal/features"
	"chaos-gate-unlocker/internal/objects"
	"chaos-gate-unlocker/internal/ui"

	"image/color"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	criticalColor = color.RGBA{R: 255, G: 50, B: 50, A: 255}
	moderateColor = color.RGBA{R: 255, G: 165, B: 0, A: 255}
	lightColor    = color.RGBA{R: 255, G: 255, B: 127, A: 255}

	re = regexp.MustCompile("([A-Z][a-z]*)")
)

type ListItem struct {
	widget.BaseWidget

	iconClass *Icon
	imgLvl    *canvas.Image

	textName   *canvas.Text
	textLvl    *canvas.Text
	textStatus *canvas.Text
}

func NewListItem() fyne.CanvasObject {
	i := &ListItem{
		iconClass:  NewIcon(),
		imgLvl:     canvas.NewImageFromResource(ui.GetWidgetUnitLevelIcon()),
		textName:   canvas.NewText("", color.White),
		textLvl:    canvas.NewText("", color.Black),
		textStatus: canvas.NewText("", color.White),
	}

	i.imgLvl.FillMode = canvas.ImageFillContain
	i.imgLvl.SetMinSize(fyne.NewSize(32, 32))

	i.textName.TextStyle = fyne.TextStyle{Bold: true}
	i.textLvl.TextStyle = fyne.TextStyle{Bold: true}
	i.textStatus.TextSize = 12

	i.ExtendBaseWidget(i)
	return i
}

func (i *ListItem) CreateRenderer() fyne.WidgetRenderer {
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
		container.NewBorder(nil, nil, container.NewHBox(
			classContainer,
			lvlContainer,
			nameContainer), nil,
		))
}

func (i *ListItem) MinSize() fyne.Size {
	return fyne.NewSize(0, 54)
}

func (i *ListItem) Bind(val interface{}) {
	var name, class, lvl string
	var healthStatus int
	var noPilot, underRepair, sideMission bool

	switch object := val.(type) {
	case *objects.KnightState:
		givenName := object.GivenName
		if object.GivenNameOverride != "" {
			givenName = object.GivenNameOverride
		}
		name = givenName + " " + features.Surnames[object.SurnameIndex]
		if object.SurnameOverride != "" {
			name = givenName + " " + object.SurnameOverride
		}
		class, lvl = parseClassLvl(object.CurrentLevelData.Key)
		healthStatus = object.HealthState.Status
		sideMission = object.CurrentSideMission.MissionID != ""
	case *objects.DreadnoughtState:
		givenName := object.GivenName
		if object.GivenNameOverride != "" {
			givenName = object.GivenNameOverride
		}
		name = givenName + " " + features.Surnames[object.SurnameIndex]
		class, lvl = parseClassLvl(object.CurrentLevelData.Key)
		healthStatus = object.HealthState.Status
		noPilot = !object.HasPilot
		underRepair = object.HealthState.RecoveryTimeLeft > 0
		sideMission = object.CurrentSideMission.MissionID != ""
	case *objects.AssassinState:
		givenName := object.GivenName
		if object.GivenNameOverride != "" {
			givenName = object.GivenNameOverride
		}
		name = givenName + " " + features.AssassinSurnames[object.SurnameIndex]
		class, lvl = parseClassLvl(object.CurrentLevelData.Key)
		healthStatus = object.HealthState.Status
		sideMission = object.CurrentSideMission.MissionID != ""
	}

	i.iconClass.SetResource(ui.GetIconByName(class))
	i.iconClass.SetToolTip(splitOnCapital(class))

	i.textName.Text = name
	i.textLvl.Text = lvl
	i.textStatus.Color = fyne.CurrentApp().Settings().Theme().Color(theme.ColorNamePrimary, 0)
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

func parseClassLvl(s string) (string, string) {
	splits := strings.Split(s, "_")
	if len(splits) == 2 {
		class := splits[0]
		lvl := splits[1]
		return class, lvl
	}
	return "", ""
}

func splitOnCapital(s string) string {
	if s == features.TechmarineClass {
		return "Techmarine"
	}
	return strings.Join(re.FindAllString(s, -1), " ")
}
