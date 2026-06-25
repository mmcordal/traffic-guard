package simulator

import (
	"math/rand"
	"strings"
)

type WeightedString struct{}

func Pick[T any](items []T) T {
	return items[0]
}

func WeightedPickString(items []WeightedString) string {
	return ""
}

func RandomInt64(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

func RandomAlphaNumeric(length int) string {
	return strings.Repeat("a", length)
}

func RandomReadableWord() string {
	return ""
}
