package simulator

import (
	"math/rand"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func Pick[T any](items []T) T {
	return items[rand.Intn(len(items))]
}

func WeightedPickString(items []WeightedString) string {
	total := 0
	for _, item := range items {
		total += item.Weight
	}
	n := rand.Intn(total)

	for _, item := range items {
		if n < item.Weight {
			return item.Value
		}
		n -= item.Weight
	}
	return items[len(items)-1].Value

}

func RandomInt64(min, max int64) int64 {
	if max <= min {
		return min
	}
	return min + rand.Int63n(max-min+1)
}

func RandomAlphaNumeric(length int) string {
	result := make([]rune, length)

	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

func RandomReadableWord() string {
	return Pick(subDomains)
}
