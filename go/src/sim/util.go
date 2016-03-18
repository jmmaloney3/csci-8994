package sim

import "crypto/rand"
import "math/big"

type Rep int

const (
  GOOD Rep = iota
  BAD Rep = iota
)

// Generate a random boolean
func randbool() bool {
  num, _ := rand.Int(rand.Reader, big.NewInt(2))
  if (num.Int64() == 0) {
    return false
  } else {
    return true
  }
}
