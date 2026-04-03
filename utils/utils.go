package utils

import (
	"fmt"
)

const DEFAULT_COLOR = 0x2b2d31

func RandWord(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[i%len(charset)]
	}
	return string(result)
}

func HexToDecimal(hex string) (int, error) {
	var result int
	_, err := fmt.Sscanf(hex, "%x", &result)
	return result, err
}

func DecimalToHex(decimal int) string {
	return fmt.Sprintf("%06x", decimal)
}

func If[T any](cond bool, vTrue, vFalse T) T {
	if cond {
		return vTrue
	}
	return vFalse
}
