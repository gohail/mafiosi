package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_Main(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	list := rand.Perm(6)
	fmt.Println(list)
}
