package main

import (
	"chaos-gate-unlocker/internal/features"
	"chaos-gate-unlocker/internal/files"
	"chaos-gate-unlocker/internal/ui"
	"chaos-gate-unlocker/internal/ui/anim"
	"chaos-gate-unlocker/internal/ui/widgets/dragscroll"
	"chaos-gate-unlocker/internal/ui/widgets/dropdown"
	"chaos-gate-unlocker/internal/ui/widgets/progress"
	"chaos-gate-unlocker/internal/ui/widgets/statuslabel"
	"chaos-gate-unlocker/internal/ui/widgets/tabs"
	"chaos-gate-unlocker/internal/ui/widgets/toggle"
	"chaos-gate-unlocker/internal/ui/widgets/tooltip"
	"chaos-gate-unlocker/internal/ui/widgets/unitlistitem"

	"fmt"
	"image/color"
	"net/url"
	"slices"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	appID      = "chaos.gate.unlocker"
	appName    = "ChaosGateUnlocker"
	version    = "%s.%d"
	websiteURL = "https://imsgit.github.io/chaos-gate-unlocker/"

	missionToolTip = "This save was made during a combat;\nThe changes apply to the star map only and won't affect the ongoing battle"
)

type feature struct {
	on        bool
	sw        *toggle.Widget
	apply     func()
	can       func() (bool, bool)
	icon      string
	name      string
	toolTip   string
	onToggled func(on bool)
}

var (
	appVersion, appBuild string

	featuresManager = features.NewManager()
	filesManager    = files.NewManager()

	repairDreadnought = &feature{apply: featuresManager.RepairDreadnought, can: featuresManager.CanRepairDreadnought,
		icon: "ActRepair", name: "Repair Dreadnought", toolTip: "Repairs the Dreadnought's damage;\nDreadnought access is required"}

	leftFeatures = []*feature{
		{apply: featuresManager.AuthorizeDreadnoughtMissions, can: featuresManager.CanAuthorizeDreadnoughtMissions,
			icon: "ActDread", name: "Authorize Dreadnought missions", toolTip: "Marks all regular missions as Technophage, including Hive missions;\nThe difficulty of the missions will increase, won't affect frigate missions;\nDreadnought access is required"},
		repairDreadnought,
		{apply: featuresManager.RestorePrognosticars, can: featuresManager.CanRestorePrognosticars,
			icon: "ActPrognosticars", name: "Restore prognosticars", toolTip: "Makes all attuned prognosticars available again"},
		{apply: featuresManager.CompleteCurrentResearch, can: featuresManager.CanCompleteCurrentResearch,
			icon: "ActComplete", name: "Complete current research", toolTip: "Completes current research project;\nAdvance time to take effect"},
		{apply: featuresManager.CompleteCurrentConstruction, can: featuresManager.CanCompleteCurrentConstruction,
			icon: "ActComplete", name: "Complete current construction", toolTip: "Completes current construction project;\nAdvance time to take effect"},
		{apply: featuresManager.UnequipMastercraftedArmor, can: featuresManager.CanUnequipMastercraftedArmor,
			icon: "ActUnequip", name: "Unequip mastercrafted armor", toolTip: "Unequips all mastercrafted armor and frees the slots;\nAlso unlocks unavailable armor, won't affect frigate missions"},
		{apply: featuresManager.UnequipMastercraftedWeapons, can: featuresManager.CanUnequipMastercraftedWeapons,
			icon: "ActUnequip", name: "Unequip mastercrafted weapons", toolTip: "Unequips all mastercrafted weapons;\nAlso unlocks unavailable weapons, won't affect frigate missions"},
	}

	rightFeatures = []*feature{
		{apply: featuresManager.UnlockPreorderItems, can: featuresManager.CanUnlockPreorderItems,
			icon: "ActPreorder", name: "Unlock pre-order items", toolTip: "Unlocks the Domina Liber Daemonica tome and Destroyer of Crys'yllix hammer"},
		{apply: featuresManager.UnlockAdvancedClasses, can: featuresManager.CanUnlockAdvancedClasses,
			icon: "Librarian", name: "Unlock advanced classes", toolTip: "Unlocks the Librarian, Paladin, Chaplain and Purifier classes;\nAdvance time to take effect"},
		{apply: featuresManager.UnlockGarranCrowe, can: featuresManager.CanUnlockGarranCrowe,
			icon: "GarranCrowe", name: "Unlock Garran Crowe", toolTip: "Unlocks castellan Garran Crowe;\nDLC access is required;\nAdvance time to take effect"},
		{apply: featuresManager.UnlockAssassins, can: featuresManager.CanUnlockAssassins,
			icon: "ActAssassins", name: "Unlock assassins", toolTip: "Unlocks imperial assassins;\nDLC access is required;\nAdvance time to take effect"},
		{apply: featuresManager.UnlockGladiusFrigate, can: featuresManager.CanUnlockGladiusFrigate,
			icon: "ActFrigate", name: "Unlock Gladius frigate", toolTip: "Unlocks the Gladius frigate, the Cleanse mission will still appear as expected;\nDLC access is required;\nAdvance time to take effect"},
		{apply: featuresManager.UnlockPuritySeals, can: featuresManager.CanUnlockPuritySeals,
			icon: "ActSeals", name: "Unlock purity seals", toolTip: "Unlocks purity seals upgrades;\nPoxus seeds access is required;\nAdvance time to take effect"},
		{apply: featuresManager.UnlockInfiniteCampaign, can: featuresManager.CanUnlockInfiniteCampaign,
			icon: "ActReaper", name: "Unlock infinite campaign", toolTip: "Removes the Exterminatus deadline so the purge can continue indefinitely;\nEndgame access is required"},
	}

	allFeatures = slices.Concat(leftFeatures, rightFeatures)

	currUnit       any
	healUnits      = map[any]bool{}
	retrainUnits   = map[any]bool{}
	talentsUnits   = map[any][][]string{}
	augmeticsUnits = map[any][][]string{}
)

