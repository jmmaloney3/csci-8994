package sim

import "crypto/rand"
import "math/big"

type Rep int

const (
  GOOD Rep = iota
  BAD Rep = iota
)

// Generate a random boolean
func RandBool() bool {
  num, err := rand.Int(rand.Reader, big.NewInt(2))
  if (err != nil) {
    panic("sim.RandBool(): " + err.Error())
  }
  if (num.Int64() == 0) {
    return false
  } else {
    return true
  }
}

// Generate a random integer in the range [0, max].
func RandInt(max int64) int64 {
  num, err := rand.Int(rand.Reader, big.NewInt(max))
  if (err != nil) {
    panic("sim.RandInt(): " + err.Error())
  }
  return num.Int64()
}
