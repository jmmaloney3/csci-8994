package agent

import "fmt"

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

type Tribe struct {
  agents []Agent
  numAgents int
  totalPayouts int32
}

func MakeTribe(numAgents int) Tribe {
  agents := make([]Agent, numAgents)
  // create agents
  for i := 0; i < numAgents; i++ {
    fmt.Println("    make agent ", i)
    agents[i] = MakeAgent()
  }
  return Tribe { agents: agents, numAgents: numAgents, totalPayouts: 0 }
}

type SimEngine struct {
  tribes []Tribe
  numTribes int
  totalPayouts int32
}

func MakeSimEngine(numTribes int, numAgents int) SimEngine {
  tribes := make([]Tribe, numTribes)
  // create tribes
  for i := 0; i < numTribes; i++ {
    fmt.Println("  make tribe ", i)
    tribes[i] = MakeTribe(numAgents)
  }
  return SimEngine { tribes: tribes, numTribes: numTribes, totalPayouts: 0 }
}

func (self SimEngine) PlayRounds(cost int32, benefit int32) int32 {
  var total_payout int32 = 0
  for i := 0; i < self.numTribes; i++ {
    for j := 0; j < self.tribes[i].numAgents; j++ {
      for k := j; k < self.tribes[i].numAgents; k++ {
        total_payout += self.tribes[i].agents[j].PlayRound(self.tribes[i].agents[k], 1, 3)
      }
    }
  }
  return total_payout
}
