package sim

import "crypto/rand"
import "math/big"

type Rep int
const (
  GOOD Rep = iota
  BAD Rep = iota
)

type Act int
const (
  DONATE Act = iota
  REFUSE Act = iota
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

// Generate a random number between zero and 1
func RandPercent() float64 {
  i := RandInt(int64(100001))
  return float64(i)/float64(100000)
}

// Generate a randon Rep
func RandRep() Rep {
  if (RandBool()) {
    return GOOD
  } else {
    return BAD
  }
}
