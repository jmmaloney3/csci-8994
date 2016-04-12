package sim

import "math/rand"
import "fmt"

// An agent that use an action module to decide how to act
type Agent struct {
  tribe *Tribe
  rep Rep
  actMod *ActionModule
  payout int32
  numGames int8
  pactmut float64 // mu_s - action module bit mutation probability
}

// Create a new agent.  By default the agent has a GOOD reputation.
func NewAgent(t *Tribe, pactmut float64, pexeerr float32, rnGen *rand.Rand) *Agent {
  var actm = NewActionModule(RandBool(rnGen), RandBool(rnGen),
                             RandBool(rnGen), RandBool(rnGen), pexeerr)
  return &Agent { tribe: t, rep: GOOD, payout: 0, numGames: 0, pactmut: pactmut, actMod: actm }
}

// Create a child of this agent to be part of the specified next generation
// tribe.  The new agent's action module is a clone of its parent's action
// module (possibly with mutations).  The new agent is added to the next generation
// tribe.  The other attributes are set to default values.
func (parent *Agent) CreateChild(nextGen *Tribe, rnGen *rand.Rand) *Agent {
  inheritedActMod := parent.actMod.CloneWithMutations(parent.pactmut, rnGen)
  return &Agent { tribe: nextGen, rep: GOOD, payout: 0, numGames: 0,
                  actMod: inheritedActMod, pactmut: parent.pactmut }
}

// generate string representation of an agent
func (a *Agent) String() string {
  str := "{\n"
  str = str + fmt.Sprintf("  \"type\":%T\n", a)
  str = str + fmt.Sprintf("  \"addr\":%p\n", a)
  str = str + fmt.Sprintf("  \"tribe\":%p\n", a.tribe)
  str = str + fmt.Sprintf("  \"rep\":%v\n", a.rep)
  str = str + fmt.Sprintf("  \"act-mod\":%v \n", a.actMod)
  str = str + fmt.Sprintf("  \"payout\":%d \n", a.payout)
  str = str + fmt.Sprintf("  \"num-games\":%d \n", a.numGames)
  str = str + "}"
  return str
}


// Reset the agent's internal state to prepare for participation in the
// next generation.
func (self *Agent) Reset() {
  self.rep = GOOD
  self.payout = 0
  self.numGames = 0
}

// Ask the agent to choose whether it will donate to the recipient agent.
// Returns true if the agent chooses to donate and false otherwise.
func (self *Agent) ChooseDonate(recipient *Agent, rnGen *rand.Rand) bool {
  return self.actMod.ChooseDonate(self.rep, recipient.rep, rnGen)
}

// Play a round of the IR game with this agent playing the role of the
// donor agent.  The total payout earned by both agents is returned.
func (self *Agent) PlayRound(recipient *Agent, cost int32, benefit int32, rnGen *rand.Rand) int32 {
  // increase number of games played
  self.numGames += 1
  recipient.numGames += 1

  // keep track of total payment earned by agents
  var totalPayout int32 = 0

  // set default action
  action := REFUSE

  // play round
  if (self.ChooseDonate(recipient, rnGen)) {
    // donor donates
    action = DONATE
    // -- recipient receives benefit
    recipient.payout += benefit
    // -- donor pays cost
    self.payout -= cost
    // update total payout
    totalPayout += (benefit - cost)
  }

  // update donor's reputation
  self.rep = self.tribe.assessMod.AssignRep(self.rep, recipient.rep, action, rnGen)

  // to prevent negative payout, each agent receives cost
  self.payout += cost
  recipient.payout += cost
  totalPayout += (2*cost)

  // return total payout earned by both agents
  return totalPayout
}

func (self *Agent) WriteSimParams() {
  fmt.Printf("  \"pactmut\":%.5f,\n", self.pactmut)
  self.actMod.WriteSimParams()
}
