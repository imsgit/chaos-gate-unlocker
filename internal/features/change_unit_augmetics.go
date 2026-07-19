package features

import (
	"chaos-gate-unlocker/internal/objects"
)

type Augmetic struct {
	ID          string
	Name        string
	Description string
}

var augmetics = []Augmetic{
	{"Augmetic_AugmeticEye", "(Head) Augmetic Eye", "This Knight gains +5% CRIT for their Ranged Attacks"},
	{"Augmetic_CerebralImplant", "(Head) Cerebral Implant", "This Knight earns +25% XP Rewards"},
	{"Augmetic_CortexImplant", "(Head) Cortex Implant", "This Knight gains +5% Focus"},
	{"Augmetic_Psybooster", "(Head) Psy-Booster", "This Knight gains +1 Max WP"},
	{"Augmetic_AugmeticHand", "(Left arm) Augmetic Hand", "This Knight gains +5% CRIT for their Ranged Attacks"},
	{"Augmetic_ElbowActuator", "(Left arm) Elbow Actuator", "This Knight gains +1 Range when making targeted Ranged Attacks"},
	{"Augmetic_MuscleCasing", "(Left arm) Muscle Casing", "This Knight gains +1 CRIT DMG for their Melee Attacks"},
	{"Augmetic_AugmeticElbow", "(Right arm) Augmetic Elbow", "This Knight gains +1 CRIT DMG for their Melee Attacks"},
	{"Augmetic_Locomotion", "(Right arm) Locomotion Augmetics", "This Knight gains +5% Resistance"},
	{"Augmetic_Synthmuscle", "(Right arm) Synthmuscle", "This Knight gains +5% CRIT for their Melee Attacks"},
	{"Augmetic_SubskinBodyArmour", "(Torso) Armour Reinforcement", "This Knight gains +1 Armour"},
	{"Augmetic_AugmeticHeart", "(Torso) Augmetic Heart", "This Knight gains +1 Max WP"},
	{"Augmetic_Autosanguine_Event", "(Torso) Autosanguine", "This Knight gains +2 Max HP and +1 Armour"},
	{"Augmetic_RespiraryFilterImplant", "(Torso) Respirary Filter", "This Knight gains +5% Resistance"},
	{"Augmetic_EnhancedKneeJoint", "(Left leg) Enhanced Knee Joint", "This Knight gains +5% Focus"},
	{"Augmetic_SubskinLegArmourLeft", "(Left leg) Subskin Leg Armour", "This Knight gains +2 Max HP"},
	{"Augmetic_AugmeticFoot", "(Right leg) Augmetic Foot", "This Knight gains +5% Focus"},
	{"Augmetic_SubskinLegArmourRight", "(Right leg) Subskin Leg Armour", "This Knight gains +2 Max HP"},
}

var augmeticsStrings, augmeticsByName, augmeticsByID = func() ([]string, map[string]string, map[string]Augmetic) {
	names := make([]string, len(augmetics))
	byName := make(map[string]string, len(augmetics))
	byID := make(map[string]Augmetic, len(augmetics))
	for i, a := range augmetics {
		names[i] = a.Name
		byName[a.Name] = a.ID
		byID[a.ID] = a
	}
	return names, byName, byID
}()

const skipOption = "(Torso) Autosanguine"

func IsSkipOption(option string) bool {
	return option == skipOption
}

func (m *Manager) ChangeUnitAugmetics(unit any, changedAugmetics []string) {
	switch object := unit.(type) {
	case *objects.KnightState:
		object.Augmetics = object.Augmetics[:0]
		for _, augmetic := range changedAugmetics {
			if augmetic != "" {
				object.Augmetics = append(object.Augmetics, &objects.StringValue{Key: augmeticsByName[augmetic]})
			}
		}
	}
}

func (m *Manager) UnitSupportsAugmetics(unit any) bool {
	if object, ok := unit.(*objects.KnightState); ok {
		return getClass(object.CurrentLevelData.Key) != GarranCrowClass
	}
	return false
}

func (m *Manager) CanChangeUnitAugmetics(unit any, idx int, heal bool) (bool, Augmetic, []string) {
	switch object := unit.(type) {
	case *objects.KnightState:
		var curr Augmetic
		if idx >= 0 && idx < len(object.Augmetics) {
			curr = augmeticsByID[object.Augmetics[idx].Key]
		}

		class := getClass(object.CurrentLevelData.Key)
		return class != GarranCrowClass &&
				((object.LostResilience > idx && (object.HealthState.Status < 3 || heal)) ||
					len(object.Augmetics) > idx),
			curr, augmeticsStrings
	}

	return false, Augmetic{}, nil
}

func (m *Manager) AugmeticByName(name string) Augmetic {
	return augmeticsByID[augmeticsByName[name]]
}
