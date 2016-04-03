package sim

import "math/rand"
import "fmt"

// A tribe of agents that uses an assessment module to assign reputations
// to agents.
type Tribe struct {
  agents []*Agent
  assessMod *AssessModule
  numAgents int
  totalPayouts int32
}

// Create a new tribe.
func NewTribe(numAgents int, passerr float32, pactmut float32, pexeerr float32, rnGen *rand.Rand) *Tribe {
  // create the tribe
  var assm = NewAssessModule(RandRep(rnGen), RandRep(rnGen), RandRep(rnGen), RandRep(rnGen),
                             RandRep(rnGen), RandRep(rnGen), RandRep(rnGen), RandRep(rnGen),
                             passerr)
  t := &Tribe { assessMod: assm, numAgents: numAgents, totalPayouts: 0 }
  // create the tribe's agents
  t.agents = make([]*Agent, numAgents)
  // create agents
  for i := 0; i < numAgents; i++ {
    t.agents[i] = NewAgent(t, pactmut, pexeerr, rnGen)
  }

  return t
}

// return the list of agents
func (t *Tribe) GetAgents() []*Agent {
  return t.agents
}

// generate string representation of a tribe
func (t *Tribe) String() string {
  str := "{\n"
  str = str + fmt.Sprintf("  \"type\":%T\n", t)
  str = str + fmt.Sprintf("  \"addr\":%p\n", t)
  str = str + fmt.Sprintf("  \"num-agents\":%d\n", t.numAgents)
  str = str + fmt.Sprintf("  \"assess-mod\":%v\n", t.assessMod)
  str = str + fmt.Sprintf("  \"total-payouts\":%d \n", t.totalPayouts)
  str = str + fmt.Sprintf("  \"agents\":%d \n", t.agents)
  str = str + "}"
  return str
}

// make a shallow copy of the tribe
// -- agents will be transferred from original tribe to copy
func (t *Tribe) ShallowCopy() *Tribe {
  copy := &Tribe { assessMod: t.assessMod, numAgents: t.numAgents,
                   totalPayouts: t.totalPayouts, agents: t.agents }
  // link agents with this new shallow copy
  for i := 0; i < copy.numAgents; i++ {
    copy.agents[i].tribe = copy
  }
  return copy
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
  random_idx := rnGen.Perm(self.numAgents)
  for idx, i := range random_idx {
    for _, j := range random_idx[idx+1:] {
      // randomly assign the agents to roles
      donor, recipient = self.AssignRoles(self.agents[i], self.agents[j], rnGen)

      // play the round
      self.totalPayouts += donor.PlayRound(recipient, cost, benefit, rnGen)
    }
  }

  // return the total payouts for use by the sim engine
  return self.totalPayouts
}

// Randomly assign the agents to the donor and recipient roles
func (self *Tribe) AssignRoles(a1 *Agent, a2 *Agent, rnGen *rand.Rand) (donor, recipient *Agent) {
  // randomly assign the agents to roles
  if (RandBool(rnGen)) {
    // agent 1 is donor and agent 2 is recipient
    donor = a1
    recipient = a2
  } else {
    // agent 2 is donor and agent 1 is recipient
    donor = a2
    recipient = a1
  }
  return donor,recipient
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

// Create the next generation by propagating action modules to the next
// generation based on the fitness those modules achieved.
func (currentGen *Tribe) CreateNextGen(rnGen *rand.Rand) *Tribe {
  // create the next generation tribe
  nextGen := &Tribe { assessMod: currentGen.assessMod, numAgents: currentGen.numAgents }
  // create the next generation of agents
  nextGen.agents = make([]*Agent, nextGen.numAgents)
  var parent *Agent
  for i := 0; i < nextGen.numAgents; i++ {
    // select parent from the current generation
    parent = currentGen.SelectParent(rnGen)
    // create a child of the parent and add to next generation
    nextGen.agents[i] = parent.CreateChild(nextGen, rnGen)
  }
  return nextGen
}

// Return the average payout for an agent in this tribe
func (self *Tribe) AvgPayout() float64 {
  return float64(self.totalPayouts)/float64(self.numAgents)
}

func (self *Tribe) WriteSimParams() {
  fmt.Printf("  \"nagents\":%d,\n", self.numAgents)
  // write assess module parameters
  self.assessMod.WriteSimParams()
  // write agent parameters
  self.agents[0].WriteSimParams()
}

// type and function to support sorting tribes by payouts
// -- tribes are sorted in ascending order of payouts
type SortTribesByPayouts []*Tribe
func (tribes SortTribesByPayouts) Len() int {
  return len(tribes)
}
func (tribes SortTribesByPayouts) Swap(i, j int) {
  tribes[i], tribes[j] = tribes[j], tribes[i]
}
func (tribes SortTribesByPayouts) Less(i, j int) bool {
  return tribes[i].totalPayouts < tribes[j].totalPayouts
}
