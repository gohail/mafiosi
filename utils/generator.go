package utils

import (
	"math/rand"
	"time"
)

const (
	MinID int = 10000
	MaxID int = 99999
)

func GenerateID() int {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(MaxID-MinID) + MinID
	return r
}

func RandomIntArr(num int) []int {
	list := rand.Perm(num)
	return list
}
