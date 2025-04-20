package features

import (
	"chaos-gate-unlocker/internal/objects"
)

var (
	augmeticsByName = map[string]string{
		"(Head) Augmetic Eye":              "Augmetic_AugmeticEye",
		"(Head) Cerebral Implant":          "Augmetic_CerebralImplant",
		"(Head) Cortex Implant":            "Augmetic_CortexImplant",
		"(Head) Psy-Booster":               "Augmetic_Psybooster",
		"(Left arm) Augmetic Hand":         "Augmetic_AugmeticHand",
		"(Left arm) Elbow Actuator":        "Augmetic_ElbowActuator",
		"(Left arm) Muscle Casing":         "Augmetic_MuscleCasing",
		"(Right arm) Augmetic Elbow":       "Augmetic_AugmeticElbow",
		"(Right arm) Locomotion Augmetics": "Augmetic_Locomotion",
		"(Right arm) Synthmuscle":          "Augmetic_Synthmuscle",
		"(Torso) Armour Reinforcement":     "Augmetic_SubskinBodyArmour",
		"(Torso) Augmetic Heart":           "Augmetic_AugmeticHeart",
		"(Torso) Autosanguine":             "Augmetic_Autosanguine_Event",
		"(Torso) Respirary Filter":         "Augmetic_RespiraryFilterImplant",
		"(Left leg) Enhanced Knee Joint":   "Augmetic_EnhancedKneeJoint",
		"(Left leg) Subskin Leg Armour":    "Augmetic_SubskinLegArmourLeft",
		"(Right leg) Augmetic Foot":        "Augmetic_AugmeticFoot",
		"(Right leg) Subskin Leg Armour":   "Augmetic_SubskinLegArmourRight",
	}

	augmeticsByID = map[string]Augmetic{
		"Augmetic_AugmeticEye": {
			ID:          "Augmetic_AugmeticEye",
			Name:        "(Head) Augmetic Eye",
			Description: "This Knight gains +5% CRIT for their Ranged Attacks",
		},
		"Augmetic_CerebralImplant": {
			ID:          "Augmetic_CerebralImplant",
			Name:        "(Head) Cerebral Implant",
			Description: "This Knight earns +25% XP Rewards",
		},
		"Augmetic_CortexImplant": {
			ID:          "Augmetic_CortexImplant",
			Name:        "(Head) Cortex Implant",
			Description: "This Knight gains +5% Focus",
		},
		"Augmetic_Psybooster": {
			ID:          "Augmetic_Psybooster",
			Name:        "(Head) Psy-Booster",
			Description: "This Knight gains +1 Max WP",
		},
		"Augmetic_AugmeticHand": {
			ID:          "Augmetic_AugmeticHand",
			Name:        "(Left arm) Augmetic Hand",
			Description: "This Knight gains +5% CRIT for their Ranged Attacks",
		},
		"Augmetic_ElbowActuator": {
			ID:          "Augmetic_ElbowActuator",
			Name:        "(Left arm) Elbow Actuator",
			Description: "This Knight gains +1 Range when making targeted Ranged Attacks",
		},
		"Augmetic_MuscleCasing": {
			ID:          "Augmetic_MuscleCasing",
			Name:        "(Left arm) Muscle Casing",
			Description: "This Knight gains +1 CRIT DMG for their Melee Attacks",
		},
		"Augmetic_AugmeticElbow": {
			ID:          "Augmetic_AugmeticElbow",
			Name:        "(Right arm) Augmetic Elbow",
			Description: "This Knight gains +1 CRIT DMG for their Melee Attacks",
		},
		"Augmetic_Locomotion": {
			ID:          "Augmetic_Locomotion",
			Name:        "(Right arm) Locomotion Augmetics",
			Description: "This Knight gains +5% Resistance",
		},
		"Augmetic_Synthmuscle": {
			ID:          "Augmetic_Synthmuscle",
			Name:        "(Right arm) Synthmuscle",
			Description: "This Knight gains +5% CRIT for their Melee Attacks",
		},
		"Augmetic_SubskinBodyArmour": {
			ID:          "Augmetic_SubskinBodyArmour",
			Name:        "(Torso) Armour Reinforcement",
			Description: "This Knight gains +1 Armour",
		},
		"Augmetic_AugmeticHeart": {
			ID:          "Augmetic_AugmeticHeart",
			Name:        "(Torso) Augmetic Heart",
			Description: "This Knight gains +1 Max WP",
		},
		"Augmetic_Autosanguine_Event": {
			ID:          "Augmetic_Autosanguine_Event",
			Name:        "(Torso) Autosanguine",
			Description: "This Knight gains +2 Max HP and +1 Armour",
		},
		"Augmetic_RespiraryFilterImplant": {
			ID:          "Augmetic_RespiraryFilterImplant",
			Name:        "(Torso) Respirary Filter",
			Description: "This Knight gains +5% Resistance",
		},
		"Augmetic_EnhancedKneeJoint": {
			ID:          "Augmetic_EnhancedKneeJoint",
			Name:        "(Left leg) Enhanced Knee Joint",
			Description: "This Knight gains +5% Focus",
		},
		"Augmetic_SubskinLegArmourLeft": {
			ID:          "Augmetic_SubskinLegArmourLeft",
			Name:        "(Left leg) Subskin Leg Armour",
			Description: "This Knight gains +2 Max HP",
		},
		"Augmetic_AugmeticFoot": {
			ID:          "Augmetic_AugmeticFoot",
			Name:        "(Right leg) Augmetic Foot",
			Description: "This Knight gains +5% Focus",
		},
		"Augmetic_SubskinLegArmourRight": {
			ID:          "Augmetic_SubskinLegArmourRight",
			Name:        "(Right leg) Subskin Leg Armour",
			Description: "This Knight gains +2 Max HP",
		},
	}

	augmeticsStrings = []string{
		"(Head) Augmetic Eye",
		"(Head) Cerebral Implant",
		"(Head) Cortex Implant",
		"(Head) Psy-Booster",
		"(Left arm) Augmetic Hand",
		"(Left arm) Elbow Actuator",
		"(Left arm) Muscle Casing",
		"(Right arm) Augmetic Elbow",
		"(Right arm) Locomotion Augmetics",
		"(Right arm) Synthmuscle",
		"(Torso) Armour Reinforcement",
		"(Torso) Augmetic Heart",
		"(Torso) Autosanguine",
		"(Torso) Respirary Filter",
		"(Left leg) Enhanced Knee Joint",
		"(Left leg) Subskin Leg Armour",
		"(Right leg) Augmetic Foot",
		"(Right leg) Subskin Leg Armour",
	}
)

type Augmetic struct {
	ID          string
	Name        string
	Description string
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

func (m *Manager) CanChangeUnitAugmetics(unit any, idx int, heal bool) (bool, Augmetic, []string) {
	switch object := unit.(type) {
	case *objects.KnightState:
		var curr Augmetic
		for i, augmetic := range object.Augmetics {
			if i == idx {
				curr = augmeticsByID[augmetic.Key]
			}
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
