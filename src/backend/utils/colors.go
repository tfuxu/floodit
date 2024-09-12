package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// HexToCairoRGB converts hexadecimal color codes to float64 representations of red, green and blue color channels.
func HexToCairoRGB(hexCode string) ([3]float64, error) {
	if !strings.HasPrefix(hexCode, "#") {
		return [3]float64{}, errors.New("provided string isn't prefixed with hash (#) character")
	}

	if len(hexCode[1:]) > 6 {
		return [3]float64{}, errors.New("provided string exceeds the length of 6 characters")
	}

	i, err := strconv.ParseInt(hexCode[1:], 16, 32)
	if err != nil {
		return [3]float64{}, fmt.Errorf("an error occurred during parsing of the hexadecimal code: %v", err)
	}

	r := float64(i>>16) / 255
	g := float64((i>>8)&0xff) / 255
	b := float64(i&0xff) / 255

	return [3]float64{r, g, b}, nil
}
