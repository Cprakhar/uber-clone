package util

import (
	"fmt"
	"math/rand/v2"
)

func GenerateRandomPlate() string {
	alpha := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	num := "0123456789"

	plate := fmt.Sprintf("%c%c%c-%c%c%c",
		alpha[rand.IntN(len(alpha))],
		alpha[rand.IntN(len(alpha))],
		alpha[rand.IntN(len(alpha))],
		num[rand.IntN(len(num))],
		num[rand.IntN(len(num))],
		num[rand.IntN(len(num))],
	)
	return plate
}
