package main

import (
	"maps"
	"slices"
	"strings"

	"chaos-gate-unlocker/internal/save"
	"chaos-gate-unlocker/internal/ui/widgets/dragscroll"
	"chaos-gate-unlocker/internal/ui/widgets/savelistitem"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func slotOf(name string) string {
	prefix, _, ok := strings.Cut(name, "_")
	if !ok {
		return ""
	}
	return prefix
}

type oneTwoLayout struct {
	btns fyne.CanvasObject
}

func (oneTwoLayout) MinSize(o []fyne.CanvasObject) fyne.Size {
	a, b := o[0].MinSize(), o[1].MinSize()
	return fyne.NewSize(a.Width+b.Width+theme.Padding(), fyne.Max(a.Height, b.Height))
}

func (l *oneTwoLayout) Layout(o []fyne.CanvasObject, s fyne.Size) {
	pad := theme.Padding()
	lw := (s.Width - pad) / 4
	if l.btns != nil {
		if w := s.Width/2 - l.btns.MinSize().Width/2 - pad; w > pad {
			lw = w
		}
	}
	o[0].Resize(fyne.NewSize(lw, s.Height))
	o[0].Move(fyne.NewPos(0, 0))
	o[1].Resize(fyne.NewSize(s.Width-pad-lw, s.Height))
	o[1].Move(fyne.NewPos(lw+pad, 0))
}

func showSavePicker(w fyne.Window, infos map[string]save.Info, onPick func(name string), onOpenDir func()) {
	bySlot := map[string][]string{}
	for n := range infos {
		s := slotOf(n)
		bySlot[s] = append(bySlot[s], n)
	}
	slots := slices.Sorted(maps.Keys(bySlot))
	for _, s := range slots {
		slices.Sort(bySlot[s])
	}

	var current []string

	savesList := widget.NewList(
		func() int { return len(current) },
		savelistitem.New,
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if item, ok := o.(*savelistitem.Widget); ok {
				file := current[i]
				in := infos[file]
				title := in.Title
				if title == "" {
					title = file[strings.IndexByte(file, '_')+1:]
					if rest, ok := strings.CutPrefix(title, "Slot"); ok {
						title = "Save" + rest
					}
					title = strings.TrimSuffix(title, ".gksave")
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
				item.Bind(save.SlotLabel(slots[i]), "")
			}
		},
	)
	slotsList.HideSeparators = true

	layout := &oneTwoLayout{}
	body := container.New(
		layout,
		dragscroll.List(slotsList), dragscroll.List(savesList),
	)

	d := dialog.NewCustomWithoutButtons(" Save selection", body, w)
	browse := widget.NewButtonWithIcon("Browse", theme.FolderOpenIcon(), onOpenDir)
	browse.Importance = widget.HighImportance
	btnRow := container.NewGridWithRows(1, widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), d.Hide), browse)
	layout.btns = btnRow
	d.SetButtons([]fyne.CanvasObject{btnRow})
	d.Resize(fyne.NewSize(520, 440))

	slotsList.OnSelected = func(i widget.ListItemID) {
		current = bySlot[slots[i]]
		savesList.Refresh()
		savesList.ScrollToTop()
	}

	savesList.OnSelected = func(i widget.ListItemID) {
		name := current[i]
		d.Hide()
		onPick(name)
	}

	slotsList.Select(0)

	d.Show()
}

func showSaveConfirm(w fyne.Window, do func()) {
	msg := widget.NewLabelWithStyle(
		"\n\n\nThis will override the existing save file. Are you sure?\nPlease make a backup if needed.\n\n\n",
		fyne.TextAlignCenter, fyne.TextStyle{})

	d := dialog.NewCustomWithoutButtons(" Save confirmation", msg, w)
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
