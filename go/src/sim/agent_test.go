package sim

import "testing"

func TestNewAgent(u *testing.T) {
  t := NewTribe(1, NewRandNumGen())
  a := t.agents[0]
  AssertRepEqual(u, a.rep, GOOD)
  AssertInt32Equal(u, a.payout, 0)
  AssertInt8Equal(u, a.numGames, 0)
}

func TestPlayround(t *testing.T) {
  cost := int32(1)
  benefit := int32(3)
  rnGen := NewRandNumGen()

  u := NewTribe(2, rnGen)
  // configure tribe assess module
  u.assessMod = NewAssessModule(GOOD, BAD, BAD, GOOD, GOOD, BAD, BAD, GOOD)

  don := u.agents[0]
  AssertRepEqual(t, don.rep, GOOD)
  AssertInt32Equal(t, don.payout, 0)
  rec := u.agents[1]
  AssertRepEqual(t, rec.rep, GOOD);
  AssertInt32Equal(t, rec.payout, 0)

  // configure donor action module
  don.actMod = NewActionModule(true, false, true, true);

  // GOOD GOOD
  AssertTrue(t, don.ChooseDonate(rec, rnGen))
  AssertInt32Equal(t, don.PlayRound(rec, cost, benefit, rnGen), benefit-cost+2*cost);
  AssertInt32Equal(t, don.payout, 0)
  AssertInt8Equal(t, don.numGames, 1)
  AssertRepEqual(t, don.rep, GOOD)
  AssertInt32Equal(t, rec.payout, 4)
  AssertInt8Equal(t, rec.numGames, 1)

  // GOOD BAD
  rec.rep = BAD
  AssertFalse(t, don.ChooseDonate(rec, rnGen))
  AssertInt32Equal(t, don.PlayRound(rec, cost, benefit, rnGen), 2*cost);
  AssertInt32Equal(t, don.payout, 1)
  AssertInt8Equal(t, don.numGames, 2)
  AssertRepEqual(t, don.rep, GOOD)
  AssertInt32Equal(t, rec.payout, 5)
  AssertInt8Equal(t, rec.numGames, 2)

  // BAD BAD
  don.rep = BAD
  AssertTrue(t, don.ChooseDonate(rec, rnGen))
  AssertInt32Equal(t, don.PlayRound(rec, cost, benefit, rnGen), benefit-cost+2*cost);
  AssertInt32Equal(t, don.payout, 1)
  AssertInt8Equal(t, don.numGames, 3)
  AssertRepEqual(t, don.rep, BAD)
  AssertInt32Equal(t, rec.payout, 9)
  AssertInt8Equal(t, rec.numGames, 3)

  // BAD GOOD
  rec.rep = GOOD
  AssertTrue(t, don.ChooseDonate(rec, rnGen))
  AssertInt32Equal(t, don.PlayRound(rec, cost, benefit, rnGen), benefit-cost+2*cost);
  AssertInt32Equal(t, don.payout, 1)
  AssertInt8Equal(t, don.numGames, 4)
  AssertRepEqual(t, don.rep, GOOD)
  AssertInt32Equal(t, rec.payout, 13)
  AssertInt8Equal(t, rec.numGames, 4)

  // reset
  don.Reset()
  AssertInt32Equal(t, don.payout, 0)
  AssertInt8Equal(t, don.numGames, 0)
  rec.Reset()
  AssertInt32Equal(t, rec.payout, 0)
  AssertInt8Equal(t, rec.numGames, 0)
}
