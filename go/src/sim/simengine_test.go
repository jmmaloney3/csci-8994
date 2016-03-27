package sim

import "testing"
import "math"

func TestSim(u *testing.T) {
  runSimTest(u, false)
}

func TestSimMP(u *testing.T) {
  runSimTest(u, true)
}

func runSimTest(u *testing.T, useMP bool) {
  // create parameter map
  var params = make(map[string]float64)

  cost := int32(1)
  benefit := int32(3)

  // populate arg maps
  // -- use default assessment error probability
  params[PASSE_F] = PASSERR
  // -- use default execution error probability
  params[PEXEE_F] = PEXEERR
  // -- set probability of conflict to 1 so tribes always evolve
  params[PCON_F]  = float64(1)
  // -- set Beta to infinity so tribe with largest payout always wins
  params[BETA_F]  = math.Inf(int(1))
  // -- set eta to 1 so that loser bits are always shifted to winner bits
  params[ETA_F]   = float64(1)
  // -- set migration probability to???
  params[PMIG_F]  = PMIG
  // -- set probability of mutation to zero so that winner bits are copied faithfully
  params[PASSM_F] = float64(0)
  // -- set mutation probability to zero so only fittest stratgies are copied forward
  params[PACTM_F] = float64(0)

  s := NewSimEngine(2, 2, params, useMP)

  // set up tribe 0
  s.tribes[0].assessMod = NewAssessModule(GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, PASSERR)
  // unconditional cooperators
  s.tribes[0].agents[0].actMod = NewActionModule(true, true, true, true, PEXEERR)
  s.tribes[0].agents[1].actMod = NewActionModule(true, true, true, true, PEXEERR)

  // set up tribe 1
  s.tribes[1].assessMod = NewAssessModule(BAD, BAD, BAD, BAD, BAD, BAD, BAD, BAD, PASSERR)
  // unconditional defectors
  s.tribes[1].agents[0].actMod = NewActionModule(false, false, false, false, PEXEERR)
  s.tribes[1].agents[1].actMod = NewActionModule(false, false, false, false, PEXEERR)

  s.PlayRounds(cost, benefit)
  for i := 0; i < s.numTribes; i++ {
    u.Logf("tribe %d totalPayouts = %d\n", i, s.tribes[i].totalPayouts)
    for j := 0; j < s.tribes[i].numAgents; j++ {
      u.Logf("  agent %d payout = %d\n", j, s.tribes[i].agents[j].payout)
    }
  }

  s.EvolveTribes()

  AssertTrue(u, s.tribes[0].assessMod.SameBits(s.tribes[1].assessMod))
  am := NewAssessModule(GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, PASSERR)
  AssertTrue(u, s.tribes[0].assessMod.SameBits(am))
  AssertTrue(u, s.tribes[1].assessMod.SameBits(am))
}
