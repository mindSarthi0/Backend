package lib

import (
	"fmt"
	"strconv"
)

// ConvertToInt converts a string to an integer and handles errors.
func ConvertToInt(s string) (int, error) {
	convertedInt, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to convert '%s' to int: %w", s, err)
	}
	return convertedInt, nil
}
