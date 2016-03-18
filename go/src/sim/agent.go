package sim

// An agent that use an action module to decide how to act
type Agent struct {
  rep Rep
  actMod *ActionModule
  payout int32
  numGames int8
}

// Create a new agent.  By default the agent has a GOOD reputation.
func NewAgent() *Agent {
  var actm = NewActionModule()
  return &Agent { rep: GOOD, payout: 0, numGames: 0, actMod: actm }
}

// Reset the agent's internal state to prepare for participation in the
// next generation.
func (self *Agent) Reset() {
  self.payout = 0
  self.numGames = 0
}

// Ask the agent to choose whether it will donate to the recipient agent.
// Returns true if the agent chooses to donate and false otherwise.
func (self *Agent) ChooseDonate(recipient *Agent) bool {
  return self.actMod.ChooseDonate(self, recipient)
}

// Play a round of the IR game with this agent playing the role of the
// donor agent.  The total payout earned by both agents is returned.
func (self *Agent) PlayRound(recipient *Agent, cost int32, benefit int32) int32 {
  // increase number of games played
  self.numGames += 1
  recipient.numGames += 1

  // keep track of total payment earned by agents
  var totalPayout int32 = 0

  // play round
  if (self.ChooseDonate(recipient)) {
    // donor donates
    // -- recipient receives benefit
    recipient.payout += benefit
    // -- donor pays cost
    self.payout -= cost
    // update total payout
    totalPayout += (benefit - cost)
  }

  // to prevent negative payout, each agent receives cost
  self.payout += cost
  recipient.payout += cost
  totalPayout += (2*cost)

  // return total payout earned by both agents
  return totalPayout
}
