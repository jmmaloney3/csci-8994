package sim

import "math"
import "math/rand"
import "time"
import "runtime"
import "fmt"

// A simulation engine for simulating the indirect reciprocity game
// played among agents divided into tribes.
type SimEngine struct {
  tribes []*Tribe
  numTribes int
  totalPayouts int32
  rnGen *rand.Rand // hold a RN generator for sequential processing
  useMP bool
  numCpu int
  cpuTasks []int // when using MP, num tasks to assign to each CPU
  cpuRNG []*rand.Rand // a separate random number generator for each CPU
  pcon float32 // prob of tribal conflict: recommended 0.01
  beta float64 // selection strength varies from 10^0 to 10^5
  eta  float64 // recommended <= 0.2 (used 0.1 in supporting materials)
  pmig float32 // prob of migration: recommended 0.005
  passmut float32 // prob of assess module bit mutation: recommended 0.0001
  // define some standard action modules for comparison
  ALLD *ActionModule // a constant used for comparison during stats gathering
  ALLC *ActionModule // a constant used for comparison during stats gathering
}

func NewDefaultSimEngine(numTribes int, numAgents int, useMP bool) *SimEngine {
  // create parameter map
  var params = make(map[string]float64)

  // populate arg maps
  params[PASSE_F] = PASSERR
  params[PACTM_F] = PACTMUT
  params[PEXEE_F] = PEXEERR
  params[PCON_F]  = PCON
  params[BETA_F]  = BETA
  params[ETA_F]   = ETA
  params[PMIG_F]  = PMIG
  params[PASSM_F] = PASSMUT

  // create simulation engine with default values
  return NewSimEngine(numTribes, numAgents, params, useMP)
}

// Make a new simulation engine.
func NewSimEngine(numTribes int, numAgents int, params map[string]float64, useMP bool) *SimEngine {
  // get parameters
  passerr, ok := params[PASSE_F]
  if (!ok) { passerr = PASSERR }
  pactmut, ok := params[PACTM_F]
  if (!ok) { pactmut = PACTMUT }
  pexeerr, ok := params[PEXEE_F]
  if (!ok) { pexeerr = PEXEERR }
  pcon, ok := params[PCON_F]
  if (!ok) { pcon = PCON }
  beta, ok := params[BETA_F]
  if (!ok) { beta = BETA }
  eta, ok := params[ETA_F]
  if (!ok) { eta = ETA }
  pmig, ok := params[PMIG_F]
  if (!ok) { pmig = PMIG }
  passmut, ok := params[PASSM_F]
  if (!ok) { passmut = PASSMUT }

  // create tribes
  tribes := make([]*Tribe, numTribes)
  rnGen := rand.New(rand.NewSource(time.Now().UnixNano()))
  for i := 0; i < numTribes; i++ {
    tribes[i] = NewTribe(numAgents, float32(passerr), float32(pactmut), float32(pexeerr), rnGen)
  }
  // figure out multiprocessing parameters if MP enabled
  ncpu := runtime.NumCPU()
  cpuTasks := make([]int, ncpu)
  cpuRNG := make([]*rand.Rand, ncpu)
  if (useMP) {
    // figure out tasks per cpu - tasks might not evenly divide among CPUs
    tasksPerCpu := int(math.Ceil(float64(numTribes)/float64(ncpu)))
    taskSum := 0
    for i := 0; i < ncpu; i++ {
      cpuRNG[i] = rand.New(rand.NewSource(time.Now().UnixNano()))
      if ((numTribes - taskSum) > tasksPerCpu) {
        cpuTasks[i] = tasksPerCpu
        taskSum += tasksPerCpu
      } else {
        cpuTasks[i] = (numTribes - taskSum)
        taskSum += (numTribes - taskSum)
      }
    }
  }

  // create sim engine
  return &SimEngine { tribes: tribes, numTribes: numTribes, totalPayouts: 0,
                      pcon: float32(pcon), beta: beta, eta: eta, pmig: float32(pmig),
                      useMP: useMP, numCpu: ncpu, cpuTasks: cpuTasks, cpuRNG: cpuRNG,
                      rnGen: rnGen, passmut: float32(passmut),
                      ALLC: NewActionModule(true, true, true, true, 0),
                      ALLD: NewActionModule(false, false, false, false, 0) }
}

