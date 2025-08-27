package randomness

import "github.com/antithesishq/antithesis-sdk-go/random"

// Float64 converts Antithesis uint64 random to float64 in range [0, 1)
func Float64() float64 {
	return float64(random.GetRandom()) / float64(1<<64)
}

// Intn converts Antithesis random to int in range [0, n)
func Intn(n int) int {
	if n <= 0 {
		return 0
	}
	return int(random.GetRandom() % uint64(n))
}

// Choice returns a randomly chosen item from a list of options using Antithesis randomness
func Choice[T any](items []T) T {
	return random.RandomChoice(items)
}