package main

import (
	"chaos-gate-unlocker/internal/features"
	"chaos-gate-unlocker/internal/files"
	"chaos-gate-unlocker/internal/ui"
	"chaos-gate-unlocker/internal/ui/widgets"

	"context"
	"reflect"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	fynetooltip "github.com/dweymouth/fyne-tooltip"
)

const (
	version = "Version: 1.0.0.46 | Author: imsgit | 2025-04-20"
)

var (
	featuresManager = features.NewManager()
	filesManager    = files.NewManager()

	removeMarketingWeapons       bool
	unlockPreorderItems          bool
	unlockAdvancedClasses        bool
	unlockPuritySeals            bool
	unlockAssassins              bool
	reattunePrognosticars        bool
	unlockGarranCrowe            bool
	authorizeDreadnoughtMissions bool
	repairDreadnought            bool
	unlockGladiusFrigate         bool
	completeCurrentResearch      bool
	completeCurrentConstruction  bool
	unequipMastercraftedWeapons  bool
	unequipMastercraftedArmor    bool

	currUnit       any
	healUnits      = map[any]bool{}
	talentsUnits   = map[any][][]string{}
	augmeticsUnits = map[any][][]string{}

	saveButton *widget.Button
)

func main() {
	a := app.NewWithID("chaos.gate.unlocker")
	a.Settings().SetTheme(ui.Theme{})
	w := a.NewWindow("Chaos Gate Unlocker")

	filesManager.OnLoadState(featuresManager.ApplyState())

	authorizeDreadnoughtMissionsSwitch := widgets.NewSwitch(func(on bool) {
		authorizeDreadnoughtMissions = on
		refreshSaveButton()
	}, "ActDread", "Authorize Dreadnought missions", "Marks all regular missions as Technophage, including Hive missions;\nThe difficulty of the missions will increase, but it won't affect the frigate's chances of winning;\nDreadnought access is required")

	var repairDamageSwitch *widgets.Switch
	repairDreadnoughtSwitch := widgets.NewSwitch(func(on bool) {
		repairDreadnought = on
		refreshSaveButton()
		if repairDamageSwitch.Visible() {
			repairDamageSwitch.SetState(on, false)
		}
	}, "ActRepair", "Repair Dreadnought", "Repairs the Dreadnought's damage for free;\nDreadnought access is required")

	reattunePrognosticarsSwitch := widgets.NewSwitch(func(on bool) {
		reattunePrognosticars = on
		refreshSaveButton()
	}, "ActPrognosticars", "Reattune prognosticars", "Restores all attuned prognosticars, making them available again")

	completeCurrentResearchSwitch := widgets.NewSwitch(func(on bool) {
		completeCurrentResearch = on
		refreshSaveButton()
	}, "ActComplete", "Complete current research", "Completes current research project;\nAdvance current day to unlock")

	completeCurrentConstructionSwitch := widgets.NewSwitch(func(on bool) {
		completeCurrentConstruction = on
		refreshSaveButton()
	}, "ActComplete", "Complete current construction", "Completes current construction project;\nAdvance current day to unlock")

	unequipMastercraftedArmorSwitch := widgets.NewSwitch(func(on bool) {
		unequipMastercraftedArmor = on
		refreshSaveButton()
	}, "ActUnequip", "Unequip mastercrafted armor", "Unequips all mastercrafted armor, clearing the slots if necessary;\nAlso unlocks currently unavailable armor, but it won't affect the frigate's chances of winning")

	unequipMastercraftedWeaponsSwitch := widgets.NewSwitch(func(on bool) {
		unequipMastercraftedWeapons = on
		refreshSaveButton()
	}, "ActUnequip", "Unequip mastercrafted weapons", "Unequips all mastercrafted weapons;\nAlso unlocks currently unavailable weapons, but it won't affect the frigate's chances of winning")

	unlockPreorderItemsSwitch := widgets.NewSwitch(func(on bool) {
		unlockPreorderItems = on
		refreshSaveButton()
	}, "ActPreorder", "Unlock pre-order items", "Unlocks the Domina Liber Daemonica tome and Destroyer of Crys'yllix hammer")

	unlockAdvancedClassesSwitch := widgets.NewSwitch(func(on bool) {
		unlockAdvancedClasses = on
		refreshSaveButton()
	}, "Librarian", "Unlock advanced classes", "Unlocks the Librarian, Paladin, Chaplain and Purifier classes;\nAdvance current day to unlock")

	unlockGarranCroweSwitch := widgets.NewSwitch(func(on bool) {
		unlockGarranCrowe = on
		refreshSaveButton()
	}, "GarranCrowe", "Unlock Garran Crowe", "Unlocks castellan Garran Crowe;\nDLC access is required;\nAdvance current day to unlock")

	unlockAssassinsSwitch := widgets.NewSwitch(func(on bool) {
		unlockAssassins = on
		refreshSaveButton()
	}, "ActAssassins", "Unlock assassins", "Unlocks imperial assassins;\nDLC access is required;\nAdvance current day to unlock")

	unlockGladiusFrigateSwitch := widgets.NewSwitch(func(on bool) {
		unlockGladiusFrigate = on
		refreshSaveButton()
	}, "ActFrigate", "Unlock Gladius frigate", "Unlocks the Gladius frigate, but the Cleanse mission will still appear as expected;\nDLC access is required;\nAdvance current day to unlock")

	unlockPuritySealsSwitch := widgets.NewSwitch(func(on bool) {
		unlockPuritySeals = on
		refreshSaveButton()
	}, "ActSeals", "Unlock purity seals", "Unlocks purity seals upgrades;\nPoxus seeds access is required;\nAdvance current day to unlock")

	removeMarketingWeaponsSwitch := widgets.NewSwitch(func(on bool) {
		removeMarketingWeapons = on
		refreshSaveButton()
	}, "ActMarketing", "Remove marketing weapons", "Unequips and removes all weapons classified as Twitch drops")

	unitsBox := container.NewVBox()
	unitsScrollBox := container.NewVScroll(unitsBox)

	healWoundSwitch := widgets.NewSwitch(func(on bool) {
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
	}, "ActHeal", "Heal wound", "Heals the wound;\nIf the wound was critical, you can also select a new augmetic for your knight")
	healWoundSwitch.Hide()

	repairDamageSwitch = widgets.NewSwitch(func(on bool) {
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
		widgets.NewListItem,
		func(item binding.DataItem, o fyne.CanvasObject) {
			f := reflect.ValueOf(item).Elem().FieldByName("index")
			if f.IsValid() && f.CanInt() {
				val, _ := unitsProvider.GetValue(int(f.Int()))
				listItem, _ := o.(*widgets.ListItem)
				listItem.Bind(val)
			}
		})
	unitsList.HideSeparators = true

	unitsList.OnSelected = func(id widget.ListItemID) {
		currUnit, _ = unitsProvider.GetValue(id)

		enable, showHeal := featuresManager.CanHealUnit(currUnit)
		if showHeal {
			healWoundSwitch.Show()
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
		talentsBox := container.NewVBox()
		unitsBox.Objects[2] = talentsBox
		for i := 0; ; i++ {
			if sel := renderTalent(i, initTalents); sel != nil {
				talentsBox.Objects = append(talentsBox.Objects, sel)
				continue
			}
			break
		}

		initAugmetics := len(augmeticsUnits[currUnit]) == 0
		if initAugmetics {
			augmeticsUnits[currUnit] = append(augmeticsUnits[currUnit], []string{}, []string{})
		}
		augmeticsBox := container.NewVBox()
		unitsBox.Objects[3] = augmeticsBox
		for i := 0; ; i++ {
			if sel := renderAugmetic(i, initAugmetics); sel != nil {
				augmeticsBox.Objects = append(augmeticsBox.Objects, sel)
				continue
			}
			break
		}
	}

	unitsList.OnUnselected = func(_ widget.ListItemID) {
		healWoundSwitch.Hide()
		repairDamageSwitch.Hide()
		unitsBox.Objects[2] = container.NewVBox()
		unitsBox.Objects[3] = container.NewVBox()
		unitsScrollBox.ScrollToTop()
	}

	var cancel context.CancelFunc
	back := canvas.NewImageFromResource(ui.GetAppBackgroundIcon())
	back.FillMode = canvas.ImageFillContain
	back.Translucency = 0.96

	mainTab := container.NewTabItemWithIcon("Main", ui.GetAppTabMainIcon(),
		container.NewGridWithColumns(2,
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
				removeMarketingWeaponsSwitch)))
	unitsTab := container.NewTabItemWithIcon("Units", ui.GetAppTabUnitsIcon(),
		container.NewGridWithColumns(2, unitsList, unitsScrollBox))
	aboutTab := container.NewTabItemWithIcon("About", ui.GetAppTabAboutIcon(),
		container.NewBorder(nil, nil,
			widget.NewRichTextFromMarkdown(`
[> Visit Nexus Mods for more information](https://www.nexusmods.com/warhammer40kchaosgatedaemonhunters/mods/5)

[> Visit Reddit for discussion](https://www.reddit.com/r/ChaosGateGame/comments/1hz3s5g/chaosgateunlocker)
`),
			widget.NewRichTextFromMarkdown(version)))
	layoutTabs := container.NewAppTabs(mainTab, unitsTab, aboutTab)
	layoutTabs.SetTabLocation(container.TabLocationTrailing)
	layoutTabs.OnSelected = func(item *container.TabItem) {
		if item.Text == "About" {
			var ctx context.Context
			ctx, cancel = context.WithCancel(context.Background())
			go animateAbout(ctx, back)
		} else {
			if cancel != nil {
				cancel()
			}
			back.Translucency = 0.96
		}
	}
	layoutTabs.Hide()

	leftAquila := canvas.NewImageFromResource(ui.GetAppLeftAquilaIcon())
	leftAquila.SetMinSize(fyne.NewSize(100, 0))
	leftAquila.Translucency = 1

	rightAquila := canvas.NewImageFromResource(ui.GetAppRightAquilaIcon())
	rightAquila.SetMinSize(fyne.NewSize(100, 0))
	rightAquila.Translucency = 1

	progress := canvas.NewRectangle(fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameBackground, 0))
	progress.SetMinSize(fyne.NewSize(0, 4))

	openButton := widget.NewButton("Open", func() {
		fileDialog := dialog.NewFileOpen(func(rc fyne.URIReadCloser, err error) {
			if rc == nil {
				return
			}

			go func() {
				animateTop(leftAquila, rightAquila, progress, true)
				fyne.DoAndWait(func() {
					layoutTabs.Show()
					_ = status.Set(filesManager.Status())
				})
			}()

			healUnits = map[any]bool{}
			augmeticsUnits = map[any][][]string{}
			talentsUnits = map[any][][]string{}

			layoutTabs.Hide()
			layoutTabs.SelectIndex(0)
			unitsList.UnselectAll()
			_ = status.Set("")

			for unitsProvider.Length() > 0 {
				unit, _ := unitsProvider.GetValue(0)
				_ = unitsProvider.Remove(unit)
			}

			err = filesManager.Load(rc)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			for _, unit := range featuresManager.Units() {
				_ = unitsProvider.Append(unit)
			}

			unlockAdvancedClassesSwitch.Enable()
			unlockAdvancedClassesSwitch.SetState(false, true)
			if canUnlock, unlocked := featuresManager.CanUnlockAdvancedClasses(); !canUnlock {
				unlockAdvancedClassesSwitch.Disable()
				unlockAdvancedClassesSwitch.SetState(unlocked, false)
			}

			repairDreadnoughtSwitch.Enable()
			repairDreadnoughtSwitch.SetState(false, true)
			if canUnlock, unlocked := featuresManager.CanRepairDreadnought(); !canUnlock {
				repairDreadnoughtSwitch.Disable()
				repairDreadnoughtSwitch.SetState(unlocked, false)
			}

			unlockPuritySealsSwitch.Enable()
			unlockPuritySealsSwitch.SetState(false, true)
			if canUnlock, unlocked := featuresManager.CanUnlockPuritySeals(); !canUnlock {
				unlockPuritySealsSwitch.Disable()
				unlockPuritySealsSwitch.SetState(unlocked, false)
			}

			reattunePrognosticarsSwitch.Enable()
			reattunePrognosticarsSwitch.SetState(false, true)
			if canUnlock, unlocked := featuresManager.CanReattunePrognosticars(); !canUnlock {
				reattunePrognosticarsSwitch.Disable()
				reattunePrognosticarsSwitch.SetState(unlocked, false)
			}

			unlockGarranCroweSwitch.Enable()
			unlockGarranCroweSwitch.SetState(false, true)
			if canUnlock, unlocked := featuresManager.CanUnlockGarranCrowe(); !canUnlock {
				unlockGarranCroweSwitch.Disable()
				unlockGarranCroweSwitch.SetState(unlocked, false)
			}

			unlockPreorderItemsSwitch.Enable()
			unlockPreorderItemsSwitch.SetState(false, true)
			if canUnlock := featuresManager.CanUnlockPreorderItems(); !canUnlock {
				unlockPreorderItemsSwitch.Disable()
				unlockPreorderItemsSwitch.SetState(true, false)
			}

			authorizeDreadnoughtMissionsSwitch.Enable()
			authorizeDreadnoughtMissionsSwitch.SetState(false, true)
			if canUnlock, unlocked := featuresManager.CanAuthorizeDreadnoughtMissions(); !canUnlock {
				authorizeDreadnoughtMissionsSwitch.Disable()
				authorizeDreadnoughtMissionsSwitch.SetState(unlocked, false)
			}

			unlockGladiusFrigateSwitch.Enable()
			unlockGladiusFrigateSwitch.SetState(false, true)
			if canUnlock, unlocked := featuresManager.CanUnlockGladiusFrigate(); !canUnlock {
				unlockGladiusFrigateSwitch.Disable()
				unlockGladiusFrigateSwitch.SetState(unlocked, false)
			}

			completeCurrentResearchSwitch.Enable()
			completeCurrentResearchSwitch.SetState(false, true)
			if canUnlock, unlocked := featuresManager.CanCompleteCurrentResearch(); !canUnlock {
				completeCurrentResearchSwitch.Disable()
				completeCurrentResearchSwitch.SetState(unlocked, false)
			}

			completeCurrentConstructionSwitch.Enable()
			completeCurrentConstructionSwitch.SetState(false, true)
			if canUnlock, unlocked := featuresManager.CanCompleteCurrentConstruction(); !canUnlock {
				completeCurrentConstructionSwitch.Disable()
				completeCurrentConstructionSwitch.SetState(unlocked, false)
			}

			unequipMastercraftedWeaponsSwitch.Enable()
			unequipMastercraftedWeaponsSwitch.SetState(false, true)
			if !featuresManager.CanUnequipMastercraftedWeapons() {
				unequipMastercraftedWeaponsSwitch.Disable()
				unequipMastercraftedWeaponsSwitch.SetState(true, false)
			}

			unequipMastercraftedArmorSwitch.Enable()
			unequipMastercraftedArmorSwitch.SetState(false, true)
			if !featuresManager.CanUnequipMastercraftedArmor() {
				unequipMastercraftedArmorSwitch.Disable()
				unequipMastercraftedArmorSwitch.SetState(true, false)
			}

			removeMarketingWeaponsSwitch.Enable()
			removeMarketingWeaponsSwitch.SetState(false, true)
			if !featuresManager.CanRemoveMarketingWeapons() {
				removeMarketingWeaponsSwitch.Disable()
				removeMarketingWeaponsSwitch.SetState(true, false)
			}

			unlockAssassinsSwitch.Enable()
			unlockAssassinsSwitch.SetState(false, true)
			if canUnlock, unlocked := featuresManager.CanUnlockAssassins(); !canUnlock {
				unlockAssassinsSwitch.Disable()
				unlockAssassinsSwitch.SetState(unlocked, false)
			}
		}, w)

		l, _ := storage.ListerForURI(storage.NewFileURI(filesManager.GetCurrentPath()))
		fileDialog.SetLocation(l)
		fileDialog.Resize(fyne.NewSize(800, 600))
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".gksave"}))
		fileDialog.Show()
	})

	saveButton = widget.NewButton("Save", func() {
		dialog.ShowConfirm(
			"Save confirmation",
			"\n\nThis will override the existing save file. Are you sure?\nPlease make a backup if needed.\n\n",
			func(response bool) {
				if response {
					go animateTop(leftAquila, rightAquila, progress, false)

					applyChanges()

					saveButton.Disable()
					layoutTabs.Hide()
					layoutTabs.SelectIndex(0)
					unitsList.UnselectAll()
					_ = status.Set("")

					err := filesManager.Save()
					if err != nil {
						dialog.ShowError(err, w)
					}
				}
			}, w)
	})
	saveButton.Disable()

	content := container.NewBorder(
		container.NewBorder(nil, nil, leftAquila, rightAquila,
			container.NewVBox(openButton, saveButton, progress)),
		widget.NewLabelWithData(status),
		nil, nil,
		back,
		layoutTabs,
	)

	w.Resize(fyne.NewSize(800, 600))
	w.SetContent(fynetooltip.AddWindowToolTipLayer(content, w.Canvas()))
	w.CenterOnScreen()
	w.ShowAndRun()
}

