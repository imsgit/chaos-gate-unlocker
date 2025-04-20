package internal

import "github.com/goccy/go-json"

type SavedTimeStamp struct {
	Days    int `json:"days"`
	Months  int `json:"months"`
	Years   int `json:"years"`
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
	Seconds int `json:"seconds"`
}

type Header struct {
	Version        string         `json:"version"`
	SavedTimeStamp SavedTimeStamp `json:"savedTimeStamp"`
	GameDays       int            `json:"gameDays"`
	Location       string         `json:"location"`
	SlotType       int            `json:"slotType"`
	CampaignID     string         `json:"campaignID"`
	Difficulty     int            `json:"difficulty"`
	IronMan        bool           `json:"ironMan"`
	SaveName       string         `json:"saveName"`
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
	SerializedObject   interface{}
}

type linearRecord struct {
	TypeName           string          `json:"typeName"`
	AssetName          string          `json:"assetName"`
	SerializedContents json.RawMessage `json:"serializedContents"`
}