func main() {
	validateScale()

	if appVersion != "" {
		build, _ := strconv.Atoi(appBuild)
		app.SetMetadata(fyne.AppMetadata{ID: appID, Name: appName, Version: appVersion, Build: build})
	}
	a := app.NewWithID(appID)

	md := a.Metadata()
	if md.Migrations == nil {
		md.Migrations = map[string]bool{}
	}
	md.Migrations["fyneDo"] = true
	app.SetMetadata(md)

	a.Settings().SetTheme(ui.Theme{})
	w := a.NewWindow("Chaos Gate Unlocker")

	var saveButton *widget.Button
	refreshSaveButton = func() {
		if len(healUnits) > 0 || len(retrainUnits) > 0 || anyDirty(augmeticsUnits) || anyDirty(talentsUnits) ||
			slices.ContainsFunc(allFeatures, func(f *feature) bool { return f.on }) {
			saveButton.Enable()
		} else {
			saveButton.Disable()
		}
	}

	filesManager.OnLoadState(featuresManager.SetState)

	for _, f := range allFeatures {
		f.sw = toggle.New(func(on bool) {
			f.on = on
			refreshSaveButton()
			if f.onToggled != nil {
				f.onToggled(on)
			}
		}, f.icon, f.name, f.toolTip)
	}

	unitsBox := container.NewVBox()
	unitsScrollBox := container.NewVScroll(unitsBox)

	healWoundSwitch := toggle.New(func(on bool) {
		delete(healUnits, currUnit)
		if on {
			healUnits[currUnit] = on
		}
		augmeticsBox := container.NewVBox()
		unitsBox.Objects[4] = augmeticsBox
		for i := 0; ; i++ {
			initNewAugmetic := len(augmeticsUnits[currUnit][1]) == i
			if sel := renderAugmetic(i, initNewAugmetic); sel != nil {
				augmeticsBox.Objects = append(augmeticsBox.Objects, sel)
				continue
			}
			if len(augmeticsUnits[currUnit][1]) > i {
				augmeticsUnits[currUnit][1][i] = ""
			}
			break
		}
		refreshSaveButton()
	}, "ActHeal", "Heal wound", "Heals the wound")
	healWoundSwitch.Hide()

	repairDamageSwitch := toggle.New(func(on bool) {
		repairDreadnought.sw.SetState(on, true)
	}, "ActRepair", "Repair damage", "Repairs the Dreadnought's damage")
	repairDamageSwitch.Hide()

	repairDreadnought.onToggled = func(on bool) {
		if repairDamageSwitch.Visible() {
			repairDamageSwitch.SetState(on, false)
		}
	}

	retrainSwitch := toggle.New(func(on bool) {
		delete(retrainUnits, currUnit)
		if on {
			retrainUnits[currUnit] = on
		}
		refreshSaveButton()
	}, "Talent_ZealousScholar", "Retrain abilities", "Refunds spent ability points;\nExtra points gained by communing with fallen knights aren't refunded")
	retrainSwitch.Hide()

	unitsBox.Objects = append(unitsBox.Objects,
		healWoundSwitch, repairDamageSwitch, retrainSwitch, container.NewVBox(), container.NewVBox())

	var units []any
	statusLabel := statuslabel.New()

	unitsList := widget.NewList(
		func() int { return len(units) },
		unitlistitem.New,
		func(id widget.ListItemID, o fyne.CanvasObject) {
			if listItem, ok := o.(*unitlistitem.Widget); ok && id < len(units) {
				listItem.Bind(units[id])
			}
		})
	unitsList.HideSeparators = true
	showSwitch := func(sw *toggle.Widget, enable, on bool) {
		sw.Show()
		if enable {
			sw.Enable()
		} else {
			sw.Disable()
			on = true
		}
		sw.SetState(on, false)
	}
	unitsList.OnSelected = func(id widget.ListItemID) {
		currUnit = units[id]

		enable, showHeal := featuresManager.CanHealUnit(currUnit)
		if showHeal {
			tip := "Heals the wound"
			if featuresManager.UnitSupportsAugmetics(currUnit) {
				tip += ";\nIf the wound was critical, you can also select a new augmetic"
			}
			healWoundSwitch.SetToolTip(tip)
			showSwitch(healWoundSwitch, enable, healUnits[currUnit])
		} else {
			showSwitch(repairDamageSwitch, enable, repairDreadnought.on)
		}

		if canRetrain, showRetrain := featuresManager.CanRetrainUnit(currUnit); showRetrain {
			showSwitch(retrainSwitch, canRetrain, retrainUnits[currUnit])
		}

		unitsBox.Objects[3] = fillDropdownBox(talentsUnits, renderTalent)
		unitsBox.Objects[4] = fillDropdownBox(augmeticsUnits, renderAugmetic)
	}
	unitsList.OnUnselected = func(widget.ListItemID) {
		healWoundSwitch.Hide()
		repairDamageSwitch.Hide()
		retrainSwitch.Hide()
		unitsBox.Objects[3] = container.NewVBox()
		unitsBox.Objects[4] = container.NewVBox()
		unitsScrollBox.ScrollToTop()
	}

	bgImg := ui.Decode(ui.AppBackgroundIcon())
	back := canvas.NewImageFromImage(bgImg)
	back.FillMode = canvas.ImageFillContain
	back.ScaleMode = canvas.ImageScaleFastest

	dim, _ := ui.Theme{}.Color(theme.ColorNameBackground, theme.VariantDark).(color.NRGBA)
	dim.A = 245
	cover := canvas.NewRectangle(dim)

	eyeGlow := anim.NewEyeGlow(bgImg)
	eyeGlowOverlay := eyeGlow.Overlay()
	eyeGlow.Animate()

	featureColumn := func(fs []*feature) *fyne.Container {
		box := container.NewVBox()
		for _, f := range fs {
			box.Add(f.sw)
		}
		return box
	}
	mainTab := &tabs.Item{Title: "Main", Icon: ui.AppTabMainIcon(),
		Content: container.NewGridWithColumns(2, featureColumn(leftFeatures), featureColumn(rightFeatures))}
	unitsTab := &tabs.Item{Title: "Units", Icon: ui.AppTabUnitsIcon(),
		Content: container.NewGridWithColumns(2, dragscroll.List(unitsList), dragscroll.Scroll(unitsScrollBox))}
	nexusURL, _ := url.Parse("https://www.nexusmods.com/warhammer40kchaosgatedaemonhunters/mods/5")
	aboutTab := &tabs.Item{Title: "About", Icon: ui.AppTabAboutIcon(),
		Content: container.NewBorder(nil, nil,
			container.NewVBox(
				widget.NewHyperlink("> Visit Nexus Mods for more information", nexusURL)),
			widget.NewLabel(fmt.Sprintf(version, a.Metadata().Version, a.Metadata().Build)))}

	var aboutAnim *fyne.Animation
	layoutTabs := tabs.New(mainTab, unitsTab, aboutTab)
	layoutTabs.OnSelected = func(item *tabs.Item) {
		switch item {
		case aboutTab:
			aboutAnim = anim.AnimateAbout(cover, dim)
		default:
			if aboutAnim != nil {
				aboutAnim.Stop()
			}
			cover.FillColor = dim
			cover.Refresh()
		}
	}
	layoutTabs.Hide()

	newAquilaImage := func() *canvas.Image {
		img := canvas.NewImageFromImage(nil)
		img.ScaleMode = canvas.ImageScaleFastest
		img.SetMinSize(fyne.NewSize(100, 0))
		img.Translucency = 1
		return img
	}
	leftAquila := newAquilaImage()
	rightAquila := newAquilaImage()

	aquila := anim.NewAquila(ui.AppLeftAquilaIcon(), ui.AppRightAquilaIcon())
	aquila.Prewarm()

	progressLine := progress.New()
	var openButton *widget.Button
	animateTop := func(open bool, onDone func()) func() {
		stop := aquila.Animate(leftAquila, rightAquila, progressLine, open, func() {
			openButton.Enable()
			if onDone != nil {
				onDone()
			}
		})
		return func() {
			stop()
			openButton.Enable()
		}
	}

	resetUI := func() {
		openButton.Disable()
		saveButton.Disable()
		layoutTabs.Hide()
		layoutTabs.SelectIndex(0)
		unitsList.UnselectAll()
		statusLabel.Set("", "")
	}

	var loadCancel func()
	beginLoad := func() {
		eyeGlow.Flash()
		loadCancel = animateTop(true, layoutTabs.Show)

		currUnit = nil
		healUnits = map[any]bool{}
		retrainUnits = map[any]bool{}
		augmeticsUnits = map[any][][]string{}
		talentsUnits = map[any][][]string{}

		resetUI()
	}

	loadData := func(name string, data []byte, err error) {
		if err == nil {
			err = filesManager.LoadBytes(name, data)
		}

		fyne.Do(func() {
			if err != nil {
				if loadCancel != nil {
					loadCancel()
				}
				dialog.ShowError(err, w)
				return
			}

			units = featuresManager.Units()
			unitsList.Refresh()

			for _, f := range allFeatures {
				toggle.Reset(f.sw, f.can)
			}

			toolTip := ""
			if filesManager.InMission() {
				toolTip = missionToolTip
			}
			statusLabel.Set(filesManager.Status(), toolTip)
		})
	}

	openButton = widget.NewButton("Open", func() {
		openFile(w, beginLoad, loadData)
	})

	saveButton = widget.NewButton("Save", func() {
		confirmSave(w, func() {
			cancel := animateTop(false, nil)

			applyChanges()

			resetUI()

			go func() {
				if err := saveFile(); err != nil {
					fyne.Do(func() {
						cancel()
						dialog.ShowError(err, w)
					})
				}
			}()
		})
	})
	saveButton.Disable()

	var bottomBar fyne.CanvasObject = statusLabel
	if showTryOnline() {
		siteURL, _ := url.Parse(websiteURL)
		tryLink := widget.NewHyperlink("> Try it online", siteURL)
		tryLink.OnTapped = func() { openWebsite(siteURL) }
		bottomBar = container.NewBorder(nil, nil, nil, tryLink, statusLabel)
	}

	content := container.NewBorder(
		container.NewBorder(nil, nil, leftAquila, rightAquila,
			container.NewVBox(openButton, saveButton, progressLine)),
		bottomBar,
		nil, nil,
		back,
		cover,
		eyeGlowOverlay,
		layoutTabs,
	)

	w.Resize(fyne.NewSize(800, 600))
	w.SetContent(tooltip.AddWindowToolTipLayer(content, w.Canvas()))
	w.CenterOnScreen()
	w.ShowAndRun()
}

