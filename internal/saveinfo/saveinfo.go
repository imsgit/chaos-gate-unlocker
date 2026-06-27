package saveinfo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ParseFile(path string) Info {
	f, err := os.Open(path)
	if err != nil {
		return Info{}
	}
	defer f.Close()

	buf := make([]byte, 64<<10)
	n, _ := io.ReadFull(f, buf)
	return Parse(buf[:n])
}

type Info struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type header struct {
	SaveName       string `json:"saveName"`
	GameDays       int    `json:"gameDays"`
	Difficulty     int    `json:"difficulty"`
	IronMan        bool   `json:"ironMan"`
	SavedTimeStamp struct {
		Days    int `json:"days"`
		Months  int `json:"months"`
		Years   int `json:"years"`
		Hours   int `json:"hours"`
		Minutes int `json:"minutes"`
	} `json:"savedTimeStamp"`
}

func Parse(data []byte) Info {
	first := data
	if i := bytes.Index(data, []byte("\r\n")); i >= 0 {
		first = data[:i]
	}

	var h header
	_ = json.Unmarshal(first, &h)

	ts := h.SavedTimeStamp
	return Info{
		Title: h.SaveName,
		Detail: fmt.Sprintf("Day %d   ·   %s   ·   %04d-%02d-%02d %02d:%02d",
			h.GameDays, difficultyName(h.Difficulty, h.IronMan),
			ts.Years, ts.Months, ts.Days, ts.Hours, ts.Minutes),
	}
}

func difficultyName(d int, ironMan bool) string {
	var s string
	switch d {
	case 3:
		s = "Legendary"
	case 2:
		s = "Ruthless"
	case 1:
		s = "Standard"
	case 0:
		s = "Merciful"
	default:
		s = "Unknown"
	}
	if ironMan {
		s += " Ironman"
	}
	return s
}
