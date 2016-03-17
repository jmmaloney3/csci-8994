package sim

type Rep int

const (
  GOOD Rep = iota
  BAD Rep = iota
)

type ActionModule struct {
  bits [4]bool
}

func MakeActionModule() ActionModule {
  return ActionModule { bits: [4]bool{true, false, true, false} }
}

func (self ActionModule) ChooseDonate(donor Agent, recipient Agent) bool {
  if (donor.rep == GOOD) {
    if (recipient.rep == GOOD) {
      return self.bits[0]
    } else {
      return self.bits[1]
    }
  } else {
    if (recipient.rep == GOOD) {
      return self.bits[2]
    } else {
      return self.bits[3]
    }
  }
}

type Agent struct {
  rep Rep
  actm ActionModule
  payout int32
  num_games int8
}

func MakeAgent() Agent {
  var actm = MakeActionModule()
  return Agent { rep: GOOD, payout: 0, num_games: 0, actm: actm }
}

func (self Agent) ChooseDonate(recipient Agent) bool {
  return self.actm.ChooseDonate(self, recipient)
}

func (self Agent) PlayRound(recipient Agent, cost int32, benefit int32) int32 {
  // increase number of games played
  self.num_games += 1
  recipient.num_games += 1

  // keep track of total payment earned by agents
  var total_payout int32 = 0

  // play round
  if (self.ChooseDonate(recipient)) {
    // donor donates
    // -- recipient receives benefit
    recipient.payout += benefit
    // -- donor pays cost
    self.payout -= cost
    // update total payout
    total_payout += (benefit - cost)
  }

  // to prevent negative payout, each agent receives cost
  self.payout += cost
  recipient.payout += cost
  total_payout += (2*cost)

  // return total payout earned by both agents
  return total_payout
}
