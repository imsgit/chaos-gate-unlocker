package features

import (
	"chaos-gate-unlocker/internal/objects"

	"strings"
)

type Talent struct {
	ID          string
	Name        string
	Description string
}

var talents = []Talent{
	{"Talent_Guerilla", "Aegis Adept", "This Knight gains an additional +2 Armour when using their Aegis Shield Ability"},
	{"Talent_BattleProdigy", "Battle Prodigy", "This Knight gains +50% XP Rewards, but only starts with 1 Ability Point at Rank 1"},
	{"Talent_Blademaster", "Blademaster", "This Knight gains +5% CRIT for all their Melee Attacks"},
	{"Talent_CrackShot", "Crack Shot", "This Knight gains +2 CRIT DMG for all their Ranged Attacks"},
	{"Talent_Cultbane", "Cultbane", "The Knight has +10% CRIT for all Melee and Ranged Attacks against Organic enemies"},
	{"Talent_Daemonbane", "Daemonbane", "The Knight has +10% CRIT for all Melee and Ranged Attacks against Daemonic enemies"},
	{"Talent_Deathless", "Deathless", "This Knight cannot permanently die. He will not lose Resilience when he suffers a Critical Wound. He starts with -2 Max WP"},
	{"Talent_DevotedPractitioner", "Devoted Practitioner", "This Knight gains +1 Ability Point at Ranks 3, 6 and 9, but earns -25% XP Rewards"},
	{"Talent_Duelist", "Duelist", "This Knight gains +10% Afflict Chance for their Melee Attacks"},
	{"Talent_EagleEye", "Eagle Eye", "This Knight gains +10% CRIT for all their Ranged Attacks"},
	{"Talent_Enginebane", "Enginebane", "The Knight has +10% CRIT for all Melee and Ranged Attacks against Mechanical enemies"},
	{"Talent_Farseer", "Farseer", "This Knight gains +3 Range for Ranged Attacks with a Storm Bolter"},
	{"Talent_FastRecovery", "Fast Recovery", "This Knight gains +2 HP whenever they are targeted with Heal"},
	{"Talent_GreatDestiny", "Great Destiny", "This Knight gains +2 Max HP, +1 Max WP, and +1 DMG for his Melee Attacks, he cannot gain Resilience"},
	{"Talent_Indomitable", "Indomitable", "This Knight starts with +2 Max HP at Rank 1"},
	{"Talent_LightningReflexes", "Lightning Reflexes", "This Knight gains +10% Focus. Focus increases their chance to trigger Autos and Afflictions"},
	{"Talent_OmnissiahsChosen", "Omnissiah's Chosen", "This Knight gains double the effects from Augmetics installed due to Critical Wounds"},
	{"Talent_Provident", "Provident", "This Knight gains +1 Ammo for their equipped Ranged Weapon"},
	{"Talent_Quartermaster", "Quartermaster", "This Knight gains +1 Ammo for their equipped Grenades"},
	{"Talent_Resilient", "Resilient", "This Knight gains +10% Resistance. Resistance increases their chance to resist Afflictions"},
	{"Talent_SkullKeeper", "Skull Keeper", "This Knight gains +1 Use for their equipped Servo Skulls"},
	{"Talent_SureStrike", "Sure Strike", "This Knight gains +1 CRIT DMG for all their Melee Attacks"},
	{"Talent_ThrowingArm", "Throwing Arm", "This Knight gains +5 Range for their equipped Grenades"},
	{"Talent_UndyingApothecary", "Undying Apothecary", "The Venerable Dreadnought's Melee Attacks collect all Bloom Seeds from the target"},
	{"Talent_UndyingChaplain", "Undying Chaplain", "The Venerable Dreadnought gains +1 Armour"},
	{"Talent_UndyingInterceptor", "Undying Interceptor", "The Venerable Dreadnought gains +1 Move Speed"},
	{"Talent_UndyingJusticar", "Undying Justicar", "The Venerable Dreadnought gains +1 DMG for his Melee Attacks"},
	{"Talent_UndyingLibrarian", "Undying Librarian", "The Venerable Dreadnought gains +2 Max WP"},
	{"Talent_UndyingPaladin", "Undying Paladin", "The Venerable Dreadnought gains +4 Max HP"},
	{"Talent_UndyingPurgator", "Undying Purgator", "The Venerable Dreadnought gains +1 Ammo for all equipped Ranged Weapons"},
	{"Talent_UndyingPurifier", "Undying Purifier", "The Venerable Dreadnought is Immune to Hazards"},
	{"Talent_UndyingTechMarine", "Undying Techmarine", "The Venerable Dreadnought gains +2 HP whenever he is targeted with Heal"},
	{"Talent_VenerableSoul", "Venerable Soul", "This Knight starts with +1 Max WP at Rank 1"},
	{"Talent_ZealousScholar", "Zealous Scholar", "This Knight gains +25% to all XP Rewards"},
}

var techMarineExcluded = map[string]bool{
	"Talent_Guerilla":      true,
	"Talent_Duelist":       true,
	"Talent_Farseer":       true,
	"Talent_Quartermaster": true,
	"Talent_SkullKeeper":   true,
	"Talent_ThrowingArm":   true,
}

var talentsByID, talentsByName, knightTalentsStrings, techMarineTalentsStrings, dreadnoughtTalentsStrings = func() (map[string]Talent, map[string]string, []string, []string, []string) {
	byID := make(map[string]Talent, len(talents))
	byName := make(map[string]string, len(talents))
	var knight, techMarine, dreadnought []string
	for _, t := range talents {
		byID[t.ID] = t
		byName[t.Name] = t.ID
		switch {
		case strings.HasPrefix(t.ID, "Talent_Undying"):
			dreadnought = append(dreadnought, t.Name)
		case techMarineExcluded[t.ID]:
			knight = append(knight, t.Name)
		default:
			knight = append(knight, t.Name)
			techMarine = append(techMarine, t.Name)
		}
	}
	return byID, byName, knight, techMarine, dreadnought
}()

func (m *Manager) ChangeUnitTalents(unit any, changedTalents []string) {
	switch object := unit.(type) {
	case *objects.KnightState:
		object.Talents = keysFromNames(object.Talents, changedTalents, talentsByName)
	case *objects.DreadnoughtState:
		object.Talents = keysFromNames(object.Talents, changedTalents, talentsByName)
	}
}

func (m *Manager) CanChangeUnitTalents(unit any, idx int) (bool, Talent, []string) {
	switch object := unit.(type) {
	case *objects.KnightState:
		talents := knightTalentsStrings
		class := stem(object.CurrentLevelData.Key)
		if class == TechmarineClass {
			talents = techMarineTalentsStrings
		}

		return class != GarranCrowClass && (idx == 0 || len(object.Talents) > idx),
			lookupAt(object.Talents, idx, talentsByID), talents
	case *objects.DreadnoughtState:
		return object.HasPilot && (idx == 0 || len(object.Talents) > idx),
			lookupAt(object.Talents, idx, talentsByID), dreadnoughtTalentsStrings
	}

	return false, Talent{}, nil
}

func (m *Manager) TalentByName(name string) Talent {
	return talentsByID[talentsByName[name]]
}
