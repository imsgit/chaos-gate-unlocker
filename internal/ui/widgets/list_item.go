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

func (i *ListItem) MinSize() fyne.Size {
	i.ExtendBaseWidget(i)
	return fyne.NewSize(0, 54)
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

func (i *ListItem) Bind(val interface{}) {
	var name, class, lvl string
	var healthStatus int
	var noPilot, underRepair, sideMission bool

	switch object := val.(type) {
	case *objects.KnightState:
		givenName := object.GivenName
		if object.GivenNameOverride != "" {
			givenName = decodeRussianUnicode([]byte(object.GivenNameOverride))
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
			givenName = decodeRussianUnicode([]byte(object.GivenNameOverride))
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
			givenName = decodeRussianUnicode([]byte(object.GivenNameOverride))
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

var russianUnicodeMappings = map[[3]byte]string{
	{228, 128, 149}: "ё", {228, 130, 147}: "й", {228, 129, 164}: "ц", {228, 128, 180}: "у",
	{228, 130, 163}: "к", {228, 129, 147}: "е", {228, 131, 147}: "н", {228, 128, 179}: "г",
	{228, 130, 132}: "ш", {228, 130, 148}: "щ", {228, 129, 179}: "з", {228, 129, 148}: "х",
	{228, 130, 164}: "ъ", {228, 129, 132}: "ф", {228, 130, 180}: "ы", {228, 128, 163}: "в",
	{228, 128, 131}: "а", {228, 131, 179}: "п", {228, 128, 132}: "р", {228, 131, 163}: "о",
	{228, 130, 179}: "л", {228, 129, 131}: "д", {228, 129, 163}: "ж", {228, 131, 148}: "э",
	{228, 131, 180}: "я", {228, 129, 180}: "ч", {228, 128, 148}: "с", {228, 131, 131}: "м",
	{228, 130, 131}: "и", {228, 128, 164}: "т", {228, 131, 132}: "ь", {228, 128, 147}: "б",
	{228, 131, 164}: "ю", {228, 128, 144}: "Ё", {228, 130, 145}: "Й", {228, 129, 162}: "Ц",
	{228, 128, 178}: "У", {228, 130, 161}: "К", {228, 129, 145}: "Е", {228, 131, 145}: "Н",
	{228, 128, 177}: "Г", {228, 130, 130}: "Ш", {228, 130, 146}: "Щ", {228, 129, 177}: "З",
	{228, 129, 146}: "Х", {228, 130, 162}: "Ъ", {228, 129, 130}: "Ф", {228, 130, 178}: "Ы",
	{228, 128, 161}: "В", {228, 128, 129}: "А", {228, 131, 177}: "П", {228, 128, 130}: "Р",
	{228, 131, 161}: "О", {228, 130, 177}: "Л", {228, 129, 129}: "Д", {228, 129, 161}: "Ж",
	{228, 131, 146}: "Э", {228, 131, 178}: "Я", {228, 129, 178}: "Ч", {228, 128, 146}: "С",
	{228, 131, 129}: "М", {228, 130, 129}: "И", {228, 128, 162}: "Т", {228, 131, 130}: "Ь",
	{228, 128, 145}: "Б", {228, 131, 162}: "Ю", {225, 137, 161}: "№",
}

func decodeRussianUnicode(data []byte) string {
	var result []rune
	for i := 0; i < len(data); i++ {
		if (data[i] == 228 || data[i] == 225) && i+2 < len(data) {
			key := [3]byte{data[i], data[i+1], data[i+2]}
			if val, ok := russianUnicodeMappings[key]; ok {
				result = append(result, []rune(val)...)
				i += 2
				continue
			}
		}
		result = append(result, rune(data[i]))
	}
	return string(result)
}
