package cordutils

import (
	"fmt"
)

func HexToDecimal(hex string) (int, error) {
	var result int
	_, err := fmt.Sscanf(hex, "%x", &result)
	return result, err
}

func DecimalToHex(decimal int) string {
	return fmt.Sprintf("%06x", decimal)
}
