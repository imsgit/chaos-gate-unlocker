package main

import (
	"chaos-gate-unlocker/internal/features"
	"chaos-gate-unlocker/internal/files"
	"chaos-gate-unlocker/internal/ui"
	"chaos-gate-unlocker/internal/ui/anim"
	"chaos-gate-unlocker/internal/ui/widgets/dropdown"
	"chaos-gate-unlocker/internal/ui/widgets/listitem"
	"chaos-gate-unlocker/internal/ui/widgets/progress"
	"chaos-gate-unlocker/internal/ui/widgets/toggle"
	"chaos-gate-unlocker/internal/ui/widgets/tooltip"

	"context"
	"fmt"
	"net/url"
	"runtime"
	"runtime/debug"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	version    = "Ver: %s.%d | Author: imsgit | 2026-06-20"
	websiteURL = "https://imsgit.github.io/chaos-gate-unlocker/"
)

var (
	featuresManager = features.NewManager()
	filesManager    = files.NewManager()

	removeMarketingWeapons, unlockPreorderItems, unlockAdvancedClasses, unlockPuritySeals, unlockAssassins,
	restorePrognosticars, unlockGarranCrowe, authorizeDreadnoughtMissions, repairDreadnought, unlockGladiusFrigate,
	completeCurrentResearch, completeCurrentConstruction, unequipMastercraftedWeapons, unequipMastercraftedArmor bool

	currUnit       any
	healUnits      = map[any]bool{}
	retrainUnits   = map[any]bool{}
	talentsUnits   = map[any][][]string{}
	augmeticsUnits = map[any][][]string{}
)

var featureActions = []struct {
	flag  *bool
	apply func()
}{
	{&unlockPuritySeals, featuresManager.UnlockPuritySeals},
	{&repairDreadnought, featuresManager.RepairDreadnought},
	{&unlockPreorderItems, featuresManager.UnlockPreorderItems},
	{&unlockAdvancedClasses, featuresManager.UnlockAdvancedClasses},
	{&restorePrognosticars, featuresManager.RestorePrognosticars},
	{&unlockGarranCrowe, featuresManager.UnlockGarranCrowe},
	{&authorizeDreadnoughtMissions, featuresManager.AuthorizeDreadnoughtMissions},
	{&unlockGladiusFrigate, featuresManager.UnlockGladiusFrigate},
	{&completeCurrentConstruction, featuresManager.CompleteCurrentConstruction},
	{&completeCurrentResearch, featuresManager.CompleteCurrentResearch},
	{&unequipMastercraftedWeapons, featuresManager.UnequipMastercraftedWeapons},
	{&unequipMastercraftedArmor, featuresManager.UnequipMastercraftedArmor},
	{&removeMarketingWeapons, featuresManager.RemoveMarketingWeapons},
	{&unlockAssassins, featuresManager.UnlockAssassins},
}

