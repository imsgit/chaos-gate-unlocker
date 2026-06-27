package display

import (
	"os/exec"
	"regexp"
	"strconv"
)

var resolutionRe = regexp.MustCompile(`resolution:\s+(\d+)x`)

func IsHiDPI() bool {
	out, err := exec.Command("xdpyinfo").Output()
	if err != nil {
		return false
	}
	m := resolutionRe.FindStringSubmatch(string(out))
	if len(m) != 2 {
		return false
	}
	dpi, _ := strconv.Atoi(m[1])
	return dpi > 96
}
