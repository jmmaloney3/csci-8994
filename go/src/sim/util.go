package sim

//import "crypto/rand"
//import "math/big"
import "math/rand"
import "time"

// default parameter values
const (
 NUMGENS = 10 // default number of generations per simulation
 GENS_F = "g"
 COST = 1  // default donation cost
 COST_F = "c"
 BENEFIT = 3 // default donation benefit
 BEN_F = "b"
 NUMTRIBES = 64 // default number of tribes in a simulation
 TRIBES_F = "t"
 BETA = 1.2 // default value of conflict selection strength
 BETA_F = "beta"
 ETA = 0.15 // default value of bit switching selection strength
 ETA_F = "eta"
 PCON = 0.01 // default probability of tribal conflict
 PCON_F = "pcon"
 PMIG = 0.005 // default probability of migration
 PMIG_F = "pmig"
 PASSMUT = 0.0001 // default probability of assess module bit mutation
 PASSM_F = "passmut"
 NUMAGENTS = 64 // default number of agents per tribe
 AGENTS_F = "a"
 PACTMUT = 0.01 // default probability of action module bit mutation
 PACTM_F = "pactmut"
 PASSERR = 0.001 // default probability of assessment error
 PASSE_F = "passerr"
 PEXEERR = 0.001 // defaul tprobability of execution error
 PEXEE_F = "pexeerr"
 FNAME = "stats.csv"
 FNAME_F = "f"
 USEMP = false
 USEMP_F = "mp"
)

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
