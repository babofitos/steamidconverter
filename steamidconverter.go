//Package steamidconverter helps convert various Steam ID formats
package steamidconverter

import (
	"fmt"
	"strconv"
	"strings"
)

//Steam identifier
var identifier uint64 = 76561197960265728

//ConvertToText converts a Steam64 ID to STEAM_0:Y:Z format
func ConvertToText(w uint64) string {
	y := w % 2
	v := identifier
	z := (w - y - v) / 2
	return fmt.Sprintf("STEAM_0:%d:%d", y, z)
}

//ConvertTo64 converts STEAM_X:Y:Z to a Steam64 format
func ConvertTo64(text string) (uint64, error) {
	sidSlice := strings.Split(text, ":")
	z, err := strconv.ParseUint(sidSlice[2], 10, 0)
	if err != nil {
		return 0, err
	}
	y, err := strconv.ParseUint(sidSlice[1], 10, 0)
	if err != nil {
		return 0, err
	}
	v := identifier
	w := z*2 + v + y
	return w, nil
}

//ConvertToSteam3 converts STEAM_X:Y:Z to a Steam3 [U:1:W] format
func ConvertToSteam3(text string) (string, error) {
	sidSlice := strings.Split(text, ":")
	z, err := strconv.ParseUint(sidSlice[2], 10, 0)
	if err != nil {
		return "", err
	}
	y, err := strconv.ParseUint(sidSlice[1], 10, 0)
	if err != nil {
		return "", err
	}
	w := z*2 + y
	return fmt.Sprintf("[U:1:%d]", w), nil
}
