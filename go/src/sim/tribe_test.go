package sim

import "testing"
import "sort"

func TestNewTribe(u *testing.T) {
  numAgents := 3
  t := NewTribe(numAgents, PASSERR, PACTMUT, PEXEERR, NewRandNumGen())
  AssertIntEqual(u, t.numAgents, numAgents)
  AssertInt32Equal(u, t.totalPayouts, 0)
  AssertIntEqual(u, len(t.agents), numAgents)
  var agent *Agent
  for i := 0; i < numAgents; i++ {
    agent = t.agents[i]
    AssertTrue(u, agent.tribe == t)
    AssertRepEqual(u, agent.rep, GOOD)
    AssertInt32Equal(u, agent.payout, 0)
    AssertInt8Equal(u, agent.numGames, 0)
  }
}

func TestPlayRounds(u * testing.T) {
  cost := int32(1)
  benefit := int32(3)
  rnGen := NewRandNumGen()
  // make agent actions deterministic
  pexeerr := float32(0)
  passerr := float32(0)

  t := NewTribe(3, passerr, PACTMUT, pexeerr, rnGen)
  t.assessMod = NewAssessModule(GOOD, BAD, BAD, GOOD, GOOD, BAD, BAD, GOOD, passerr)

  // all agents use CO action model
  co := NewActionModule(true, false, true, false, PACTMUT, pexeerr)
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

func TestPlayRounds2(u *testing.T) {
  cost := int32(1)
  benefit := int32(3)
  rnGen := NewRandNumGen()
  // make agent actions deterministic
  pexeerr := float32(0)
  passerr := float32(0)

  t := NewTribe(3, passerr, PACTMUT, pexeerr, rnGen)
  // stern-judging assessment module
  t.assessMod = NewAssessModule(GOOD, BAD, BAD, GOOD, GOOD, BAD, BAD, GOOD, passerr)

  // all agents use CO action model
  co := NewActionModule(true, false, true, false, PACTMUT, pexeerr)
  for i := 0; i < len(t.agents); i++ {
    t.agents[i].actMod = co
  }

  // set agent reps to BAD
  for i := 0; i < len(t.agents); i++ {
    t.agents[i].rep = BAD
  }

  // play rounds using stern-judging and CO
  t.PlayRounds(cost, benefit, rnGen)

  // count number of good agents
  numGood := 0
  for i := 0; i < 3; i++ {
    if (t.agents[i].rep == GOOD) { numGood++ }
  }
  AssertTrue(u, (numGood == 2) || (numGood == 3))

  // check results
  tp := t.totalPayouts
  switch {
  case (numGood == 2):
    AssertTrue(u, (tp == 6) || (tp == 8))
  case (numGood == 3):
    AssertTrue(u, (tp ==8) || (tp == 10))
  }
}

func TestSelectParent(u *testing.T) {
  rnGen := NewRandNumGen()

  // set tribe's probability of action module bit mutation to zero
  pactmut := float32(0)
  t := NewTribe(3, PASSERR, pactmut, PEXEERR, rnGen)

  allc := NewActionModule(true, true, true, true, pactmut, PEXEERR)
  alld := NewActionModule(false, false, false, false, pactmut, PEXEERR)

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
  // -- probability of action module bit mutation was set to zero above
  nextGen := t.CreateNextGen(rnGen)
  AssertAssModEqual(u, t.assessMod, nextGen.assessMod)
  AssertIntEqual(u, t.numAgents, nextGen.numAgents)
  AssertInt32Equal(u, nextGen.totalPayouts, 0)
  for i := 0; i < len(nextGen.agents); i++ {
    agent := nextGen.agents[i]
    AssertTrue(u, agent.tribe == nextGen)
    AssertRepEqual(u, agent.rep, GOOD)
    AssertActModEqual(u, agent.actMod, alld)
    AssertInt32Equal(u, agent.payout, 0)
    AssertInt8Equal(u, agent.numGames, 0)
  }

  // set tribe's probability of action module bit mutation to one
  for i := 0; i < t.numAgents; i++ {
    t.agents[i].actMod.pactmut = float32(1)
  }

  // all agents in next generation will inherit from agent #2 again
  // because payouts have not changed, but the probability of mutation
  // is now 1 so the child agent's modules will be allc
  nextGen = t.CreateNextGen(rnGen)
  AssertAssModEqual(u, t.assessMod, nextGen.assessMod)
  AssertIntEqual(u, t.numAgents, nextGen.numAgents)
  AssertInt32Equal(u, nextGen.totalPayouts, 0)
  for i := 0; i < len(nextGen.agents); i++ {
    agent := nextGen.agents[i]
    AssertTrue(u, agent.tribe == nextGen)
    AssertRepEqual(u, agent.rep, GOOD)
    AssertActModEqual(u, agent.actMod, allc)
    AssertInt32Equal(u, agent.payout, 0)
    AssertInt8Equal(u, agent.numGames, 0)
  }
}

func TestSortByPayout(u *testing.T) {
  rnGen := NewRandNumGen()
  numTribes := 50
  numAgents := 2
  tribes := make([]*Tribe, numTribes)
  payouts := rnGen.Perm(numTribes)
  for i := 0; i < numTribes; i++ {
    tribes[i] = NewTribe(numAgents, PASSERR, PACTMUT, PEXEERR, rnGen)
    tribes[i].totalPayouts = int32(payouts[i])
  }
  // sort the tribes by their payouts
  sort.Sort(SortTribesByPayouts(tribes))
  // test that the tribes are sorted correctly
  current := int32(-1)
  for i := 0; i < numTribes; i++ {
    AssertInt32GT(u, tribes[i].totalPayouts, current)
    current = tribes[i].totalPayouts
  }
}
