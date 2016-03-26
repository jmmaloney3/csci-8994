package sim

import "math/rand"

// A tribe of agents that uses an assessment module to assign reputations
// to agents.
type Tribe struct {
  agents []*Agent
  assessMod *AssessModule
  numAgents int
  totalPayouts int32
  mutP float32 // mu_s - mutation probability
}

// Create a new tribe.
func NewTribe(numAgents int, rnGen *rand.Rand) *Tribe {
  // create the tribe
  var assm = NewAssessModule(RandRep(rnGen), RandRep(rnGen), RandRep(rnGen), RandRep(rnGen),
                             RandRep(rnGen), RandRep(rnGen), RandRep(rnGen), RandRep(rnGen))
  t := &Tribe { assessMod: assm, numAgents: numAgents, totalPayouts: 0, mutP: 0.01 }
  // create the tribe's agents
  t.agents = make([]*Agent, numAgents)
  // create agents
  for i := 0; i < numAgents; i++ {
    t.agents[i] = NewAgent(t, rnGen)
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
func (self *Tribe) PlayRounds(cost int32, benefit int32, rnGen *rand.Rand) int32 {
  var donor *Agent
  var recipient *Agent
  // randomize the order of the agents
  outer_idx := rnGen.Perm(self.numAgents)
  for _, i := range outer_idx {
    // randomize the inner loop
    inner_idx := rnGen.Perm(self.numAgents - (i + 1))
    for _, j := range inner_idx {
      // adjust index j to be in right range
      j = j + (i + 1)
      // randomly assign the agents to roles
      if (RandBool(rnGen)) {
        // agent i is donor and agent j is recipient
        donor = self.agents[i]
        recipient = self.agents[j]
      } else {
        // agent j is donor and agent i is recipient
        donor = self.agents[j]
        recipient = self.agents[i]
      }

      // play the round
      self.totalPayouts += donor.PlayRound(recipient, cost, benefit, rnGen)
    }
  }

  // return the total payouts for use by the sim engine
  return self.totalPayouts
}

// Randomly select an agent from the local population.  The chance that an
// agent is selected is proportional to its fitness.
func (self *Tribe) SelectParent(rnGen *rand.Rand) *Agent {
  ri := int32(RandInt(rnGen, int64(self.totalPayouts)))
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

// Randomly select an agent from the local population.  Each agent has an equal
// chance of being selected.
func (self *Tribe) SelectMutationParent(rnGen *rand.Rand) *Agent {
  // select the index of the agent
  i := RandInt(rnGen, int64(self.numAgents))
  return self.agents[i]
}

// Create the next generation by propagating action modules to the next
// generation based on the fitness those modules achieved.
func (self *Tribe) CreateNextGen(rnGen *rand.Rand) {
  var parent *Agent
  for i := 0; i < self.numAgents; i++ {
    // select parent
    if (RandPercent(rnGen) < float64(self.mutP)) {
      parent = self.SelectMutationParent(rnGen)
    } else {
      parent = self.SelectParent(rnGen)
    }
    // inherit parent's action module
    self.agents[i].actMod = parent.actMod;
  }
}

// Return the average payout for an agent in this tribe
func (self *Tribe) AvgPayout() float64 {
  return float64(self.totalPayouts)/float64(self.numAgents)
}
