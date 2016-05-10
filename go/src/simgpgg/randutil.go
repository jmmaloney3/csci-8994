package simgpgg

import "math/rand"
import "math"
import "time"

// return the value of the Fermi distribution for the specified values
// of beta, p1 and p2: 1 / (1 + e^(-beta*(p1-p2)))
func Fermi(beta, p1, p2 float64) float64 {
  if (beta < 0) {
    panic("beta < 0")
  }
  one := float64(1)
  if (math.IsInf(beta, +1)) {
    return one
  } else if (beta == float64(0)) {
    return float64(0.5)
  } else {
    exp := (-beta)*(p1 - p2)
    return one/(one + math.Exp(exp))
  }
}

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
