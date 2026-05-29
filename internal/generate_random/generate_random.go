package generate_random

import (
	"math/rand/v2"
	"strconv"
)

func GenerateNumbers(length int) int {
	resStr := ""
	for len(resStr) < length {
		randomer := rand.IntN(58)
		if randomer > 47 && randomer < 58 {
			if len(resStr) == 0 && randomer == '0' {
				continue
			}
			resStr += string(byte(randomer))
		}
	}
	num, _ := strconv.Atoi(resStr)
	return num
}
func GenerateString(length int) string {
	resStr := ""
	for len(resStr) < length {
		randomer := rand.IntN(123)
		if (randomer > 47 && randomer < 58) || (randomer > 64 && randomer < 91) || randomer > 96 {
			resStr += string(byte(randomer))
		}
	}
	return resStr
}
