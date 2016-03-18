package sim

// A simulation engine for simulating the indirect reciprocity game
// played among agents divided into tribes.
type SimEngine struct {
  tribes []*Tribe
  numTribes int
  totalPayouts int32
}

// Make a new simulation engine.
func NewSimEngine(numTribes int, numAgents int) *SimEngine {
  tribes := make([]*Tribe, numTribes)
  // create tribes
  for i := 0; i < numTribes; i++ {
    tribes[i] = NewTribe(numAgents)
  }
  return &SimEngine { tribes: tribes, numTribes: numTribes, totalPayouts: 0 }
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
