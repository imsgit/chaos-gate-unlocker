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
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"runtime/debug"
	"slices"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

const (
	version = "Ver: 1.0.0.%d | Author: imsgit | 2026-06-06"
)

var (
	featuresManager = features.NewManager()
	filesManager    = files.NewManager()

	removeMarketingWeapons, unlockPreorderItems, unlockAdvancedClasses, unlockPuritySeals, unlockAssassins,
	reattunePrognosticars, unlockGarranCrowe, authorizeDreadnoughtMissions, repairDreadnought, unlockGladiusFrigate,
	completeCurrentResearch, completeCurrentConstruction, unequipMastercraftedWeapons, unequipMastercraftedArmor bool

	currUnit       any
	healUnits      = map[any]bool{}
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
	{&reattunePrognosticars, featuresManager.ReattunePrognosticars},
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

	a := app.NewWithID("chaos.gate.unlocker")
	a.Settings().SetTheme(ui.Theme{})
	w := a.NewWindow("Chaos Gate Unlocker")

	var saveButton *widget.Button
	refreshSaveButton = func() {
		var augmeticsChanged bool
		for _, augmetics := range augmeticsUnits {
			if !slices.Equal(augmetics[0], augmetics[1]) {
				augmeticsChanged = true
				break
			}
		}

		var talentsChanged bool
		for _, talents := range talentsUnits {
			if !slices.Equal(talents[0], talents[1]) {
				talentsChanged = true
				break
			}
		}

		canApplyChanges := len(healUnits) > 0 || augmeticsChanged || talentsChanged
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

	authorizeDreadnoughtMissionsSwitch := boolSwitch(&authorizeDreadnoughtMissions, "ActDread", "Authorize Dreadnought missions", "Marks all regular missions as Technophage, including Hive missions;\nThe difficulty of the missions will increase, won't affect the frigate's chances of winning;\nDreadnought access is required")
	reattunePrognosticarsSwitch := boolSwitch(&reattunePrognosticars, "ActPrognosticars", "Reattune prognosticars", "Restores all attuned prognosticars, making them available again")
	completeCurrentResearchSwitch := boolSwitch(&completeCurrentResearch, "ActComplete", "Complete current research", "Completes current research project;\nAdvance current day to unlock")
	completeCurrentConstructionSwitch := boolSwitch(&completeCurrentConstruction, "ActComplete", "Complete current construction", "Completes current construction project;\nAdvance current day to unlock")
	unequipMastercraftedArmorSwitch := boolSwitch(&unequipMastercraftedArmor, "ActUnequip", "Unequip mastercrafted armor", "Unequips all mastercrafted armor, clearing the slots if necessary;\nAlso unlocks currently unavailable armor, won't affect the frigate's chances of winning")
	unequipMastercraftedWeaponsSwitch := boolSwitch(&unequipMastercraftedWeapons, "ActUnequip", "Unequip mastercrafted weapons", "Unequips all mastercrafted weapons;\nAlso unlocks currently unavailable weapons, won't affect the frigate's chances of winning")
	unlockPreorderItemsSwitch := boolSwitch(&unlockPreorderItems, "ActPreorder", "Unlock pre-order items", "Unlocks the Domina Liber Daemonica tome and Destroyer of Crys'yllix hammer")
	unlockAdvancedClassesSwitch := boolSwitch(&unlockAdvancedClasses, "Librarian", "Unlock advanced classes", "Unlocks the Librarian, Paladin, Chaplain and Purifier classes;\nAdvance current day to unlock")
	unlockGarranCroweSwitch := boolSwitch(&unlockGarranCrowe, "GarranCrowe", "Unlock Garran Crowe", "Unlocks castellan Garran Crowe;\nDLC access is required;\nAdvance current day to unlock")
	unlockAssassinsSwitch := boolSwitch(&unlockAssassins, "ActAssassins", "Unlock assassins", "Unlocks imperial assassins;\nDLC access is required;\nAdvance current day to unlock")
	unlockGladiusFrigateSwitch := boolSwitch(&unlockGladiusFrigate, "ActFrigate", "Unlock Gladius frigate", "Unlocks the Gladius frigate, the Cleanse mission will still appear as expected;\nDLC access is required;\nAdvance current day to unlock")
	unlockPuritySealsSwitch := boolSwitch(&unlockPuritySeals, "ActSeals", "Unlock purity seals", "Unlocks purity seals upgrades;\nPoxus seeds access is required;\nAdvance current day to unlock")
	removeMarketingWeaponsSwitch := boolSwitch(&removeMarketingWeapons, "ActMarketing", "Remove marketing weapons", "Unequips and removes all weapons classified as Twitch drops")

	var repairDamageSwitch *toggle.Widget
	repairDreadnoughtSwitch := toggle.New(func(on bool) {
		repairDreadnought = on
		refreshSaveButton()
		if repairDamageSwitch.Visible() {
			repairDamageSwitch.SetState(on, false)
		}
	}, "ActRepair", "Repair Dreadnought", "Repairs the Dreadnought's damage for free;\nDreadnought access is required")

	unitsBox := container.NewVBox()
	unitsScrollBox := container.NewVScroll(unitsBox)

	healWoundSwitch := toggle.New(func(on bool) {
		delete(healUnits, currUnit)
		if on {
			healUnits[currUnit] = on
		}

		augmeticsBox := container.NewVBox()
		unitsBox.Objects[3] = augmeticsBox
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
	}, "ActRepair", "Repair damage", "Repairs the Dreadnought's damage for free")
	repairDamageSwitch.Hide()

	unitsBox.Objects = append(unitsBox.Objects,
		healWoundSwitch, repairDamageSwitch, container.NewVBox(), container.NewVBox())

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
				repairDamageSwitch.SetState(repairDreadnought, false)
			} else {
				repairDamageSwitch.Disable()
				repairDamageSwitch.SetState(true, false)
			}
		}

		initTalents := len(talentsUnits[currUnit]) == 0
		if initTalents {
			talentsUnits[currUnit] = append(talentsUnits[currUnit], []string{}, []string{})
		}
		unitsBox.Objects[2] = fillDropdownBox(renderTalent, initTalents)

		initAugmetics := len(augmeticsUnits[currUnit]) == 0
		if initAugmetics {
			augmeticsUnits[currUnit] = append(augmeticsUnits[currUnit], []string{}, []string{})
		}
		unitsBox.Objects[3] = fillDropdownBox(renderAugmetic, initAugmetics)
	}
	unitsList.OnUnselected = func(widget.ListItemID) {
		healWoundSwitch.Hide()
		repairDamageSwitch.Hide()
		unitsBox.Objects[2] = container.NewVBox()
		unitsBox.Objects[3] = container.NewVBox()
		unitsScrollBox.ScrollToTop()
	}

	back := canvas.NewImageFromResource(ui.GetAppBackgroundIcon())
	back.FillMode = canvas.ImageFillContain
	back.ScaleMode = canvas.ImageScaleFastest
	back.Translucency = 0.96

	eyeGlow := anim.NewEyeGlow(ui.GetAppBackgroundIcon())
	eyeGlowOverlay := eyeGlow.Overlay()
	eyeGlow.Start()

	mainTab := container.NewTabItemWithIcon("Main", ui.GetAppTabMainIcon(),
		container.NewThemeOverride(container.NewGridWithColumns(2,
			container.NewVBox(
				authorizeDreadnoughtMissionsSwitch,
				repairDreadnoughtSwitch,
				reattunePrognosticarsSwitch,
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
				removeMarketingWeaponsSwitch)), ui.Theme{}))
	unitsTab := container.NewTabItemWithIcon("Units", ui.GetAppTabUnitsIcon(),
		container.NewThemeOverride(
			container.NewGridWithColumns(2, unitsList, unitsScrollBox), ui.Theme{}))
	aboutTab := container.NewTabItemWithIcon("About", ui.GetAppTabAboutIcon(),
		container.NewThemeOverride(container.NewBorder(nil, nil,
			widget.NewRichTextFromMarkdown(`
[> Visit Nexus Mods for more information](https://www.nexusmods.com/warhammer40kchaosgatedaemonhunters/mods/5)

[> Visit Fyne.io for app details](https://apps.fyne.io/apps/chaos.gate.unlocker.html)

[> Visit Reddit for discussion](https://www.reddit.com/r/ChaosGateGame/comments/1hz3s5g/chaosgateunlocker)`),
			widget.NewRichTextFromMarkdown(fmt.Sprintf(version, a.Metadata().Build))), ui.Theme{}))

	var acancel context.CancelFunc
	var tabsThemed *container.ThemeOverride
	layoutTabs := container.NewAppTabs(mainTab, unitsTab, aboutTab)
	layoutTabs.SetTabLocation(container.TabLocationTrailing)
	layoutTabs.OnSelected = func(item *container.TabItem) {
		switch item {
		case aboutTab:
			var actx context.Context
			actx, acancel = context.WithCancel(context.Background())
			go func() {
				anim.AnimateAbout(actx, back)
				acancel()
			}()
		default:
			if acancel != nil {
				acancel()
			}
			back.Translucency = 0.96
		}
	}
	layoutTabs.Hide()

	leftAquila := canvas.NewImageFromResource(ui.GetAppLeftAquilaIcon())
	leftAquila.ScaleMode = canvas.ImageScaleFastest
	leftAquila.SetMinSize(fyne.NewSize(100, 0))
	leftAquila.Translucency = 1

	rightAquila := canvas.NewImageFromResource(ui.GetAppRightAquilaIcon())
	rightAquila.ScaleMode = canvas.ImageScaleFastest
	rightAquila.SetMinSize(fyne.NewSize(100, 0))
	rightAquila.Translucency = 1

	aquila := anim.NewAquila(ui.GetAppLeftAquilaIcon(), ui.GetAppRightAquilaIcon())
	aquila.Prewarm()

	progressLine := progress.New()

	var openButton *tooltip.Button

	animateTop := func(open bool, onDone func()) context.CancelFunc {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			aquila.AnimateTop(ctx, leftAquila, rightAquila, progressLine, open)
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

	openButton = tooltip.NewButton("Open", func() {
		fileDialog := dialog.NewFileOpen(func(rc fyne.URIReadCloser, err error) {
			if rc == nil {
				return
			}

			cancel := animateTop(true, func() {
				layoutTabs.Show()
				if tabsThemed != nil {
					tabsThemed.Refresh()
				}
				status.Set(filesManager.Status())
			})

			healUnits = map[any]bool{}
			augmeticsUnits = map[any][][]string{}
			talentsUnits = map[any][][]string{}

			resetUI()

			err = filesManager.Load(rc)
			if err != nil {
				cancel()
				dialog.ShowError(err, w)
				return
			}

			unitsProvider.Set(featuresManager.Units())

			resetSwitch(unlockAdvancedClassesSwitch, featuresManager.CanUnlockAdvancedClasses)
			resetSwitch(repairDreadnoughtSwitch, featuresManager.CanRepairDreadnought)
			resetSwitch(unlockPuritySealsSwitch, featuresManager.CanUnlockPuritySeals)
			resetSwitch(reattunePrognosticarsSwitch, featuresManager.CanReattunePrognosticars)
			resetSwitch(unlockGarranCroweSwitch, featuresManager.CanUnlockGarranCrowe)
			resetSwitch(authorizeDreadnoughtMissionsSwitch, featuresManager.CanAuthorizeDreadnoughtMissions)
			resetSwitch(unlockGladiusFrigateSwitch, featuresManager.CanUnlockGladiusFrigate)
			resetSwitch(completeCurrentResearchSwitch, featuresManager.CanCompleteCurrentResearch)
			resetSwitch(completeCurrentConstructionSwitch, featuresManager.CanCompleteCurrentConstruction)
			resetSwitch(unlockAssassinsSwitch, featuresManager.CanUnlockAssassins)
			resetSwitchOn(unlockPreorderItemsSwitch, featuresManager.CanUnlockPreorderItems())
			resetSwitchOn(unequipMastercraftedWeaponsSwitch, featuresManager.CanUnequipMastercraftedWeapons())
			resetSwitchOn(unequipMastercraftedArmorSwitch, featuresManager.CanUnequipMastercraftedArmor())
			resetSwitchOn(removeMarketingWeaponsSwitch, featuresManager.CanRemoveMarketingWeapons())
		}, w)

		l, _ := storage.ListerForURI(storage.NewFileURI(filesManager.GetCurrentPath()))
		fileDialog.SetTitleText("Open game save file ../" + filesManager.SaveDir())
		fileDialog.SetConfirmText("Open")
		fileDialog.SetDismissText("Cancel")
		fileDialog.SetLocation(l)
		fileDialog.Resize(fyne.NewSize(800, 600))
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".gksave"}))
		fileDialog.Show()
	})
	openButton.SetToolTip("Can't find your save? It's in:\n" + filesManager.DefaultLocationHint())

	saveButton = widget.NewButton("Save", func() {
		confirmDialog := dialog.NewConfirm(
			"Save confirmation",
			"\n\n\nThis will override the existing save file. Are you sure?\nPlease make a backup if needed.",
			func(r bool) {
				if !r {
					return
				}

				cancel := animateTop(false, nil)

				applyChanges()

				resetUI()

				err := filesManager.Save()
				if err != nil {
					cancel()
					dialog.ShowError(err, w)
				}
			}, w)

		confirmDialog.SetConfirmText("Save")
		confirmDialog.SetDismissText("Cancel")
		confirmDialog.Show()
	})
	saveButton.Disable()

	tabsThemed = container.NewThemeOverride(layoutTabs, ui.TabTheme{})
	content := container.NewBorder(
		container.NewBorder(nil, nil, leftAquila, rightAquila,
			container.NewVBox(openButton, saveButton, progressLine)),
		widget.NewLabelWithData(status),
		nil, nil,
		back,
		eyeGlowOverlay,
		tabsThemed,
	)

	validateScale()

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

	sel.SetResource(ui.GetIconByName(item.ID))
	sel.SetToolTip(item.Description)
	sel.SetOptionToolTip(func(opt string) string {
		return spec.lookup(opt).Description
	})
	sel.SetOptionIcon(func(opt string) fyne.Resource {
		return ui.GetIconByName(spec.lookup(opt).ID)
	})

	sel.OnChanged(func(newVal string) {
		changed := spec.lookup(newVal)
		sel.SetResource(ui.GetIconByName(changed.ID))
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
	for unit, augmetics := range augmeticsUnits {
		featuresManager.ChangeUnitAugmetics(unit, augmetics[1])
	}
	for unit, talents := range talentsUnits {
		featuresManager.ChangeUnitTalents(unit, talents[1])
	}
}

func validateScale() {
	if runtime.GOOS == "windows" {
		return
	}

	cmd := exec.Command("xdpyinfo")
	out, err := cmd.Output()
	if err != nil {
		return
	}

	re := regexp.MustCompile(`resolution:\s+(\d+)x`)
	match := re.FindStringSubmatch(string(out))
	if len(match) == 2 {
		if dpi, _ := strconv.Atoi(match[1]); dpi > 96 {
			os.Setenv("FYNE_SCALE", "2.0")
		}
	}
}

func boolSwitch(flag *bool, icon, name, tooltip string) *toggle.Widget {
	return toggle.New(func(on bool) {
		*flag = on
		refreshSaveButton()
	}, icon, name, tooltip)
}

func resetSwitch(sw *toggle.Widget, status func() (bool, bool)) {
	sw.Enable()
	sw.SetState(false, true)
	if available, state := status(); !available {
		sw.Disable()
		sw.SetState(state, false)
	}
}

func resetSwitchOn(sw *toggle.Widget, available bool) {
	resetSwitch(sw, func() (bool, bool) { return available, true })
}
