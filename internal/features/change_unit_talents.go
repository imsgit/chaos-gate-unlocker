package features

import (
	"chaos-gate-unlocker/internal/objects"
)

var (
	talentsByName = map[string]string{
		"Aegis Adept":          "Talent_Guerilla",
		"Battle Prodigy":       "Talent_BattleProdigy",
		"Blademaster":          "Talent_Blademaster",
		"Crack Shot":           "Talent_CrackShot",
		"Cultbane":             "Talent_Cultbane",
		"Daemonbane":           "Talent_Daemonbane",
		"Deathless":            "Talent_Deathless",
		"Devoted Practitioner": "Talent_DevotedPractitioner",
		"Duelist":              "Talent_Duelist",
		"Eagle Eye":            "Talent_EagleEye",
		"Enginebane":           "Talent_Enginebane",
		"Farseer":              "Talent_Farseer",
		"Fast Recovery":        "Talent_FastRecovery",
		"Great Destiny":        "Talent_GreatDestiny",
		"Indomitable":          "Talent_Indomitable",
		"Lightning Reflexes":   "Talent_LightningReflexes",
		"Omnissiah's Chosen":   "Talent_OmnissiahsChosen",
		"Provident":            "Talent_Provident",
		"Quartermaster":        "Talent_Quartermaster",
		"Resilient":            "Talent_Resilient",
		"Skull Keeper":         "Talent_SkullKeeper",
		"Sure Strike":          "Talent_SureStrike",
		"Throwing Arm":         "Talent_ThrowingArm",
		"Venerable Soul":       "Talent_VenerableSoul",
		"Zealous Scholar":      "Talent_ZealousScholar",
		"Undying Apothecary":   "Talent_UndyingApothecary",
		"Undying Chaplain":     "Talent_UndyingChaplain",
		"Undying Interceptor":  "Talent_UndyingInterceptor",
		"Undying Justicar":     "Talent_UndyingJusticar",
		"Undying Librarian":    "Talent_UndyingLibrarian",
		"Undying Paladin":      "Talent_UndyingPaladin",
		"Undying Purgator":     "Talent_UndyingPurgator",
		"Undying Purifier":     "Talent_UndyingPurifier",
		"Undying Techmarine":   "Talent_UndyingTechMarine",
	}

	talentsByID = map[string]Talent{
		"Talent_Guerilla": {
			ID:          "Talent_Guerilla",
			Name:        "Aegis Adept",
			Description: "This Knight gains an additional +2 Armour when using their Aegis Shield Ability",
		},
		"Talent_BattleProdigy": {
			ID:          "Talent_BattleProdigy",
			Name:        "Battle Prodigy",
			Description: "This Knight gains +50% XP Rewards, but only starts with 1 Ability Point at Rank 1",
		},
		"Talent_Blademaster": {
			ID:          "Talent_Blademaster",
			Name:        "Blademaster",
			Description: "This Knight gains +5% CRIT for all their Melee Attacks",
		},
		"Talent_CrackShot": {
			ID:          "Talent_CrackShot",
			Name:        "Crack Shot",
			Description: "This Knight gains +2 CRIT DMG for all their Ranged Attacks",
		},
		"Talent_Cultbane": {
			ID:          "Talent_Cultbane",
			Name:        "Cultbane",
			Description: "The Knight has +10% CRIT for all Melee and Ranged Attacks against Organic enemies",
		},
		"Talent_Daemonbane": {
			ID:          "Talent_Daemonbane",
			Name:        "Daemonbane",
			Description: "The Knight has +10% CRIT for all Melee and Ranged Attacks against Daemonic enemies",
		},
		"Talent_Deathless": {
			ID:          "Talent_Deathless",
			Name:        "Deathless",
			Description: "This Knight cannot permanently die. He will not lose Resilience when he suffers a Critical Wound. He starts with -2 Max WP",
		},
		"Talent_DevotedPractitioner": {
			ID:          "Talent_DevotedPractitioner",
			Name:        "Devoted Practitioner",
			Description: "This Knight gains +1 Ability Point at Ranks 3, 6 and 9, but earns -25% XP Rewards",
		},
		"Talent_Duelist": {
			ID:          "Talent_Duelist",
			Name:        "Duelist",
			Description: "This Knight gains +10% Afflict Chance for their Melee Attacks",
		},
		"Talent_EagleEye": {
			ID:          "Talent_EagleEye",
			Name:        "Eagle Eye",
			Description: "This Knight gains +10% CRIT for all their Ranged Attacks",
		},
		"Talent_Enginebane": {
			ID:          "Talent_Enginebane",
			Name:        "Enginebane",
			Description: "The Knight has +10% CRIT for all Melee and Ranged Attacks against Mechanical enemies",
		},
		"Talent_Farseer": {
			ID:          "Talent_Farseer",
			Name:        "Farseer",
			Description: "This Knight gains +3 Range for Ranged Attacks with a Storm Bolter",
		},
		"Talent_FastRecovery": {
			ID:          "Talent_FastRecovery",
			Name:        "Fast Recovery",
			Description: "This Knight gains +2 HP whenever they are targeted with Heal",
		},
		"Talent_GreatDestiny": {
			ID:          "Talent_GreatDestiny",
			Name:        "Great Destiny",
			Description: "This Knight gains +2 Max HP, +1 Max WP, and +1 DMG for his Melee Attacks, he cannot gain Resilience",
		},
		"Talent_Indomitable": {
			ID:          "Talent_Indomitable",
			Name:        "Indomitable",
			Description: "This Knight starts with +2 Max HP at Rank 1",
		},
		"Talent_LightningReflexes": {
			ID:          "Talent_LightningReflexes",
			Name:        "Lightning Reflexes",
			Description: "This Knight gains +10% Focus. Focus increases their chance to trigger Autos and Afflictions",
		},
		"Talent_OmnissiahsChosen": {
			ID:          "Talent_OmnissiahsChosen",
			Name:        "Omnissiah's Chosen",
			Description: "This Knight gains double the effects from Augmetics installed due to Critical Wounds",
		},
		"Talent_Provident": {
			ID:          "Talent_Provident",
			Name:        "Provident",
			Description: "This Knight gains +1 Ammo for their equipped Ranged Weapon",
		},
		"Talent_Quartermaster": {
			ID:          "Talent_Quartermaster",
			Name:        "Quartermaster",
			Description: "This Knight gains +1 Ammo for their equipped Grenades",
		},
		"Talent_Resilient": {
			ID:          "Talent_Resilient",
			Name:        "Resilient",
			Description: "This Knight gains +10% Resistance. Resistance increases their chance to resist Afflictions",
		},
		"Talent_SkullKeeper": {
			ID:          "Talent_SkullKeeper",
			Name:        "Skull Keeper",
			Description: "This Knight gains +1 Use for their equipped Servo Skulls",
		},
		"Talent_SureStrike": {
			ID:          "Talent_SureStrike",
			Name:        "Sure Strike",
			Description: "This Knight gains +1 CRIT DMG for all their Melee Attacks",
		},
		"Talent_ThrowingArm": {
			ID:          "Talent_ThrowingArm",
			Name:        "Throwing Arm",
			Description: "This Knight gains +5 Range for their equipped Grenades",
		},
		"Talent_VenerableSoul": {
			ID:          "Talent_VenerableSoul",
			Name:        "Venerable Soul",
			Description: "This Knight starts with +1 Max WP at Rank 1",
		},
		"Talent_ZealousScholar": {
			ID:          "Talent_ZealousScholar",
			Name:        "Zealous Scholar",
			Description: "This Knight gains +25% to all XP Rewards",
		},
		"Talent_UndyingApothecary": {
			ID:          "Talent_UndyingApothecary",
			Name:        "Undying Apothecary",
			Description: "The Venerable Dreadnought's Melee Attacks collect all Bloom Seeds from the target",
		},
		"Talent_UndyingChaplain": {
			ID:          "Talent_UndyingChaplain",
			Name:        "Undying Chaplain",
			Description: "The Venerable Dreadnought gains +1 Armour",
		},
		"Talent_UndyingInterceptor": {
			ID:          "Talent_UndyingInterceptor",
			Name:        "Undying Interceptor",
			Description: "The Venerable Dreadnought gains +1 Move Speed",
		},
		"Talent_UndyingJusticar": {
			ID:          "Talent_UndyingJusticar",
			Name:        "Undying Justicar",
			Description: "The Venerable Dreadnought gains +1 DMG for his Melee Attacks",
		},
		"Talent_UndyingLibrarian": {
			ID:          "Talent_UndyingLibrarian",
			Name:        "Undying Librarian",
			Description: "The Venerable Dreadnought gains +2 Max WP",
		},
		"Talent_UndyingPaladin": {
			ID:          "Talent_UndyingPaladin",
			Name:        "Undying Paladin",
			Description: "The Venerable Dreadnought gains +4 Max HP",
		},
		"Talent_UndyingPurgator": {
			ID:          "Talent_UndyingPurgator",
			Name:        "Undying Purgator",
			Description: "The Venerable Dreadnought gains +1 Ammo for all equipped Ranged Weapons",
		},
		"Talent_UndyingPurifier": {
			ID:          "Talent_UndyingPurifier",
			Name:        "Undying Purifier",
			Description: "The Venerable Dreadnought is Immune to Hazards",
		},
		"Talent_UndyingTechMarine": {
			ID:          "Talent_UndyingTechMarine",
			Name:        "Undying Techmarine",
			Description: "The Venerable Dreadnought gains +2 HP whenever he is targeted with Heal",
		},
	}

	knightTalentsStrings = []string{
		"Aegis Adept",
		"Battle Prodigy",
		"Blademaster",
		"Crack Shot",
		"Cultbane",
		"Daemonbane",
		"Deathless",
		"Devoted Practitioner",
		"Duelist",
		"Eagle Eye",
		"Enginebane",
		"Farseer",
		"Fast Recovery",
		"Great Destiny",
		"Indomitable",
		"Lightning Reflexes",
		"Omnissiah's Chosen",
		"Provident",
		"Quartermaster",
		"Resilient",
		"Skull Keeper",
		"Sure Strike",
		"Throwing Arm",
		"Venerable Soul",
		"Zealous Scholar",
	}

	techMarineTalentsStrings = []string{
		"Battle Prodigy",
		"Blademaster",
		"Crack Shot",
		"Cultbane",
		"Daemonbane",
		"Deathless",
		"Devoted Practitioner",
		"Eagle Eye",
		"Enginebane",
		"Fast Recovery",
		"Great Destiny",
		"Indomitable",
		"Lightning Reflexes",
		"Omnissiah's Chosen",
		"Provident",
		"Resilient",
		"Sure Strike",
		"Venerable Soul",
		"Zealous Scholar",
	}

	dreadnoughtTalentsStrings = []string{
		"Undying Apothecary",
		"Undying Chaplain",
		"Undying Interceptor",
		"Undying Justicar",
		"Undying Librarian",
		"Undying Paladin",
		"Undying Purgator",
		"Undying Purifier",
		"Undying Techmarine",
	}
)

