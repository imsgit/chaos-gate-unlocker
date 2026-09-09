package internal

import (
	"chaos-gate-unlocker/internal/save"

	"encoding/json"
)

type Header struct {
	Version        string     `json:"version"`
	SavedTimeStamp save.Stamp `json:"savedTimeStamp"`
	GameDays       int        `json:"gameDays"`
	Location       string     `json:"location"`
	SlotType       int        `json:"slotType"`
	CampaignID     string     `json:"campaignID"`
	Difficulty     int        `json:"difficulty"`
	IronMan        bool       `json:"ironMan"`
	SaveName       string     `json:"saveName"`
}

type State struct {
	TopRecord         *LinearRecord   `json:"topRecord"`
	LinearInstanceIds []int           `json:"linearInstanceIds"`
	LinearRecords     []*LinearRecord `json:"linearRecords"`
}

type LinearRecord struct {
	TypeName           string
	AssetName          string
	SerializedContents json.RawMessage
	SerializedObject   any
}

type linearRecord struct {
	TypeName           string          `json:"typeName"`
	AssetName          string          `json:"assetName"`
	SerializedContents json.RawMessage `json:"serializedContents"`
}
