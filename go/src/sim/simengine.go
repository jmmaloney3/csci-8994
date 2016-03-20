package sim

import "math"

// A simulation engine for simulating the indirect reciprocity game
// played among agents divided into tribes.
type SimEngine struct {
  tribes []*Tribe
  numTribes int
  totalPayouts int32
  pConflict float64
  beta float64 // selection strength
  eta float64
}

// Make a new simulation engine.
func NewSimEngine(numTribes int, numAgents int) *SimEngine {
  tribes := make([]*Tribe, numTribes)
  // create tribes
  for i := 0; i < numTribes; i++ {
    tribes[i] = NewTribe(numAgents)
  }
  // configure pConflict to 0.01
  return &SimEngine { tribes: tribes, numTribes: numTribes, totalPayouts: 0,
                      pConflict: 001, beta: 1.2, eta: 0.15 }
}

// Reset the simulations to prepare for participation in the next generation.
func (self *SimEngine) Reset() {
  self.totalPayouts = 0
  for i := 0; i < self.numTribes; i++ {
    self.tribes[i].Reset()
  }
}

// Play the required rounds of the IR game to complete the current generation.
func (self *SimEngine) PlayRounds(cost int32, benefit int32) int32 {
  for i := 0; i < self.numTribes; i++ {
    self.totalPayouts += self.tribes[i].PlayRounds(cost, benefit)
  }
  return self.totalPayouts
}

// Create the next generation by propagating action modules to the next
// generation based on the fitness those modules achieved.
func (self *SimEngine) CreateNextGen() {
  for i := 0; i < self.numTribes; i++ {
    self.tribes[i].CreateNextGen()
  }
}

// Evolve the tribal assessment modules based on the average payouts
// earned by each tribe during the last generation
func (self *SimEngine) EvolveTribes() {
  // iterate over the tribes and select pairs for confict
  for i := 0; i < self.numTribes; i++ {
    for j := i+1; j < self.numTribes; j++ {
      if (RandPercent() > self.pConflict) {
        winner, loser := self.Conflict(self.tribes[i], self.tribes[j])
        self.ShiftAssessMod(winner, loser)
      }
    }
  }
}

// Determine the tribe that wins the conflict
func (self *SimEngine) Conflict(tribeA *Tribe, tribeB *Tribe) (winner, loser *Tribe) {
  diff := tribeB.AvgPayout() - tribeA.AvgPayout()
  p := math.Pow(float64(1) + math.Exp(diff*(-self.beta)), float64(-1))
  if (RandPercent() > p) {
    return tribeB, tribeA
  } else {
    return tribeA, tribeB
  }
}

// Shift the loser's assessment module toward the winner's assessment module
func (self *SimEngine) ShiftAssessMod(winner *Tribe, loser *Tribe) {
  // copy assess module in case this is shared with another tribe
  loser.assessMod = CopyAssessModule(*loser.assessMod)
  // get average payouts
  poW := winner.AvgPayout()
  poL := loser.AvgPayout()
  p := (self.eta*poW)/(self.eta*poW + (1-self.eta)*poL)
  for i := 0; i < 8; i++ {
    if (RandPercent() > p) {
      loser.assessMod.bits[i] = winner.assessMod.bits[i]
    }
  }
}
