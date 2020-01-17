package utils

import (
	"fmt"
	"testing"
)

func TestGenerateId(t *testing.T) {
	random := GenerateID()

	fmt.Printf("random:%d \n", random)

	if random < MinID || random > MaxID {
		t.Fatalf("random out of a bound r:%d min:%d max:%d", random, MinID, MaxID)
	}
}
