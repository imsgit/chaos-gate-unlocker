package main

import (
	"sort"
	"strings"

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

func showSavePicker(w fyne.Window, names []string, onPick func(name string), onOpenDir func()) {
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
				name := current[i]
				if j := strings.IndexByte(name, '_'); j >= 0 {
					name = name[j+1:]
				}
				item.Bind(name)
			}
		},
	)
	savesList.HideSeparators = true

	slotsList := widget.NewList(
		func() int { return len(slots) },
		savelistitem.New,
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if item, ok := o.(*savelistitem.Widget); ok {
				item.Bind("Slot " + slots[i])
			}
		},
	)
	slotsList.HideSeparators = true

	body := container.NewGridWithColumns(2, dragscroll.List(slotsList), dragscroll.List(savesList))

	d := dialog.NewCustomWithoutButtons("Open game save file", body, w)
	buttons := make([]fyne.CanvasObject, 0, 2)
	if onOpenDir != nil {
		buttons = append(buttons, widget.NewButtonWithIcon("Open folder", theme.FolderOpenIcon(), onOpenDir))
	}
	buttons = append(buttons, widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), d.Hide))
	d.SetButtons(buttons)
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