type Talent struct {
	ID          string
	Name        string
	Description string
}

func (m *Manager) ChangeUnitTalents(unit any, changedTalents []string) {
	switch object := unit.(type) {
	case *objects.KnightState:
		object.Talents = object.Talents[:0]
		for _, talent := range changedTalents {
			if talent != "" {
				object.Talents = append(object.Talents, &objects.StringValue{Key: talentsByName[talent]})
			}
		}
	case *objects.DreadnoughtState:
		object.Talents = object.Talents[:0]
		for _, talent := range changedTalents {
			if talent != "" {
				object.Talents = append(object.Talents, &objects.StringValue{Key: talentsByName[talent]})
			}
		}
	}
}

func (m *Manager) CanChangeUnitTalents(unit any, idx int) (bool, Talent, []string) {
	switch object := unit.(type) {
	case *objects.KnightState:
		var curr Talent
		for i, talent := range object.Talents {
			if idx == i {
				curr = talentsByID[talent.Key]
			}
		}

		talents := knightTalentsStrings
		class := getClass(object.CurrentLevelData.Key)
		if class == TechmarineClass {
			talents = techMarineTalentsStrings
		}

		return class != GarranCrowClass && (idx == 0 || len(object.Talents) > idx),
			curr, talents
	case *objects.DreadnoughtState:
		var curr Talent
		for i, talent := range object.Talents {
			if idx == i {
				curr = talentsByID[talent.Key]
			}
		}

		return object.HasPilot && (idx == 0 || len(object.Talents) > idx),
			curr, dreadnoughtTalentsStrings
	}

	return false, Talent{}, nil
}

func (m *Manager) TalentByName(name string) Talent {
	return talentsByID[talentsByName[name]]
}
