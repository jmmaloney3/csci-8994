package simgpgg

import "math/rand"
import "time"

// Return a new random number generator.  This generator is NOT protected
// by a mutex lock and therefore not thread safe.
func NewRandNumGen() *rand.Rand {
  return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Generate a random integer in the range [0, max) from the provided source
func RandInt(source *rand.Rand, max int64) int64 {
  if (max == 0) { return 0 }
  return source.Int63n(max)
}

// Generate a random boolean from the provided source
func RandBool(source *rand.Rand) bool {
  num := source.Intn(2)
  if (num == 0) {
    return false
  } else {
    return true
  }
}

// Generate a random number between zero and 1 from the provided source
func RandProb(source *rand.Rand) float64 {
  return source.Float64()
}
