package utils

import (
	"strconv"
	"strings"
)

// HexToDecimal - convert hex values to decimal
func HexToDecimal(hex string) (uint64, error) {
	cleaned := strings.Replace(hex, "0x", "", -1)
	result, err := strconv.ParseUint(cleaned, 16, 64)

	return result, err
}