type dropdownItem struct {
	ID          string
	Name        string
	Description string
}

type dropdownSpec struct {
	placeholder string
	store       map[any][][]string
	canChange   func(idx int) (bool, dropdownItem, []string)
	lookup      func(name string) dropdownItem
}

func fillDropdownBox(store map[any][][]string, render func(idx int, init bool) *dropdown.IconWidget) *fyne.Container {
	init := len(store[currUnit]) == 0
	if init {
		store[currUnit] = [][]string{{}, {}}
	}
	box := container.NewVBox()
	for i := 0; ; i++ {
		sel := render(i, init)
		if sel == nil {
			break
		}
		box.Objects = append(box.Objects, sel)
	}
	return box
}

func renderDropdown(idx int, init bool, spec dropdownSpec) *dropdown.IconWidget {
	canChange, item, options := spec.canChange(idx)
	if !canChange {
		return nil
	}

	sel := dropdown.NewIconWidget()
	sel.SetPlaceHolder(spec.placeholder)
	sel.SetOptions(options)

	if init {
		sel.SetSelected(item.Name)
		spec.store[currUnit][0] = append(spec.store[currUnit][0], item.Name)
		spec.store[currUnit][1] = append(spec.store[currUnit][1], item.Name)
	} else {
		sel.SetSelected(spec.store[currUnit][1][idx])
		item = spec.lookup(sel.Selected())
	}

	sel.SetResource(ui.IconByName(item.ID))
	sel.SetToolTip(item.Description)
	sel.SetOptionToolTip(func(opt string) string {
		return spec.lookup(opt).Description
	})
	sel.SetOptionIcon(func(opt string) fyne.Resource {
		return ui.IconByName(spec.lookup(opt).ID)
	})

	sel.OnChanged(func(newVal string) {
		changed := spec.lookup(newVal)
		sel.SetResource(ui.IconByName(changed.ID))
		sel.SetToolTip(changed.Description)
		spec.store[currUnit][1][idx] = newVal
		refreshSaveButton()
	})

	sel.OnBeforeShowPopup(func() {
		var opts []string
		for _, opt := range options {
			if !containsOpt(spec.store[currUnit][1], opt) || opt == sel.Selected() {
				opts = append(opts, opt)
			}
		}
		sel.SetOptions(opts)
	})
	return sel
}

