package sim

import "testing"
import "math"
import "math/rand"

func TestNewSimEngine(u *testing.T) {
  numTribes := 2
  numAgents := 2
  useMP := true
  s := NewDefaultSimEngine(numTribes, numAgents, useMP)

  // check that engine was created correctly
  AssertIntEqual(u, s.numTribes, numTribes)
  var t *Tribe
  var a *Agent
  for i := 0; i < numTribes; i++ {
    t = s.tribes[i]
    AssertFalse(u, t.assessMod == nil)
    AssertIntEqual(u, t.numAgents, numAgents)
    AssertInt32Equal(u, t.totalPayouts, 0)
    for j := 0; j < numAgents; j++ {
      a = t.agents[j]
      AssertTrue(u, a.tribe == t)
      AssertRepEqual(u, a.rep, GOOD)
      AssertFalse(u, a.actMod == nil)
      AssertInt32Equal(u, a.payout, 0)
      AssertInt8Equal(u, a.numGames, 0)
    }
  }
  AssertInt32Equal(u, s.totalPayouts, 0)
  AssertFalse(u, s.rnGen == nil)
  AssertTrue(u, s.useMP == useMP)
  AssertFloat32Equal(u, s.pcon, PCON)
  AssertFloat64Equal(u, s.beta, BETA)
  AssertFloat64Equal(u, s.eta, ETA)
  AssertFloat32Equal(u, s.pmig, PMIG)
  AssertFloat32Equal(u, s.passmut, PASSMUT)
  alld := NewActionModule(false, false, false, false, PACTMUT, PEXEERR)
  AssertActModEqual(u, s.ALLD, alld)
  allc := NewActionModule(true, true, true, true, PACTMUT, PEXEERR)
  AssertActModEqual(u, s.ALLC, allc)
}

func TestConflict(u *testing.T) {
  rnGen := NewRandNumGen()
  numTribes := 2
  numAgents := 2
  useMP := true
  s := NewDefaultSimEngine(numTribes, numAgents, useMP)

  s.tribes[0].totalPayouts = 5
  s.tribes[1].totalPayouts = 10

  // test deterministic winner determination
  s.beta = math.Inf(int(1))
  w, l := s.Conflict(1,0,rnGen)
  AssertIntEqual(u, w, 1)
  AssertIntEqual(u, l, 0)
  // flip order of tribes and make sure answer is the same
  w, l = s.Conflict(0,1,rnGen)
  AssertIntEqual(u, w, 1)
  AssertIntEqual(u, l, 0)

  // non-deterministic winner determination
  s.beta = float64(1)
  w, l = s.Conflict(1,0,rnGen)
  AssertTrue(u, (w ==0) || (w == 1))
  AssertTrue(u, (l ==0) || (l == 1))
}

func TestShiftAssessMod(u *testing.T) {
  rnGen := NewRandNumGen()
  numTribes := 2
  numAgents := 2
  useMP := true
  s := NewDefaultSimEngine(numTribes, numAgents, useMP)

  // all GOOD
  allg := NewAssessModule(GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, PASSERR)
  // all BAD
  allb := NewAssessModule(BAD, BAD, BAD, BAD, BAD, BAD, BAD, BAD, PASSERR)

  AssertAssModOpposite(u, allb, allg)

  // first test with tribe payouts equal to zero
  AssertInt32Equal(u, s.tribes[0].totalPayouts, 0)
  AssertInt32Equal(u, s.tribes[1].totalPayouts, 0)

  testSAM(u, s, allb, allg, rnGen)

  // test with positive payouts
  s.tribes[0].totalPayouts = 50
  s.tribes[1].totalPayouts = 100

  testSAM(u, s, allb, allg, rnGen)
}

