package sim

import "math"
import "math/rand"
import "time"
import "runtime"
//import "fmt"

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
  conP float32 // prob of tribal conflict: recommended 0.01
  Beta float64 // selection strength varies from 10^0 to 10^5
  eta  float64 // recommended <= 0.2 (used 0.1 in supporting materials)
  migP float32 // prob of migration: recommended 0.005
  mutP float32 // prob of assess module bit mutation: recommended 0.0001
}

// Make a new simulation engine.
func NewSimEngine(numTribes int, numAgents int, useMP bool) *SimEngine {
  tribes := make([]*Tribe, numTribes)
  // create tribes
  rnGen := rand.New(rand.NewSource(time.Now().UnixNano()))
  for i := 0; i < numTribes; i++ {
    tribes[i] = NewTribe(numAgents, rnGen)
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
  // configure pConflict to 0.01
  return &SimEngine { tribes: tribes, numTribes: numTribes, totalPayouts: 0,
                      conP: 0.01, Beta: 1.2, eta: 0.15, migP: 0.005,
                      useMP: useMP, numCpu: ncpu, cpuTasks: cpuTasks, cpuRNG: cpuRNG,
                      rnGen: rnGen, mutP: 0.0001 }
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

// Calculate the maximum total payout that could be earned by the agents
// in a single generation
func (self *SimEngine) MaxPayouts(cost int32, benefit int32) int32 {
  var max int32 = 0
  numAgents := self.tribes[0].numAgents
  for i := 0; i < numAgents; i++ {
    for j := i+1; j < numAgents; j++ {
      // add (benefit - cost) + (2*cost)
      max += (benefit + cost)
    }
  }
  return max*int32(self.numTribes)
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
      if (RandPercent(self.rnGen) < float64(self.conP)) {
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
    if (RandPercent(rnGen) < float64(self.migP)) {
      to.agents[i].actMod = from.agents[i].actMod
    }
  }
}

// Collect statistics for the most recently completed generation
func (self *SimEngine) GetStats() (assess_stats [8]int, action_stats [4]int) {
  for i := 0; i < self.numTribes; i++ {
    // collect statistics on the tribe's assess module
    for j := 0; j < 8; j++ {
      assess_stats[j] += self.tribes[i].assessMod.GetBit(j)
    }
    // collect statistics on the agent's action modules
    for k := 0; k < self.tribes[i].numAgents; k++ {
      for m := 0; m < 4; m++ {
        action_stats[m] += self.tribes[i].agents[k].actMod.GetBit(m)
      }
    }
  }
  return assess_stats, action_stats
}

// Determine the tribe that wins the conflict
func (self *SimEngine) Conflict(tribeA *Tribe, tribeB *Tribe, rnGen *rand.Rand) (winner, loser *Tribe) {
  if (math.IsInf(self.Beta, int(1))) {
    // if Beta is infinite then tribe with higher payout always wins
    if (tribeB.AvgPayout() > tribeA.AvgPayout()) {
      return tribeB, tribeA
    } else {
      // if A payout is greater than B payout or payouts are equal
      return tribeA, tribeB
    }
  } else {
    diff := tribeB.AvgPayout() - tribeA.AvgPayout()
    p  := math.Pow(float64(1) + math.Exp(diff*(-self.Beta)), float64(-1))
    if (RandPercent(rnGen) > p) {
      return tribeB, tribeA
    } else {
      return tribeA, tribeB
    }
  }
}

// Shift the loser's assessment module toward the winner's assessment module
func (self *SimEngine) ShiftAssessMod(winner *Tribe, loser *Tribe, rnGen *rand.Rand) {
  // copy assess module in case this is shared with another tribe
  loser.assessMod = CopyAssessModule(*loser.assessMod)
  // get average payouts
  poW := winner.AvgPayout()
  poL := loser.AvgPayout()
  p  := (self.eta*poW)/((self.eta*poW) + (float64(1)-self.eta)*poL)
  //bits := loser.assessMod.GetBits()
  //wBits := winner.assessMod.GetBits()
  //fmt.Printf("before: %8b (%4d) => %8b (%4d)\n", bits, bits, wBits, wBits)
  for i := 0; i < 8; i++ {
    if (loser.assessMod.bits[i] != winner.assessMod.bits[i]) {
      if (RandPercent(rnGen) < p) {
        loser.assessMod.bits[i] = winner.assessMod.bits[i]
      }
    } else {
      if (RandPercent(rnGen) < float64(self.mutP)) {
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
