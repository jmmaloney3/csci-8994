package sim

import "testing"
import "math"

func TestEvolveTribes(u *testing.T) {
  cost := int32(1)
  benefit := int32(3)

  s := NewSimEngine(2,2,false)
  // set probability of conflict to 1 so tribes always evolve
  s.conP = float32(1)
  // set Beta to infinity so tribe with largest payout always wins
  s.Beta = math.Inf(int(1))
  // set eta to 1 so that loser bits are always shifted to winner bits
  s.eta = float64(1)
  // set migration probability to???
  // s.migP = ??
  // set probability of mutation to zero so that winner bits are copied faithfully
  s.mutP = float32(0)

  // set up tribes
  for i := 0; i < s.numTribes; i++ {
    // set mutation probability to zero so only fittest stratgies are copied forward
    s.tribes[i].mutP = float32(0)
  }

  // set up tribe 0
  s.tribes[0].assessMod = NewAssessModule(GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, GOOD)
  // unconditional cooperators
  s.tribes[0].agents[0].actMod = NewActionModule(true, true, true, true)
  s.tribes[0].agents[1].actMod = NewActionModule(true, true, true, true)

  // set up tribe 1
  s.tribes[1].assessMod = NewAssessModule(BAD, BAD, BAD, BAD, BAD, BAD, BAD, BAD)
  // unconditional defectors
  s.tribes[1].agents[0].actMod = NewActionModule(false, false, false,false)
  s.tribes[1].agents[1].actMod = NewActionModule(false, false, false,false)

  s.PlayRounds(cost, benefit)
  for i := 0; i < s.numTribes; i++ {
    u.Logf("tribe %d totalPayouts = %d\n", i, s.tribes[i].totalPayouts)
    for j := 0; j < s.tribes[i].numAgents; j++ {
      u.Logf("  agent %d payout = %d\n", j, s.tribes[i].agents[j].payout)
    }
  }

  s.EvolveTribes()

  AssertTrue(u, s.tribes[0].assessMod.SameBits(s.tribes[1].assessMod))
  am := NewAssessModule(GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, GOOD)
  AssertTrue(u, s.tribes[0].assessMod.SameBits(am))
  AssertTrue(u, s.tribes[1].assessMod.SameBits(am))
}
