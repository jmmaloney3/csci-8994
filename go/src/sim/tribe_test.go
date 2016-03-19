package sim

import "testing"

func TestInit(u *testing.T) {
  t := NewTribe(3)
  AssertInt32Equal(u, t.totalPayouts, 0)
  AssertIntEqual(u, len(t.agents), 3)
}

func TestPlayRounds(u * testing.T) {
  cost := int32(1)
  benefit := int32(3)

  t := NewTribe(3)

  // all agents use CO action model
  co := NewActionModule(true, false, true, false)
  for i := 0; i < len(t.agents); i++ {
    t.agents[i].actMod = co
  }

  // three rounds will be played
  // total payout will be 12
  AssertTrue(u, t.agents[0].ChooseDonate(t.agents[1]))
  AssertTrue(u, t.agents[0].ChooseDonate(t.agents[2]))
  AssertTrue(u, t.agents[1].ChooseDonate(t.agents[0]))
  AssertTrue(u, t.agents[1].ChooseDonate(t.agents[2]))
  AssertTrue(u, t.agents[2].ChooseDonate(t.agents[0]))
  AssertTrue(u, t.agents[2].ChooseDonate(t.agents[1]))
  AssertInt32Equal(u, t.PlayRounds(cost, benefit), 12)

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

  // three rounds will be played
  // total payout will be 0
  AssertFalse(u, t.agents[0].ChooseDonate(t.agents[1]))
  AssertFalse(u, t.agents[0].ChooseDonate(t.agents[2]))
  AssertFalse(u, t.agents[1].ChooseDonate(t.agents[0]))
  AssertFalse(u, t.agents[1].ChooseDonate(t.agents[2]))
  AssertFalse(u, t.agents[2].ChooseDonate(t.agents[0]))
  AssertFalse(u, t.agents[2].ChooseDonate(t.agents[1]))
  AssertInt32Equal(u, t.PlayRounds(cost, benefit), 6)
}

func TestSelectParent(u *testing.T) {
  t := NewTribe(3)

  allc := NewActionModule(true, true, true, true)
  alld := NewActionModule(false, false, false, false)

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
  AssertAgentEqual(u, t.SelectParent(), t.agents[2])

  // all agents in next generation will inherit from agent #2
  t.CreateNextGen()
  for i := 0; i < len(t.agents); i++ {
    AssertActModEqual(u, t.agents[i].actMod, alld)
  }
}
