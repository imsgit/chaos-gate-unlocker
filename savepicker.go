package main

import (
	"sort"
	"strconv"
	"strings"

	"chaos-gate-unlocker/internal/saveinfo"
	"chaos-gate-unlocker/internal/ui/widgets/dragscroll"
	"chaos-gate-unlocker/internal/ui/widgets/savelistitem"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func slotOf(name string) string {
	if i := strings.IndexByte(name, '_'); i >= 0 {
		return name[:i]
	}
	return ""
}

func slotLabel(s string) string {
	if n, err := strconv.Atoi(s); err == nil {
		return "Slot " + strconv.Itoa(n+1)
	}
	return s
}

type oneTwoLayout struct {
	btns func() fyne.CanvasObject
}

func (oneTwoLayout) MinSize(o []fyne.CanvasObject) fyne.Size {
	a, b := o[0].MinSize(), o[1].MinSize()
	return fyne.NewSize(a.Width+b.Width+theme.Padding(), fyne.Max(a.Height, b.Height))
}

func (l oneTwoLayout) Layout(o []fyne.CanvasObject, s fyne.Size) {
	pad := theme.Padding()
	lw := (s.Width - pad) / 4
	if l.btns != nil {
		if b := l.btns(); b != nil {
			if w := s.Width/2 - b.MinSize().Width/2 - pad; w > pad {
				lw = w
			}
		}
	}
	o[0].Resize(fyne.NewSize(lw, s.Height))
	o[0].Move(fyne.NewPos(0, 0))
	o[1].Resize(fyne.NewSize(s.Width-pad-lw, s.Height))
	o[1].Move(fyne.NewPos(lw+pad, 0))
}

func showSavePicker(w fyne.Window, names []string, info func(name string) saveinfo.Info, onPick func(name string), onOpenDir func()) {
	slots := make([]string, 0)
	bySlot := map[string][]string{}
	for _, n := range names {
		s := slotOf(n)
		if _, ok := bySlot[s]; !ok {
			slots = append(slots, s)
		}
		bySlot[s] = append(bySlot[s], n)
	}
	sort.Strings(slots)
	for _, s := range slots {
		sort.Strings(bySlot[s])
	}

	var current []string

	savesList := widget.NewList(
		func() int { return len(current) },
		savelistitem.New,
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if item, ok := o.(*savelistitem.Widget); ok {
				file := current[i]
				disp := file
				if j := strings.IndexByte(disp, '_'); j >= 0 {
					disp = disp[j+1:]
				}
				if strings.HasPrefix(disp, "Slot") {
					disp = "Save" + disp[len("Slot"):]
				}
				disp = strings.TrimSuffix(disp, ".gksave")

				in := info(file)
				title := in.Title
				if title == "" {
					title = disp
				}
				item.Bind(title, in.Detail)
			}
		},
	)
	savesList.HideSeparators = true

	slotsList := widget.NewList(
		func() int { return len(slots) },
		savelistitem.New,
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if item, ok := o.(*savelistitem.Widget); ok {
				item.Bind(slotLabel(slots[i]), "")
			}
		},
	)
	slotsList.HideSeparators = true

	var btnRow *fyne.Container
	body := container.New(
		oneTwoLayout{btns: func() fyne.CanvasObject { return btnRow }},
		dragscroll.List(slotsList), dragscroll.List(savesList),
	)

	d := dialog.NewCustomWithoutButtons("Open game save file", body, w)
	buttons := make([]fyne.CanvasObject, 0, 2)
	buttons = append(buttons, widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), d.Hide))
	if onOpenDir != nil {
		browse := widget.NewButtonWithIcon("Browse", theme.FolderOpenIcon(), onOpenDir)
		browse.Importance = widget.HighImportance
		buttons = append(buttons, browse)
	}
	btnRow = container.NewGridWithRows(1, buttons...)
	d.SetButtons([]fyne.CanvasObject{btnRow})
	d.Resize(fyne.NewSize(520, 440))

	slotsList.OnSelected = func(i widget.ListItemID) {
		current = bySlot[slots[i]]
		savesList.UnselectAll()
		savesList.Refresh()
		savesList.ScrollToTop()
	}

	savesList.OnSelected = func(i widget.ListItemID) {
		name := current[i]
		d.Hide()
		onPick(name)
	}

	if len(slots) > 0 {
		slotsList.Select(0)
	}

	d.Show()
}

func showSaveConfirm(w fyne.Window, do func()) {
	msg := widget.NewLabelWithStyle(
		"\n\n\nThis will override the existing save file. Are you sure?\nPlease make a backup if needed.\n\n\n",
		fyne.TextAlignCenter, fyne.TextStyle{})

	d := dialog.NewCustomWithoutButtons("Save confirmation", msg, w)
	d.SetIcon(theme.QuestionIcon())
	save := widget.NewButtonWithIcon("Save", theme.ConfirmIcon(), func() {
		d.Hide()
		do()
	})
	save.Importance = widget.HighImportance
	cancel := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), d.Hide)
	d.SetButtons([]fyne.CanvasObject{cancel, save})

	d.Show()
}
