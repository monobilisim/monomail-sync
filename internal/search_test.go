package internal

import (
	"math/rand"
	"strconv"
	"testing"
)

func BenchmarkSearchInQueue(b *testing.B) {
	initQueue()
	for i := 0; i < b.N; i++ {
		searchInQueue(strconv.Itoa(rand.Intn(2000)))
	}
}
