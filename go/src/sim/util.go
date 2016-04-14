package sim

//import "crypto/rand"
//import "math/big"
import "math"
import "math/rand"
import "time"
import "fmt"

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
 ETA = 0.1 // default value of bit switching selection strength
 ETA_F = "eta"
 PCON = 0.01 // default probability of tribal conflict
 PCON_F = "pcon"
 SINGLE_DEF = true // whether a tribe is limited to a single defeat per generation
 SINGLE_DEF_F = "singdef"
 PMIG = 0.005 // default probability of migration
 PMIG_F = "pmig"
 PASSMUT = 0.0001 // default probability of assess module bit mutation
 PASSM_F = "passmut"
 PASSMUT_ALL = false // default value always attempting to mutate assmod bits
 PASSMUT_ALL_F = "passmutall"
 NUMAGENTS = 64 // default number of agents per tribe
 AGENTS_F = "a"
 PACTMUT = 0.01 // default probability of action module bit mutation
 PACTM_F = "pactmut"
 PASSERR = 0.001 // default probability of assessment error
 PASSE_F = "passerr"
 PEXEERR = 0.001 // defaul tprobability of execution error
 PEXEE_F = "pexeerr"
 USEAM = false
 USEAM_F = "am"
 FNAME = "stats.csv"
 FNAME_F = "f"
 NOMP = false
 NOMP_F = "nmp"
 ALLD = 0
 ALLC = 15
)

type Rep int
const (
  GOOD Rep = iota
  BAD Rep = iota
)

func (r Rep) String() string {
  if (r == GOOD) { return "GOOD" } else { return "BAD" }
}

type Act int
const (
  DONATE Act = iota
  REFUSE Act = iota
)

const epsilon = float64(0.00000001)
// determine whether two floating point numbers are equal
// http://stackoverflow.com/questions/4915462/how-should-i-do-floating-point-comparison
func FloatAlmostEquals(f1 float64, f2 float64, epsilon float64) bool {
  if (f1 == f2) {
    return true
  } else {
    diff := math.Abs(f1 - f2)
    if (f1 == 0 || f2 == 0 || diff < math.SmallestNonzeroFloat64) {
      // f1 or f2 is zero or both are extremely close to it
      // relative error is less meaningful here
      return diff < (epsilon * math.SmallestNonzeroFloat64)
    } else { // use relative error
      return ((diff / (math.Abs(f1) + math.Abs(f2))) < epsilon)
    }
  }
}

const lowFitMutRate = float64(0.002)
// Calculate an adaptive mutation rate for a tribe based on the provided total
// payouts, min payout and max payout
func CalcAdaptTribalMutRate(totalPayouts float64, minPO, maxPO int32) float64 {
  // if the tribe earned the minimum payout then return low fit mutation rate
  if (FloatAlmostEquals(totalPayouts, float64(minPO), epsilon)) {
    return lowFitMutRate
  }

  // error checks
  if (minPO >= maxPO) {
    msg := fmt.Sprintf("minPO >= maxPO (minPO: %d, maxPO: %d)", minPO, maxPO)
    panic(msg)
  }

  // calculate the percent of possible payout earned by the tribe
  // -- calculate earned payout - the amount above the minimum payout
  earnedPO := totalPayouts - float64(minPO)
  // -- calculate the possible max earned payout - the diff between max and min payouts
  maxEarnedPO := float64(maxPO - minPO)

  // calculate mutation rate based on parent fitness
  // -- linear mutation rate
  // return lowFitMutRate - (earnedPO/maxEarnedPO)*lowFitMutRate
  // -- error rate based
  // errRate := float64(1) - earnedPO/maxEarnedPO
  // return PASSMUT + math.Pow(errRate, float64(4))*float64(0.002)
  // -- exponential mutation rate
  return lowFitMutRate * math.Exp(math.Log(0.5)*5*earnedPO/maxEarnedPO)
  //return lowFitMutRate * math.Exp(-6.2146*5*earnedPO/maxEarnedPO)
}

// Calculate the minimum and maximum total payouts that can be earned by a tribe
// in a single generation
func CalcMinMaxTribalPayouts(numAgents int, cost int32, benefit int32) (min int32, max int32) {
  max = 0
  min = 0
  for i := 0; i < numAgents; i++ {
    for j := i+1; j < numAgents; j++ {
      // add (benefit - cost) + (2*cost)
      max += (benefit + cost)
      min += 2*cost
    }
  }
  return min, max
}

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
  if (max == 0) { return 0 }
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
