package ui

import (
	"fyne.io/fyne/v2"
)

var icons map[string]fyne.Resource

func GetIconByName(name string) fyne.Resource {
	if icons == nil {
		icons = map[string]fyne.Resource{
			"Apothecary":        resourceClassApothecaryPng,
			"CallidusAssassin":  resourceClassCallidusAssassinPng,
			"Chaplain":          resourceClassChaplainPng,
			"CulexusAssassin":   resourceClassCulexusAssassinPng,
			"Dreadnought":       resourceClassDreadnoughtPng,
			"EversorAssassin":   resourceClassEversorAssassinPng,
			"GarranCrowe":       resourceClassGarranCrowePng,
			"Interceptor":       resourceClassInterceptorPng,
			"Justicar":          resourceClassJusticarPng,
			"Librarian":         resourceClassLibrarianPng,
			"Paladin":           resourceClassPaladinPng,
			"Purgator":          resourceClassPurgatorPng,
			"Purifier":          resourceClassPurifierPng,
			"TechMarine":        resourceClassTechMarinePng,
			"VindicareAssassin": resourceClassVindicareAssassinPng,

			"ActDread":         resourceActDreadPng,
			"ActComplete":      resourceActCompletePng,
			"ActFrigate":       resourceActFrigatePng,
			"ActPrognosticars": resourceActPrognosticarsPng,
			"ActSeals":         resourceActSealsPng,
			"ActUnequip":       resourceActUnequipPng,
			"ActMarketing":     resourceActMarketingPng,
			"ActHeal":          resourceActHealPng,
			"ActRepair":        resourceActRepairPng,
			"ActAssassins":     resourceActAssassinsPng,
			"ActPreorder":      resourceActPreorderPng,

			"WidgetUnitLevel":      resourceWidgetUnitLevelPng,
			"WidgetSwitchOn":       resourceWidgetSwitchOnPng,
			"WidgetSwitchOff":      resourceWidgetSwitchOffPng,
			"WidgetSwitchOnHover":  resourceWidgetSwitchOnHoverPng,
			"WidgetSwitchOffHover": resourceWidgetSwitchOffHoverPng,

			"AppTabMain":     resourceAppTabMainPng,
			"AppTabUnits":    resourceAppTabUnitsPng,
			"AppTabAbout":    resourceAppTabAboutPng,
			"AppLeftAquila":  resourceAppLeftAquilaPng,
			"AppRightAquila": resourceAppRightAquilaPng,
			"AppBackground":  resourceAppBackgroundPng,

			"Talent_BattleProdigy":       resourceTalentBattleProdigyPng,
			"Talent_Blademaster":         resourceTalentBlademasterPng,
			"Talent_CrackShot":           resourceTalentCrackShotPng,
			"Talent_Cultbane":            resourceTalentCultbanePng,
			"Talent_Daemonbane":          resourceTalentDaemonbanePng,
			"Talent_Deathless":           resourceTalentDeathlessPng,
			"Talent_DevotedPractitioner": resourceTalentDevotedPractitionerPng,
			"Talent_Duelist":             resourceTalentDuelistPng,
			"Talent_EagleEye":            resourceTalentEagleEyePng,
			"Talent_Enginebane":          resourceTalentEnginebanePng,
			"Talent_Farseer":             resourceTalentFarseerPng,
			"Talent_FastRecovery":        resourceTalentFastRecoveryPng,
			"Talent_GreatDestiny":        resourceTalentGreatDestinyPng,
			"Talent_Guerilla":            resourceTalentGuerillaPng,
			"Talent_Indomitable":         resourceTalentIndomitablePng,
			"Talent_LightningReflexes":   resourceTalentLightningReflexesPng,
			"Talent_OmnissiahsChosen":    resourceTalentOmnissiahsChosenPng,
			"Talent_Provident":           resourceTalentProvidentPng,
			"Talent_Quartermaster":       resourceTalentQuartermasterPng,
			"Talent_Resilient":           resourceTalentResilientPng,
			"Talent_SkullKeeper":         resourceTalentSkullKeeperPng,
			"Talent_SureStrike":          resourceTalentSureStrikePng,
			"Talent_ThrowingArm":         resourceTalentThrowingArmPng,
			"Talent_UndyingApothecary":   resourceTalentUndyingApothecaryPng,
			"Talent_UndyingChaplain":     resourceTalentUndyingChaplainPng,
			"Talent_UndyingInterceptor":  resourceTalentUndyingInterceptorPng,
			"Talent_UndyingJusticar":     resourceTalentUndyingJusticarPng,
			"Talent_UndyingLibrarian":    resourceTalentUndyingLibrarianPng,
			"Talent_UndyingPaladin":      resourceTalentUndyingPaladinPng,
			"Talent_UndyingPurgator":     resourceTalentUndyingPurgatorPng,
			"Talent_UndyingPurifier":     resourceTalentUndyingPurifierPng,
			"Talent_UndyingTechMarine":   resourceTalentUndyingTechMarinePng,
			"Talent_VenerableSoul":       resourceTalentVenerableSoulPng,
			"Talent_ZealousScholar":      resourceTalentZealousScholarPng,

			"Augmetic_Autosanguine_Event":     resourceAugmeticAutosanguineEventPng,
			"Augmetic_CerebralImplant":        resourceAugmeticCerebralImplantPng,
			"Augmetic_CortexImplant":          resourceAugmeticCortexImplantPng,
			"Augmetic_AugmeticElbow":          resourceAugmeticAugmeticElbowPng,
			"Augmetic_ElbowActuator":          resourceAugmeticElbowActuatorPng,
			"Augmetic_EnhancedKneeJoint":      resourceAugmeticEnhancedKneeJointPng,
			"Augmetic_AugmeticEye":            resourceAugmeticAugmeticEyePng,
			"Augmetic_AugmeticFoot":           resourceAugmeticAugmeticFootPng,
			"Augmetic_AugmeticHand":           resourceAugmeticAugmeticHandPng,
			"Augmetic_AugmeticHeart":          resourceAugmeticAugmeticHeartPng,
			"Augmetic_Locomotion":             resourceAugmeticLocomotionPng,
			"Augmetic_MuscleCasing":           resourceAugmeticMuscleCasingPng,
			"Augmetic_Psybooster":             resourceAugmeticPsyboosterPng,
			"Augmetic_RespiraryFilterImplant": resourceAugmeticRespiraryFilterImplantPng,
			"Augmetic_SubskinBodyArmour":      resourceAugmeticSubskinBodyArmourPng,
			"Augmetic_SubskinLegArmourLeft":   resourceAugmeticSubskinLegArmourLeftPng,
			"Augmetic_SubskinLegArmourRight":  resourceAugmeticSubskinLegArmourRightPng,
			"Augmetic_Synthmuscle":            resourceAugmeticSynthmusclePng,
		}
	}

	icon, exists := icons[name]
	if !exists {
		return nil
	}
	return icon
}

func GetWidgetUnitLevelIcon() fyne.Resource      { return GetIconByName("WidgetUnitLevel") }
func GetWidgetSwitchOnIcon() fyne.Resource       { return GetIconByName("WidgetSwitchOn") }
func GetWidgetSwitchOffIcon() fyne.Resource      { return GetIconByName("WidgetSwitchOff") }
func GetWidgetSwitchOnHoverIcon() fyne.Resource  { return GetIconByName("WidgetSwitchOnHover") }
func GetWidgetSwitchOffHoverIcon() fyne.Resource { return GetIconByName("WidgetSwitchOffHover") }

func GetAppLeftAquilaIcon() fyne.Resource  { return GetIconByName("AppLeftAquila") }
func GetAppRightAquilaIcon() fyne.Resource { return GetIconByName("AppRightAquila") }
func GetAppTabMainIcon() fyne.Resource     { return GetIconByName("AppTabMain") }
func GetAppTabUnitsIcon() fyne.Resource    { return GetIconByName("AppTabUnits") }
func GetAppTabAboutIcon() fyne.Resource    { return GetIconByName("AppTabAbout") }
func GetAppBackgroundIcon() fyne.Resource  { return GetIconByName("AppBackground") }