// Reset the simulations to prepare for participation in the next generation.
func (self *SimEngine) Reset() {
  self.totalPayouts = 0
  for i := 0; i < self.numTribes; i++ {
    self.tribes[i].Reset()
  }
}

// Play the required rounds of the IR game to complete the current generation
// and then create the next generation.
func (self *SimEngine) PlayRounds(cost int32, benefit int32) int32 {
  if (self.useMP) {
    // create channel to collect payouts from each parallel task
    payouts := make(chan int32, self.numCpu)
    tribeStart := 0
    tribeEnd := 0
    for i := 0; i < self.numCpu; i++ {
      tribeStart = tribeEnd
      tribeEnd = tribeStart + self.cpuTasks[i]
      go func (tribeStart int, tribeEnd int, rnGen *rand.Rand) {
        task_payouts := int32(0)
        for j := tribeStart; j < tribeEnd; j++ {
          task_payouts += self.tribes[j].PlayRounds(cost, benefit, rnGen)
          self.tribes[j].CreateNextGen(rnGen)
        }
        payouts <- task_payouts
      } (tribeStart, tribeEnd, self.cpuRNG[i])
    }
    // wait for goroutines to finish
    for i := 0; i < self.numCpu; i++ {
      self.totalPayouts += (<-payouts)
    }
  } else {
    for i := 0; i < self.numTribes; i++ {
      self.totalPayouts += self.tribes[i].PlayRounds(cost, benefit, self.rnGen)
      self.tribes[i].CreateNextGen(self.rnGen)
    }
  }
  return self.totalPayouts
}

// Calculate the maximum and minimum total payouts that could be earned by the agents
// in a single generation
func (self *SimEngine) MaxMinPayouts(cost int32, benefit int32) (max int32, min int32) {
  max = 0
  min = 0
  numAgents := self.tribes[0].numAgents
  for i := 0; i < numAgents; i++ {
    for j := i+1; j < numAgents; j++ {
      // add (benefit - cost) + (2*cost)
      max += (benefit + cost)
      min += 2*cost
    }
  }
  max = max*int32(self.numTribes)
  min = min*int32(self.numTribes)
  return max, min
}

// Create the next generation by propagating action modules to the next
// generation based on the fitness those modules achieved.
/*
func (self *SimEngine) CreateNextGen() {
  for i := 0; i < self.numTribes; i++ {
    self.tribes[i].CreateNextGen()
  }
}
*/
// Evolve the tribal assessment modules based on the average payouts
// earned by each tribe during the last generation
func (self *SimEngine) EvolveTribes() {
  // iterate over the tribes and select pairs for confict
  for i := 0; i < self.numTribes; i++ {
    for j := i+1; j < self.numTribes; j++ {
      if (RandPercent(self.rnGen) < float64(self.pcon)) {
        winner, loser := self.Conflict(self.tribes[i], self.tribes[j], self.rnGen)
        self.ShiftAssessMod(winner, loser, self.rnGen)
        self.MigrateAgents(winner, loser, self.rnGen)
      }
    }
  }
}

// Migrate some agents from the first tribe to the second tribe
func (self *SimEngine) MigrateAgents(from *Tribe, to *Tribe, rnGen *rand.Rand) {
  for i := 0; i < to.numAgents; i++ {
    if (RandPercent(rnGen) < float64(self.pmig)) {
      to.agents[i].actMod = from.agents[i].actMod
    }
  }
}

