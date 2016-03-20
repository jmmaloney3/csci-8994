package sim

// A tribe of agents that uses an assessment module to assign reputations
// to agents.
type Tribe struct {
  agents []*Agent
  assessMod *AssessModule
  numAgents int
  totalPayouts int32
}

// Create a new tribe.
func NewTribe(numAgents int) *Tribe {
  // create the tribe
  var assm = NewAssessModule(RandRep(), RandRep(), RandRep(), RandRep(),
                             RandRep(), RandRep(), RandRep(), RandRep())
  t := &Tribe { assessMod: assm, numAgents: numAgents, totalPayouts: 0 }
  // create the tribe's agents
  t.agents = make([]*Agent, numAgents)
  // create agents
  for i := 0; i < numAgents; i++ {
    t.agents[i] = NewAgent(t)
  }

  return t
}

// Reset the tribe's agents to prepare for participation in the next generation.
func (self *Tribe) Reset() {
  self.totalPayouts = 0;
  for i := 0; i < self.numAgents; i++ {
    self.agents[i].Reset()
  }
}

// Play the required rounds of the IR game to complete the current generation.
func (self *Tribe) PlayRounds(cost int32, benefit int32) int32 {
  var donor *Agent
  var recipient *Agent
  for i := 0; i < self.numAgents; i++ {
    for j := i+1; j < self.numAgents; j++ {
      // randomly assign the agents to roles
      if (RandBool()) {
        // agent i is donor and agent j is recipient
        donor = self.agents[i]
        recipient = self.agents[j]
      } else {
        // agent j is donor and agent i is recipient
        donor = self.agents[j]
        recipient = self.agents[i]
      }

      // play the round
      self.totalPayouts += donor.PlayRound(recipient, cost, benefit)
    }
  }

  // return the total payouts for use by the sim engine
  return self.totalPayouts
}

// Randomly select an agent from the local population.  The chance that an
// agent is selected is proportional to its fitness.
func (self *Tribe) SelectParent() *Agent {
  ri := int32(RandInt(int64(self.totalPayouts)))
  thresh := int32(0);
  var parent *Agent
  for i := 0; i < self.numAgents; i++ {
    thresh += self.agents[i].payout
    if (ri <= thresh) {
      parent = self.agents[i]
      break
    }
  }
  return parent
}

// Create the next generation by propagating action modules to the next
// generation based on the fitness those modules achieved.
func (self *Tribe) CreateNextGen() {
  for i := 0; i < self.numAgents; i++ {
    parent := self.SelectParent()
    self.agents[i].actMod = parent.actMod;
  }
}

// Return the average payout for an agent in this tribe
func (self *Tribe) AvgPayout() float64 {
  return float64(self.totalPayouts)/float64(self.numAgents)
}