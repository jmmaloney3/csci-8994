package sim

import "fmt"

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