func renderTalent(idx int, init bool) *widgets.SelectIcon {
	if canChange, talent, options := featuresManager.CanChangeUnitTalents(currUnit, idx); canChange {
		sel := widgets.NewSelectIcon()
		sel.SetPlaceHolder("(Select talent)")
		sel.SetOptions(options)

		if init {
			sel.SetSelected(talent.Name)
			talentsUnits[currUnit][0] = append(talentsUnits[currUnit][0], talent.Name)
			talentsUnits[currUnit][1] = append(talentsUnits[currUnit][1], talent.Name)
		} else {
			sel.SetSelected(talentsUnits[currUnit][1][idx])
			talent = featuresManager.TalentByName(sel.Selected())
		}

		sel.SetResource(ui.GetIconByName(talent.ID))
		sel.SetToolTip(talent.Description)

		sel.OnChanged(func(newVal string) {
			sel.SetResource(ui.GetIconByName(featuresManager.TalentByName(newVal).ID))
			sel.SetToolTip(featuresManager.TalentByName(newVal).Description)
			talentsUnits[currUnit][1][idx] = newVal
			refreshSaveButton()
		})

		sel.OnBeforeShowPopup(func() {
			var opts []string
			for _, opt := range options {
				if !containsOpt(talentsUnits[currUnit][1], opt) || opt == sel.Selected() {
					opts = append(opts, opt)
				}
			}
			sel.SetOptions(opts)
		})
		return sel
	}
	return nil
}

