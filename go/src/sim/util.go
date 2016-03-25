package sim

//import "crypto/rand"
//import "math/big"
import "math/rand"
import "time"

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

// Return a new random number generator.  This generator is NOT protected
// by a mutex lock and therefore not thread safe.
func NewRandNumGen() *rand.Rand {
  return rand.New(rand.NewSource(time.Now().UnixNano()))
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

// Generate a random integer in the range [0, max] from the provided source
func RandInt(source *rand.Rand, max int64) int64 {
  return source.Int63n(max)
}

// Generate a random number between zero and 1 from the provided source
func RandPercent(source *rand.Rand) float64 {
  return source.Float64()
}

// Generate a random Rep from the provided source
func RandRep(source *rand.Rand) Rep {
  if (RandBool(source)) {
    return GOOD
  } else {
    return BAD
  }
}

/*
// Generate a random boolean from the secure source
func CryptoRandBool() bool {
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

// Generate a random integer in the range [0, max] from the secure source
func CryptoRandInt(max int64) int64 {
  num, err := rand.Int(rand.Reader, big.NewInt(max))
  if (err != nil) {
    panic("sim.RandInt(): " + err.Error())
  }
  return num.Int64()
}

// Generate a random number between zero and 1 from the secure source
func CryptoRandPercent() float64 {
  i := CryptoRandInt(int64(100001))
  return float64(i)/float64(100000)
}

// Generate a random Rep from the secure source
func CryptoRandRep() Rep {
  if (CryptoRandBool()) {
    return GOOD
  } else {
    return BAD
  }
}
*/
