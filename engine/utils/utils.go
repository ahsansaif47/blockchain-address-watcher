package utils

import "strconv"

func StringToInteger(s string) (int, error) {
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1, err
	}
	return int(num), nil
}