func renderAugmetic(idx int, init bool) *widgets.SelectIcon {
	if canChange, augmetic, options := featuresManager.CanChangeUnitAugmetics(currUnit, idx, healUnits[currUnit]); canChange {
		sel := widgets.NewSelectIcon()
		sel.SetPlaceHolder("(Select augmetic)")
		sel.SetOptions(options)

		if init {
			sel.SetSelected(augmetic.Name)
			augmeticsUnits[currUnit][0] = append(augmeticsUnits[currUnit][0], augmetic.Name)
			augmeticsUnits[currUnit][1] = append(augmeticsUnits[currUnit][1], augmetic.Name)
		} else {
			sel.SetSelected(augmeticsUnits[currUnit][1][idx])
			augmetic = featuresManager.AugmeticByName(sel.Selected())
		}

		sel.SetResource(ui.GetIconByName(augmetic.ID))
		sel.SetToolTip(augmetic.Description)

		sel.OnChanged(func(newVal string) {
			sel.SetResource(ui.GetIconByName(featuresManager.AugmeticByName(newVal).ID))
			sel.SetToolTip(featuresManager.AugmeticByName(newVal).Description)
			augmeticsUnits[currUnit][1][idx] = newVal
			refreshSaveButton()
		})

		sel.OnBeforeShowPopup(func() {
			var opts []string
			for _, opt := range options {
				if !containsOpt(augmeticsUnits[currUnit][1], opt) || opt == sel.Selected() {
					opts = append(opts, opt)
				}
			}
			sel.SetOptions(opts)
		})
		return sel
	}
	return nil
}