func testSAM(u *testing.T, s *SimEngine, allb *AssessModule, allg *AssessModule, rnGen *rand.Rand) {
  // set eta to 1 so that loser bits are always shifted to winner bits
  s.eta = float64(1)

  // set assessment modules
  s.tribes[0].assessMod = allb
  s.tribes[1].assessMod = allg

  // assume that tribe 0 is the winner
  s.ShiftAssessMod(s.tribes[0], s.tribes[1], rnGen)
  AssertAssModEqual(u, s.tribes[0].assessMod, allb)
  AssertAssModEqual(u, s.tribes[1].assessMod, allb)

  // reset assessment modules
  s.tribes[0].assessMod = allb
  s.tribes[1].assessMod = allg

  // assume that tribe 1 is the winner
  s.ShiftAssessMod(s.tribes[1], s.tribes[0], rnGen)
  AssertAssModEqual(u, s.tribes[0].assessMod, allg)
  AssertAssModEqual(u, s.tribes[1].assessMod, allg)

  // set eta to 0 so that loser bits are never shifted to winner bits
  s.eta = float64(0)

  // set assessment modules
  s.tribes[0].assessMod = allb
  s.tribes[1].assessMod = allg

  // assume that tribe 0 is the winner
  s.ShiftAssessMod(s.tribes[0], s.tribes[1], rnGen)
  AssertAssModEqual(u, s.tribes[0].assessMod, allb)
  AssertAssModEqual(u, s.tribes[1].assessMod, allg)

  // reset assessment modules
  s.tribes[0].assessMod = allb
  s.tribes[1].assessMod = allg

  // assume that tribe 1 is the winner
  s.ShiftAssessMod(s.tribes[1], s.tribes[0], rnGen)
  AssertAssModEqual(u, s.tribes[0].assessMod, allb)
  AssertAssModEqual(u, s.tribes[1].assessMod, allg)
}