// Collect statistics for the most recently completed generation
func (self *SimEngine) GetStats() (assessStats [8]int, actionStats [4]int, allcCnt int, alldCnt int) {
  for i := 0; i < self.numTribes; i++ {
    // collect statistics on the tribe's assess module
    for j := 0; j < 8; j++ {
      assessStats[j] += self.tribes[i].assessMod.GetBit(j)
    }
    // collect statistics on the agent's action modules
    for k := 0; k < self.tribes[i].numAgents; k++ {
      for m := 0; m < 4; m++ {
        actionStats[m] += self.tribes[i].agents[k].actMod.GetBit(m)
      }
      // count occurences of ALLD and ALLC
      if (self.tribes[i].agents[k].actMod.SameBits(self.ALLD)) {
        alldCnt++
      } else if (self.tribes[i].agents[k].actMod.SameBits(self.ALLC)) {
        allcCnt++
      }
    }
  }
  return assessStats, actionStats, allcCnt, alldCnt
}

// Determine the tribe that wins the conflict
func (self *SimEngine) Conflict(tribeA *Tribe, tribeB *Tribe, rnGen *rand.Rand) (winner, loser *Tribe) {
  if (math.IsInf(self.beta, int(1))) {
    // if Beta is infinite then tribe with higher payout always wins
    if (tribeB.AvgPayout() > tribeA.AvgPayout()) {
      return tribeB, tribeA
    } else {
      // if A payout is greater than B payout or payouts are equal
      return tribeA, tribeB
    }
  } else {
    diff := tribeB.AvgPayout() - tribeA.AvgPayout()
    p  := math.Pow(float64(1) + math.Exp(diff*(-self.beta)), float64(-1))
    if (RandPercent(rnGen) > p) {
      return tribeB, tribeA
    } else {
      return tribeA, tribeB
    }
  }
}

// Shift the loser's assessment module toward the winner's assessment module
func (self *SimEngine) ShiftAssessMod(winner *Tribe, loser *Tribe, rnGen *rand.Rand) {
  // before changing the loser's assess module, make a copy in case
  // it is shared with another tribe
  loser.assessMod = loser.assessMod.Copy()
  // get average payouts
  poW := winner.AvgPayout()
  poL := loser.AvgPayout()
  pflip := (self.eta*poW)/((self.eta*poW) + (float64(1)-self.eta)*poL)
  //bits := loser.assessMod.GetBits()
  //wBits := winner.assessMod.GetBits()
  //fmt.Printf("before: %8b (%4d) => %8b (%4d)\n", bits, bits, wBits, wBits)
  for i := 0; i < 8; i++ {
    if (loser.assessMod.bits[i] != winner.assessMod.bits[i]) {
      if (RandPercent(rnGen) < pflip) {
        loser.assessMod.bits[i] = winner.assessMod.bits[i]
      }
    } else {
      if (RandPercent(rnGen) < float64(self.passmut)) {
        if (loser.assessMod.bits[i] == GOOD) {
          loser.assessMod.bits[i] = BAD
        } else {
          loser.assessMod.bits[i] = GOOD
        }
      }
      // mutation
    }
  }
  //bits = loser.assessMod.GetBits()
  //fmt.Printf("after:  %8b (%4d)\n", bits, bits)
}

func (self *SimEngine) WriteSimParams() {
  fmt.Printf("  \"ntribes\":%d,\n", self.numTribes)
  fmt.Printf("  \"beta\":%.5f,\n", self.beta)
  fmt.Printf("  \"eta\":%.5f,\n", self.eta)
  fmt.Printf("  \"pcon\":%.5f,\n", self.pcon)
  fmt.Printf("  \"pmig\":%.5f,\n", self.pmig)
  fmt.Printf("  \"passmut\":%.5f,\n", self.passmut)
  fmt.Printf("  \"mp\":%t,\n", self.useMP)
  fmt.Printf("  \"ncpu\":%d,\n", self.numCpu)
  // write tribe sim parameters
  self.tribes[0].WriteSimParams()
}