func main() {
	debug.SetGCPercent(50)

	validateScale()

	a := app.NewWithID("chaos.gate.unlocker")

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
		canApplyChanges := len(healUnits) > 0 || len(retrainUnits) > 0 ||
			anyDirty(augmeticsUnits) || anyDirty(talentsUnits)
		for _, action := range featureActions {
			if *action.flag {
				canApplyChanges = true
				break
			}
		}

		saveButton.Disable()
		if canApplyChanges {
			saveButton.Enable()
		}
	}

	filesManager.OnLoadState(featuresManager.ApplyState())

	authorizeDreadnoughtMissionsSwitch := boolSwitch(&authorizeDreadnoughtMissions, "ActDread", "Authorize Dreadnought missions", "Marks all regular missions as Technophage, including Hive missions;\nThe difficulty of the missions will increase, won't affect frigate missions;\nDreadnought access is required")
	restorePrognosticarsSwitch := boolSwitch(&restorePrognosticars, "ActPrognosticars", "Restore prognosticars", "Makes all attuned prognosticars available again")
	completeCurrentResearchSwitch := boolSwitch(&completeCurrentResearch, "ActComplete", "Complete current research", "Completes current research project;\nAdvance time to take effect")
	completeCurrentConstructionSwitch := boolSwitch(&completeCurrentConstruction, "ActComplete", "Complete current construction", "Completes current construction project;\nAdvance time to take effect")
	unequipMastercraftedArmorSwitch := boolSwitch(&unequipMastercraftedArmor, "ActUnequip", "Unequip mastercrafted armor", "Unequips all mastercrafted armor and frees the slots;\nAlso unlocks unavailable armor, won't affect frigate missions")
	unequipMastercraftedWeaponsSwitch := boolSwitch(&unequipMastercraftedWeapons, "ActUnequip", "Unequip mastercrafted weapons", "Unequips all mastercrafted weapons;\nAlso unlocks unavailable weapons, won't affect frigate missions")
	unlockPreorderItemsSwitch := boolSwitch(&unlockPreorderItems, "ActPreorder", "Unlock pre-order items", "Unlocks the Domina Liber Daemonica tome and Destroyer of Crys'yllix hammer")
	unlockAdvancedClassesSwitch := boolSwitch(&unlockAdvancedClasses, "Librarian", "Unlock advanced classes", "Unlocks the Librarian, Paladin, Chaplain and Purifier classes;\nAdvance time to take effect")
	unlockGarranCroweSwitch := boolSwitch(&unlockGarranCrowe, "GarranCrowe", "Unlock Garran Crowe", "Unlocks castellan Garran Crowe;\nDLC access is required;\nAdvance time to take effect")
	unlockAssassinsSwitch := boolSwitch(&unlockAssassins, "ActAssassins", "Unlock assassins", "Unlocks imperial assassins;\nDLC access is required;\nAdvance time to take effect")
	unlockGladiusFrigateSwitch := boolSwitch(&unlockGladiusFrigate, "ActFrigate", "Unlock Gladius frigate", "Unlocks the Gladius frigate, the Cleanse mission will still appear as expected;\nDLC access is required;\nAdvance time to take effect")
	unlockPuritySealsSwitch := boolSwitch(&unlockPuritySeals, "ActSeals", "Unlock purity seals", "Unlocks purity seals upgrades;\nPoxus seeds access is required;\nAdvance time to take effect")
	removeMarketingWeaponsSwitch := boolSwitch(&removeMarketingWeapons, "ActMarketing", "Remove marketing weapons", "Unequips and removes all weapons classified as Twitch drops")

	var repairDamageSwitch *toggle.Widget
	repairDreadnoughtSwitch := toggle.New(func(on bool) {
		repairDreadnought = on
		refreshSaveButton()
		if repairDamageSwitch.Visible() {
			repairDamageSwitch.SetState(on, false)
		}
	}, "ActRepair", "Repair Dreadnought", "Repairs the Dreadnought's damage;\nDreadnought access is required")

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

	repairDamageSwitch = toggle.New(func(on bool) {
		repairDreadnought = on
		refreshSaveButton()
		repairDreadnoughtSwitch.SetState(on, false)
	}, "ActRepair", "Repair damage", "Repairs the Dreadnought's damage")
	repairDamageSwitch.Hide()

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

	unitsProvider := binding.NewUntypedList()
	status := binding.NewString()

	unitsList := widget.NewListWithData(
		unitsProvider,
		listitem.New,
		func(item binding.DataItem, o fyne.CanvasObject) {
			val, ok := item.(binding.Untyped)
			listItem, ok2 := o.(*listitem.Widget)
			if ok && ok2 {
				v, _ := val.Get()
				listItem.Bind(v)
			}
		})
	unitsList.HideSeparators = true
	unitsList.OnSelected = func(id widget.ListItemID) {
		currUnit, _ = unitsProvider.GetValue(id)

		enable, showHeal := featuresManager.CanHealUnit(currUnit)
		if showHeal {
			healWoundSwitch.Show()
			if featuresManager.UnitSupportsAugmetics(currUnit) {
				healWoundSwitch.SetToolTip("Heals the wound;\nIf the wound was critical, you can also select a new augmetic")
			} else {
				healWoundSwitch.SetToolTip("Heals the wound")
			}
			if enable {
				healWoundSwitch.Enable()
				healWoundSwitch.SetState(healUnits[currUnit], false)
			} else {
				healWoundSwitch.Disable()
				healWoundSwitch.SetState(true, false)
			}
		} else {
			repairDamageSwitch.Show()
			if enable {
				repairDamageSwitch.Enable()
				repairDamageSwitch.SetState(repairDreadnought, false)
			} else {
				repairDamageSwitch.Disable()
				repairDamageSwitch.SetState(true, false)
			}
		}

		if canRetrain, showRetrain := featuresManager.CanRetrainUnit(currUnit); showRetrain {
			retrainSwitch.Show()
			if canRetrain {
				retrainSwitch.Enable()
				retrainSwitch.SetState(retrainUnits[currUnit], false)
			} else {
				retrainSwitch.Disable()
				retrainSwitch.SetState(true, false)
			}
		} else {
			retrainSwitch.Hide()
		}

		initTalents := len(talentsUnits[currUnit]) == 0
		if initTalents {
			talentsUnits[currUnit] = append(talentsUnits[currUnit], []string{}, []string{})
		}
		unitsBox.Objects[3] = fillDropdownBox(renderTalent, initTalents)

		initAugmetics := len(augmeticsUnits[currUnit]) == 0
		if initAugmetics {
			augmeticsUnits[currUnit] = append(augmeticsUnits[currUnit], []string{}, []string{})
		}
		unitsBox.Objects[4] = fillDropdownBox(renderAugmetic, initAugmetics)
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
	back.Translucency = 0.96

	eyeGlow := anim.NewEyeGlow(bgImg)
	eyeGlowOverlay := eyeGlow.Overlay()
	eyeGlow.Animate()

	mainTab := container.NewTabItemWithIcon("Main", ui.AppTabMainIcon(),
		container.NewGridWithColumns(2,
			container.NewVBox(
				authorizeDreadnoughtMissionsSwitch,
				repairDreadnoughtSwitch,
				restorePrognosticarsSwitch,
				completeCurrentResearchSwitch,
				completeCurrentConstructionSwitch,
				unequipMastercraftedArmorSwitch,
				unequipMastercraftedWeaponsSwitch),
			container.NewVBox(
				unlockPreorderItemsSwitch,
				unlockAdvancedClassesSwitch,
				unlockGarranCroweSwitch,
				unlockAssassinsSwitch,
				unlockGladiusFrigateSwitch,
				unlockPuritySealsSwitch,
				removeMarketingWeaponsSwitch)))
	unitsTab := container.NewTabItemWithIcon("Units", ui.AppTabUnitsIcon(),
		container.NewGridWithColumns(2, unitsList, unitsScrollBox))
	aboutTab := container.NewTabItemWithIcon("About", ui.AppTabAboutIcon(),
		container.NewBorder(nil, nil,
			widget.NewRichTextFromMarkdown(`
[> Visit Nexus Mods for more information](https://www.nexusmods.com/warhammer40kchaosgatedaemonhunters/mods/5)

[> Visit Fyne.io for app details](https://apps.fyne.io/apps/chaos.gate.unlocker.html)`),
			widget.NewRichTextFromMarkdown(fmt.Sprintf(version, a.Metadata().Version, a.Metadata().Build))))

	var acancel context.CancelFunc
	layoutTabs := container.NewAppTabs(mainTab, unitsTab, aboutTab)
	layoutTabs.SetTabLocation(container.TabLocationTrailing)
	layoutTabs.OnSelected = func(item *container.TabItem) {
		switch item {
		case aboutTab:
			var actx context.Context
			actx, acancel = context.WithCancel(context.Background())
			cancel := acancel
			go func() {
				anim.AnimateAbout(actx, back)
				cancel()
			}()
		default:
			if acancel != nil {
				acancel()
			}
			back.Translucency = 0.96
		}
	}
	layoutTabs.Hide()

	leftAquila := canvas.NewImageFromResource(ui.AppLeftAquilaIcon())
	leftAquila.ScaleMode = canvas.ImageScaleFastest
	leftAquila.SetMinSize(fyne.NewSize(100, 0))
	leftAquila.Translucency = 1

	rightAquila := canvas.NewImageFromResource(ui.AppRightAquilaIcon())
	rightAquila.ScaleMode = canvas.ImageScaleFastest
	rightAquila.SetMinSize(fyne.NewSize(100, 0))
	rightAquila.Translucency = 1

	aquila := anim.NewAquila(ui.AppLeftAquilaIcon(), ui.AppRightAquilaIcon())
	aquila.Prewarm()

	progressLine := progress.New()
	var openButton *tooltip.Button
	animateTop := func(open bool, onDone func()) context.CancelFunc {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			aquila.Animate(ctx, leftAquila, rightAquila, progressLine, open)
			fyne.DoAndWait(func() {
				openButton.Enable()
				if ctx.Err() == nil && onDone != nil {
					onDone()
				}
			})
			cancel()
		}()
		return cancel
	}

	resetUI := func() {
		openButton.Disable()
		saveButton.Disable()
		layoutTabs.Hide()
		layoutTabs.SelectIndex(0)
		unitsList.UnselectAll()
		status.Set("")
	}

	loadData := func(name string, data []byte) {
		cancel := animateTop(true, func() {
			layoutTabs.Show()
			status.Set(filesManager.Status())
		})

		healUnits = map[any]bool{}
		retrainUnits = map[any]bool{}
		augmeticsUnits = map[any][][]string{}
		talentsUnits = map[any][][]string{}

		resetUI()

		if err := filesManager.LoadBytes(name, data); err != nil {
			cancel()
			dialog.ShowError(err, w)
			return
		}

		unitsProvider.Set(featuresManager.Units())

		toggle.Reset(unlockAdvancedClassesSwitch, featuresManager.CanUnlockAdvancedClasses)
		toggle.Reset(repairDreadnoughtSwitch, featuresManager.CanRepairDreadnought)
		toggle.Reset(unlockPuritySealsSwitch, featuresManager.CanUnlockPuritySeals)
		toggle.Reset(restorePrognosticarsSwitch, featuresManager.CanRestorePrognosticars)
		toggle.Reset(unlockGarranCroweSwitch, featuresManager.CanUnlockGarranCrowe)
		toggle.Reset(authorizeDreadnoughtMissionsSwitch, featuresManager.CanAuthorizeDreadnoughtMissions)
		toggle.Reset(unlockGladiusFrigateSwitch, featuresManager.CanUnlockGladiusFrigate)
		toggle.Reset(completeCurrentResearchSwitch, featuresManager.CanCompleteCurrentResearch)
		toggle.Reset(completeCurrentConstructionSwitch, featuresManager.CanCompleteCurrentConstruction)
		toggle.Reset(unlockAssassinsSwitch, featuresManager.CanUnlockAssassins)
		toggle.ResetOn(unlockPreorderItemsSwitch, featuresManager.CanUnlockPreorderItems())
		toggle.ResetOn(unequipMastercraftedWeaponsSwitch, featuresManager.CanUnequipMastercraftedWeapons())
		toggle.ResetOn(unequipMastercraftedArmorSwitch, featuresManager.CanUnequipMastercraftedArmor())
		toggle.ResetOn(removeMarketingWeaponsSwitch, featuresManager.CanRemoveMarketingWeapons())
	}

	openButton = tooltip.NewButton("Open", func() {
		openFile(w, filesManager, loadData)
	})
	openButton.SetToolTip("Can't find your save? It's in:\n" + filesManager.DefaultLocationHint())

	saveButton = widget.NewButton("Save", func() {
		confirmSave(w, func() {
			cancel := animateTop(false, nil)

			applyChanges()

			resetUI()

			if err := saveFile(filesManager); err != nil {
				cancel()
				dialog.ShowError(err, w)
				return
			}
		})
	})
	saveButton.Disable()

	statusLabel := widget.NewLabelWithData(status)
	var tryLink *widget.Hyperlink
	if browserSupported {
		tryLink = widget.NewHyperlink("> Try it online", nil)
		tryLink.OnTapped = func() {
			if err := openInBrowser(); err != nil {
				dialog.ShowError(err, w)
			}
		}
	} else if runtime.GOOS != "js" {
		u, _ := url.Parse(websiteURL)
		tryLink = widget.NewHyperlink("> Try it online", u)
	}
	var bottomBar fyne.CanvasObject = statusLabel
	if tryLink != nil {
		bottomBar = container.NewBorder(nil, nil, nil, tryLink, statusLabel)
	}

	content := container.NewBorder(
		container.NewBorder(nil, nil, leftAquila, rightAquila,
			container.NewVBox(openButton, saveButton, progressLine)),
		bottomBar,
		nil, nil,
		back,
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

func fillDropdownBox(render func(idx int, init bool) *dropdown.IconWidget, init bool) *fyne.Container {
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
	if features.IsSkipOption(val) {
		return true
	}
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}

var refreshSaveButton func()

func applyChanges() {
	for _, action := range featureActions {
		if *action.flag {
			action.apply()
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

func boolSwitch(flag *bool, icon, name, tooltip string) *toggle.Widget {
	return toggle.New(func(on bool) {
		*flag = on
		refreshSaveButton()
	}, icon, name, tooltip)
}
