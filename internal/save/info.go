package save

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

type Stamp struct {
	Days    int `json:"days"`
	Months  int `json:"months"`
	Years   int `json:"years"`
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
	Seconds int `json:"seconds"`
}

func (s Stamp) String() string {
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",
		s.Years, s.Months, s.Days, s.Hours, s.Minutes, s.Seconds)
}

func Detail(gameDays, difficulty int, ironMan bool, ts Stamp) string {
	return fmt.Sprintf("Day %d   ·   %s   ·   %s",
		gameDays, DifficultyName(difficulty, ironMan), ts)
}

type header struct {
	SaveName       string `json:"saveName"`
	GameDays       int    `json:"gameDays"`
	Difficulty     int    `json:"difficulty"`
	IronMan        bool   `json:"ironMan"`
	SavedTimeStamp Stamp  `json:"savedTimeStamp"`
}

func Parse(data []byte) Info {
	first := data
	if i := bytes.Index(data, []byte("\r\n")); i >= 0 {
		first = data[:i]
	}

	var h header
	_ = json.Unmarshal(first, &h)

	return Info{
		Title:  h.SaveName,
		Detail: Detail(h.GameDays, h.Difficulty, h.IronMan, h.SavedTimeStamp),
	}
}

func DifficultyName(d int, ironMan bool) string {
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