func TestMigrateAgents(u *testing.T) {
  rnGen := NewRandNumGen()
  numTribes := 2
  numAgents := 2
  useMP := true
  s := NewDefaultSimEngine(numTribes, numAgents, useMP)

  // tribe zero has unconditional cooperators
  // tribe one has unconditional defectors
  for i := 0; i < numAgents; i++ {
    s.tribes[0].agents[i].actMod = s.ALLC
    s.tribes[1].agents[i].actMod = s.ALLD
  }

  // set probability of migration to one
  s.pmig = float32(1)

  // test migrate agents from tribe 0 to tribe 1
  s.MigrateAgents(s.tribes[0], s.tribes[1], rnGen)

  // verify that all agents are now using ALLC
  for i := 0; i < numAgents; i++ {
    AssertActModEqual(u, s.tribes[0].agents[i].actMod, s.ALLC)
    AssertActModEqual(u, s.tribes[1].agents[i].actMod, s.ALLC)
  }

  // reset to original state
  // tribe zero has unconditional cooperators
  // tribe one has unconditional defectors
  for i := 0; i < numAgents; i++ {
    s.tribes[0].agents[i].actMod = s.ALLC
    s.tribes[1].agents[i].actMod = s.ALLD
  }

  // set probability of migration to zero
  s.pmig = float32(0)

  // test migrate agents from tribe 1 to tribe 0-
  s.MigrateAgents(s.tribes[1], s.tribes[0], rnGen)

  // verify that no miration occured
  for i := 0; i < numAgents; i++ {
    AssertActModEqual(u, s.tribes[0].agents[i].actMod, s.ALLC)
    AssertActModEqual(u, s.tribes[1].agents[i].actMod, s.ALLD)
  }
}

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
  // -- set migration probability to one so that migration always occurs
  params[PMIG_F]  = float64(1)
  // -- set probability of mutation to zero so that winner bits are copied faithfully
  params[PASSM_F] = float64(0)
  // -- set mutation probability to zero so only fittest stratgies are copied forward
  pactmut := float32(0)
  params[PACTM_F] = float64(pactmut)

  s := NewSimEngine(2, 2, params, useMP)

  allg := NewAssessModule(GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, GOOD, passerr)
  allb := NewAssessModule(BAD, BAD, BAD, BAD, BAD, BAD, BAD, BAD, passerr)

  allc := NewActionModule(true, true, true, true, pactmut, pexeerr)
  alld := NewActionModule(false, false, false, false, pactmut, pexeerr)

  // set up tribe 0
  s.tribes[0].assessMod = allg
  // unconditional cooperators
  s.tribes[0].agents[0].actMod = allc
  s.tribes[0].agents[1].actMod = allc

  // set up tribe 1
  s.tribes[1].assessMod = allb
  // unconditional defectors
  s.tribes[1].agents[0].actMod = alld
  s.tribes[1].agents[1].actMod = alld

  nextGen := s.PlayRounds(cost, benefit)

  // check payouts
  tp0 := s.tribes[0].agents[0].payout + s.tribes[0].agents[1].payout
  AssertInt32Equal(u, s.tribes[0].totalPayouts, tp0)
  AssertInt32Equal(u, tp0, 4)
  tp1 := s.tribes[1].agents[0].payout + s.tribes[1].agents[1].payout
  AssertInt32Equal(u, s.tribes[1].totalPayouts, tp1)
  AssertInt32Equal(u, tp1, 2)

  s.EvolveTribes(nextGen)

  u.Logf("payouts after tribe evolution")
  for i := 0; i < s.numTribes; i++ {
    u.Logf("tribe %d totalPayouts = %d\n", i, s.tribes[i].totalPayouts)
    for j := 0; j < s.tribes[i].numAgents; j++ {
      u.Logf("  agent %d payout = %d\n", j, s.tribes[i].agents[j].payout)
    }
  }

  // Since pcon is one, the tribes always engage in conflict
  // Since beta is infinite, tribe with highest payout always wins (tribe[0])
  // Since eta is one, loser assessment module bits are always flipped to winner bits
  // Since passmut is zero, matching assessment module bits are never flipped
  AssertAssModEqual(u, s.tribes[0].assessMod, s.tribes[1].assessMod)
  AssertAssModEqual(u, s.tribes[0].assessMod, allg)
  AssertAssModEqual(u, s.tribes[1].assessMod, allg)

  // Since pmig is set to one, winner agents always migrate to loser tribe
  for i := 0; i < 2; i++ {
    AssertActModEqual(u, s.tribes[0].agents[i].actMod, allc)
    AssertActModEqual(u, s.tribes[1].agents[i].actMod, allc)
  }
  for i := 0; i < 2; i++ {
    AssertActModEqual(u, nextGen[0].agents[i].actMod, allc)
    AssertActModEqual(u, nextGen[1].agents[i].actMod, allc)
  }

  // set passmut to one so that matching assessment module bits are always flipped
  s.passmut = float32(1)

  // restore tribe 1 agents to unconditional defectors
  s.tribes[1].agents[0].actMod = alld
  s.tribes[1].agents[1].actMod = alld

  // set pmig to zero so that migration never occurs
  s.pmig = float32(0)

  s.EvolveTribes(nextGen)

  // Since passmut is one, matching assessment module bits are always flipped
  AssertAssModEqual(u, s.tribes[0].assessMod, allg)
  AssertAssModEqual(u, s.tribes[1].assessMod, allb)

  // Since pmig is set to zero, winner agents never migrate to loser tribe
  for i := 0; i < 2; i++ {
    AssertActModEqual(u, s.tribes[0].agents[i].actMod, allc)
    AssertActModEqual(u, s.tribes[1].agents[i].actMod, alld)
  }

  // set eta to zero so that loser assessment module bits are never flipped to winner bits
  s.eta = float64(0)

  // set pmig to one so that migration always occurs
  s.pmig = float32(1)

  s.EvolveTribes(nextGen)

  // Passmut doesn't come into play since the two assessment modules have no matching bits
  // Since eta is one, loser's assessment module bits are always flipped
  AssertAssModEqual(u, s.tribes[0].assessMod, allg)
  AssertAssModEqual(u, s.tribes[1].assessMod, allb)

  // Since pmig is set to one, winner agents always migrate to loser tribe
  for i := 0; i < 2; i++ {
    AssertActModEqual(u, s.tribes[0].agents[i].actMod, allc)
    AssertActModEqual(u, s.tribes[1].agents[i].actMod, allc)
  }
}