func renderTalent(idx int, init bool) *dropdown.IconWidget {
	return renderDropdown(idx, init, dropdownSpec{
		placeholder: "(Select talent)",
		store:       talentsUnits,
		canChange: func(idx int) (bool, dropdownItem, []string) {
			canChange, talent, options := featuresManager.CanChangeUnitTalents(currUnit, idx)
			return canChange, dropdownItem(talent), options
		},
		lookup: func(name string) dropdownItem {
			return dropdownItem(featuresManager.TalentByName(name))
		},
	})
}

func renderAugmetic(idx int, init bool) *dropdown.IconWidget {
	return renderDropdown(idx, init, dropdownSpec{
		placeholder: "(Select augmetic)",
		store:       augmeticsUnits,
		canChange: func(idx int) (bool, dropdownItem, []string) {
			canChange, augmetic, options := featuresManager.CanChangeUnitAugmetics(currUnit, idx, healUnits[currUnit])
			return canChange, dropdownItem(augmetic), options
		},
		lookup: func(name string) dropdownItem {
			return dropdownItem(featuresManager.AugmeticByName(name))
		},
	})
}

func anyDirty(units map[any][][]string) bool {
	for _, v := range units {
		if !slices.Equal(v[0], v[1]) {
			return true
		}
	}
	return false
}

func containsOpt(list []string, val string) bool {
	return features.IsSkipOption(val) || slices.Contains(list, val)
}

var refreshSaveButton func()

func applyChanges() {
	for _, f := range allFeatures {
		if f.on {
			f.apply()
		}
	}
	for unit := range healUnits {
		featuresManager.HealUnit(unit)
	}
	for unit := range retrainUnits {
		featuresManager.RetrainUnit(unit)
	}
	for unit, augmetics := range augmeticsUnits {
		featuresManager.ChangeUnitAugmetics(unit, augmetics[1])
	}
	for unit, talents := range talentsUnits {
		featuresManager.ChangeUnitTalents(unit, talents[1])
	}
}