const skipOpt = "(Torso) Autosanguine"

func containsOpt(list []string, val string) bool {
	for _, v := range list {
		if v == val || val == skipOpt {
			return true
		}
	}
	return false
}

func refreshSaveButton() {
	var augmeticsChanged bool
	for _, augmetics := range augmeticsUnits {
		if strings.Join(augmetics[0], ",") != strings.Join(augmetics[1], ",") {
			augmeticsChanged = true
			break
		}
	}

	var talentsChanged bool
	for _, talents := range talentsUnits {
		if strings.Join(talents[0], ",") != strings.Join(talents[1], ",") {
			talentsChanged = true
			break
		}
	}

	canApplyChanges := unlockPuritySeals ||
		unlockAdvancedClasses ||
		repairDreadnought ||
		reattunePrognosticars ||
		unlockGarranCrowe ||
		authorizeDreadnoughtMissions ||
		unlockGladiusFrigate ||
		completeCurrentConstruction ||
		completeCurrentResearch ||
		unequipMastercraftedWeapons ||
		unequipMastercraftedArmor ||
		removeMarketingWeapons ||
		unlockAssassins ||
		unlockPreorderItems ||
		len(healUnits) > 0 ||
		augmeticsChanged ||
		talentsChanged

	saveButton.Disable()
	if canApplyChanges {
		saveButton.Enable()
	}
}

