package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateId(t *testing.T) {
	random := GenerateID()

	fmt.Printf("random:%d \n", random)

	if random < MinID || random > MaxID {
		t.Fatalf("random out of a bound r:%d min:%d max:%d", random, MinID, MaxID)
	}
}

func TestRandomIntArr(t *testing.T) {
	arr := RandomIntArr(5)
	assert.Equal(t, 5, len(arr))
	assert.Equal(t, true, checkUniqInt(RandomIntArr(5)))
	assert.Equal(t, true, checkUniqInt(RandomIntArr(10)))
}

func checkUniqInt(arr []int) bool {
	countArr := make([]int, len(arr))
	for _, v := range arr {
		countArr[v]++
	}
	for _, v := range countArr {
		if v != 1 {
			return false
		}
	}
	return true
}
