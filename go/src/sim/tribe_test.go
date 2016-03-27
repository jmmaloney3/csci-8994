package sim

import "testing"

func TestInit(u *testing.T) {
  t := NewTribe(3, PASSERR, PACTMUT, PEXEERR, NewRandNumGen())
  AssertInt32Equal(u, t.totalPayouts, 0)
  AssertIntEqual(u, len(t.agents), 3)
}

func TestPlayRounds(u * testing.T) {
  cost := int32(1)
  benefit := int32(3)
  rnGen := NewRandNumGen()

  t := NewTribe(3, PASSERR, PACTMUT, PEXEERR, rnGen)
  t.assessMod = NewAssessModule(GOOD, BAD, BAD, GOOD, GOOD, BAD, BAD, GOOD, PASSERR)

  // all agents use CO action model
  co := NewActionModule(true, false, true, false, PEXEERR)
  for i := 0; i < len(t.agents); i++ {
    t.agents[i].actMod = co
  }

  // three rounds will be played
  // total payout will be 12
  AssertTrue(u, t.agents[0].ChooseDonate(t.agents[1], rnGen))
  AssertTrue(u, t.agents[0].ChooseDonate(t.agents[2], rnGen))
  AssertTrue(u, t.agents[1].ChooseDonate(t.agents[0], rnGen))
  AssertTrue(u, t.agents[1].ChooseDonate(t.agents[2], rnGen))
  AssertTrue(u, t.agents[2].ChooseDonate(t.agents[0], rnGen))
  AssertTrue(u, t.agents[2].ChooseDonate(t.agents[1], rnGen))
  AssertInt32Equal(u, t.PlayRounds(cost, benefit, rnGen), 12)

  // test reset
  t.Reset()
  AssertInt32Equal(u, t.totalPayouts, 0);
  for i := 0; i < len(t.agents); i++ {
    AssertInt32Equal(u, t.agents[i].payout, 0)
  }

  // set agent reps to BAD
  for i := 0; i < len(t.agents); i++ {
    t.agents[i].rep = BAD
  }

  // Agent 0 refuses to donate to agent 1 because [BAD, BAD] => REFUSE
  AssertFalse(u, t.agents[0].ChooseDonate(t.agents[1], rnGen))
  AssertInt32Equal(u, t.agents[0].PlayRound(t.agents[1], cost, benefit, rnGen), int32(2))
  // Agent 0 is now good because [BAD, BAD, REFUSE] => GOOD
  AssertRepEqual(u, t.agents[0].rep, GOOD)
  // Agent 0 refuses to donate to agent 2 because [GOOD, BAD] => REFUSE
  AssertFalse(u, t.agents[0].ChooseDonate(t.agents[2], rnGen))
  AssertInt32Equal(u, t.agents[0].PlayRound(t.agents[2], cost, benefit, rnGen), int32(2))
  // Agent 0 is still good because [GOOD, BAD, REFUSE] => GOOD
  AssertRepEqual(u, t.agents[0].rep, GOOD)
  // Agent 1 donates to agent 0 because [BAD, GOOD] => DONATE
  AssertTrue(u, t.agents[1].ChooseDonate(t.agents[0], rnGen))
  AssertInt32Equal(u, t.agents[1].PlayRound(t.agents[0], cost, benefit, rnGen), int32(4))
  // Agent 1 is now good because [BAD, GOOD, DONATE] => GOOD
  AssertRepEqual(u, t.agents[1].rep, GOOD)
  // Agent 1 refuses to donate to agent 2 because [GOOD, BAD] => REFUSE
  AssertFalse(u, t.agents[1].ChooseDonate(t.agents[2], rnGen))
  AssertInt32Equal(u, t.agents[1].PlayRound(t.agents[2], cost, benefit, rnGen), int32(2))
  // Agent 1 is still good because [GOOD, BAD, REFUSE] => GOOD
  AssertRepEqual(u, t.agents[1].rep, GOOD)
  // Agent 2 donates to agent 0 because [BAD, GOOD] => DONATE
  AssertTrue(u, t.agents[2].ChooseDonate(t.agents[0], rnGen))
  AssertInt32Equal(u, t.agents[2].PlayRound(t.agents[0], cost, benefit, rnGen), int32(4))
  // Agent 2 is now good because [BAD, GOOD, DONATE] => GOOD
  AssertRepEqual(u, t.agents[2].rep, GOOD)
  // Agent 2 donates to agent 0 because [GOOD, GOOD] => DONATE
  AssertTrue(u, t.agents[2].ChooseDonate(t.agents[1], rnGen))
  AssertInt32Equal(u, t.agents[2].PlayRound(t.agents[1], cost, benefit, rnGen), int32(4))
  // Agent 2 is still good because [GOOD, GOOD, DONATE] => GOOD
  AssertRepEqual(u, t.agents[2].rep, GOOD)
  // ALl agents are good, so each round resuls in a donation
  AssertInt32Equal(u, t.PlayRounds(cost, benefit, rnGen), 12)
}

func TestSelectParent(u *testing.T) {
  rnGen := NewRandNumGen()

  t := NewTribe(3, PASSERR, PACTMUT, PEXEERR, rnGen)

  allc := NewActionModule(true, true, true, true, PEXEERR)
  alld := NewActionModule(false, false, false, false, PEXEERR)

  t.agents[0].payout = -1
  t.agents[0].actMod = allc
  t.agents[1].payout = -1
  t.agents[1].actMod = allc
  t.agents[2].payout = 10
  t.agents[2].actMod = alld

  // update the tribe total payouts
  for i := 0; i < len(t.agents); i++ {
    t.totalPayouts += t.agents[i].payout
  }

  // agent with positive payout will be selected
  AssertAgentEqual(u, t.SelectParent(rnGen), t.agents[2])

  // all agents in next generation will inherit from agent #2
  t.CreateNextGen(rnGen)
  for i := 0; i < len(t.agents); i++ {
    AssertActModEqual(u, t.agents[i].actMod, alld)
  }
}
