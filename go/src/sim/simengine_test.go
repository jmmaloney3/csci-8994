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
  // -- set assessment error to zero so that agent actions are deterministic
  passerr := float32(0)
  params[PASSE_F] = float64(passerr)
  // -- set execution error to zero so that agent actions are deterministic
  pexeerr := float32(0)
  params[PEXEE_F] = float64(pexeerr)
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
  pactmut := float32(0)
  params[PACTM_F] = pactmut

  s := NewSimEngine(2, 2, params, useMP)

  // set up tribe 0
  s.tribes[0].assessMod = NewAssessModule(GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, passerr)
  // unconditional cooperators
  s.tribes[0].agents[0].actMod = NewActionModule(true, true, true, true, pactmut, pexeerr)
  s.tribes[0].agents[1].actMod = NewActionModule(true, true, true, true, pactmut, pexeerr)

  // set up tribe 1
  s.tribes[1].assessMod = NewAssessModule(BAD, BAD, BAD, BAD, BAD, BAD, BAD, BAD, passerr)
  // unconditional defectors
  s.tribes[1].agents[0].actMod = NewActionModule(false, false, false, false, pactmut, pexeerr)
  s.tribes[1].agents[1].actMod = NewActionModule(false, false, false, false, pactmut, pexeerr)

  s.PlayRounds(cost, benefit)

  u.Logf("payouts before tribe evolution")
  for i := 0; i < s.numTribes; i++ {
    u.Logf("tribe %d totalPayouts = %d\n", i, s.tribes[i].totalPayouts)
    for j := 0; j < s.tribes[i].numAgents; j++ {
      u.Logf("  agent %d payout = %d\n", j, s.tribes[i].agents[j].payout)
    }
  }

  s.EvolveTribes()

  u.Logf("payouts after tribe evolution")
  for i := 0; i < s.numTribes; i++ {
    u.Logf("tribe %d totalPayouts = %d\n", i, s.tribes[i].totalPayouts)
    for j := 0; j < s.tribes[i].numAgents; j++ {
      u.Logf("  agent %d payout = %d\n", j, s.tribes[i].agents[j].payout)
    }
  }

  AssertAssModEqual(u, s.tribes[0].assessMod, s.tribes[1].assessMod)
  am := NewAssessModule(GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, passerr)
  AssertAssModEqual(u, s.tribes[0].assessMod, am)
  AssertAssModEqual(u, s.tribes[1].assessMod, am)
}
