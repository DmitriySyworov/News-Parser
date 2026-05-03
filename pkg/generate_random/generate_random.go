package generate_random

import (
	"math/rand/v2"
	"strconv"
)

func GenerateNumbers(length int) int {
	resStr := ""
	for len(resStr) < length{
		randomer := rand.IntN(58)
		if randomer >47 && randomer <58 {
			if len(resStr)== 0 && randomer == '0'{
				continue
			}
			resStr += string(byte(randomer))
		}
	}
	num, _ := strconv.Atoi(resStr)
	return num
}
