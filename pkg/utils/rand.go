package utils

import (
	"math/rand"
	"time"
)

func RandPoint() float64 {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(1000)
	return float64(r) / 1000
}