func applyChanges() {
	if unlockPuritySeals {
		featuresManager.UnlockPuritySeals()
	}

	if repairDreadnought {
		featuresManager.RepairDreadnought()
	}

	if unlockPreorderItems {
		featuresManager.UnlockPreorderItems()
	}

	if unlockAdvancedClasses {
		featuresManager.UnlockAdvancedClasses()
	}

	if reattunePrognosticars {
		featuresManager.ReattunePrognosticars()
	}

	if unlockGarranCrowe {
		featuresManager.UnlockGarranCrowe()
	}

	if authorizeDreadnoughtMissions {
		featuresManager.AuthorizeDreadnoughtMissions()
	}

	if unlockGladiusFrigate {
		featuresManager.UnlockGladiusFrigate()
	}

	if completeCurrentConstruction {
		featuresManager.CompleteCurrentConstruction()
	}

	if completeCurrentResearch {
		featuresManager.CompleteCurrentResearch()
	}

	if unequipMastercraftedWeapons {
		featuresManager.UnequipMastercraftedWeapons()
	}

	if unequipMastercraftedArmor {
		featuresManager.UnequipMastercraftedArmor()
	}

	if removeMarketingWeapons {
		featuresManager.RemoveMarketingWeapons()
	}

	if unlockAssassins {
		featuresManager.UnlockAssassins()
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

func animateTop(im, im2 *canvas.Image, r *canvas.Rectangle, open bool) {
	currentSize := fyne.NewSize(0, 4)
	sOffset := r.Size().Width / 20
	tOffset := 0.04
	if open {
		tOffset *= -1
		im.Translucency = 1
		im2.Translucency = 1
	}

	fyne.DoAndWait(func() {
		r.FillColor = fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameBackground, 0)
	})

	ticker := time.NewTicker(24 * time.Millisecond)
	defer ticker.Stop()

	for i := 0; i < 30; i++ {
		select {
		case <-ticker.C:
			if i < 20 {
				newSize := fyne.NewSize(currentSize.Width+sOffset, 4)
				currentSize = newSize
				fyne.DoAndWait(func() {
					r.FillColor = fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameShadow, 0)
					r.Resize(newSize)
				})
			} else if !open {
				fyne.DoAndWait(func() {
					r.FillColor = fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameBackground, 0)
				})
			}

			if i > 4 {
				fyne.DoAndWait(func() {
					im.Translucency += tOffset
					if im.Translucency < 0 {
						im.Translucency = 0
					}
					if im.Translucency > 1 {
						im.Translucency = 1
					}
					im.Refresh()

					im2.Translucency += tOffset
					if im2.Translucency < 0 {
						im2.Translucency = 0
					}
					if im2.Translucency > 1 {
						im2.Translucency = 1
					}
					im2.Refresh()
				})
			}
		}
	}
}

func animateAbout(ctx context.Context, im *canvas.Image) {
	tOffset := -0.04

	ticker := time.NewTicker(24 * time.Millisecond)
	defer ticker.Stop()

	for i := 0; i < 30; i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if i > 5 {
				fyne.DoAndWait(func() {
					im.Translucency += tOffset
					if im.Translucency < 0 {
						im.Translucency = 0
					}
					im.Refresh()
				})
			}
		}
	}
}
