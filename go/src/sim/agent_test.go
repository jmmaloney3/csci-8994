package sim

import "testing"

func TestNewAgent(u *testing.T) {
  t := NewTribe(1, PASSERR, PACTMUT, PEXEERR, NewRandNumGen())
  a := t.agents[0]
  AssertRepEqual(u, a.rep, GOOD)
  AssertInt32Equal(u, a.payout, 0)
  AssertInt8Equal(u, a.numGames, 0)
}

func TestCreateChild(u *testing.T) {
  rnGen := NewRandNumGen()
  passerr := float32(0)
  pactmut := float32(0)
  pexeerr := float32(0)
  // create original tribe
  t1 := NewTribe(1, passerr, pactmut, pexeerr, rnGen)
  a1 := t1.agents[0]

  // create next generation
  t2 := NewTribe(1, passerr, pactmut, pexeerr, rnGen)
  a2 := a1.CreateChild(t2, rnGen)
  AssertTrue(u, a1.tribe == t1)
  AssertTrue(u, a2.tribe == t2)
  // since pactmut==0, action modules will be the same
  AssertActModEqual(u, a2.actMod, a1.actMod)
  AssertInt32Equal(u, a2.payout, 0)
  AssertInt8Equal(u, a2.numGames, 0)

  // create new original tribe with pactmut = 1.0
  pactmut = 1.0
  t1 = NewTribe(1, passerr, pactmut, pexeerr, rnGen)
  a1 = t1.agents[0]

  // create another generation
  t2 = NewTribe(1, passerr, pactmut, pexeerr, rnGen)
  a2 = a1.CreateChild(t2, rnGen)
  AssertTrue(u, a1.tribe == t1)
  AssertTrue(u, a2.tribe == t2)
  // since pactmut==1, action modules will be opposites
  AssertActModOpposite(u, a1.actMod, a2.actMod)
  AssertInt32Equal(u, a2.payout, 0)
  AssertInt8Equal(u, a2.numGames, 0)
}

func TestPlayround(t *testing.T) {
  cost := int32(1)
  benefit := int32(3)
  rnGen := NewRandNumGen()
  passerr := float32(0)
  pactmut := float32(0)
  pexeerr := float32(0)

  u := NewTribe(2, passerr, pactmut, pexeerr, rnGen)
  // configure tribe assess module
  u.assessMod = NewAssessModule(GOOD, BAD, BAD, GOOD, GOOD, BAD, BAD, GOOD, passerr)

  don := u.agents[0]
  AssertRepEqual(t, don.rep, GOOD)
  AssertInt32Equal(t, don.payout, 0)
  rec := u.agents[1]
  AssertRepEqual(t, rec.rep, GOOD);
  AssertInt32Equal(t, rec.payout, 0)

  // configure donor action module
  don.actMod = NewActionModule(true, false, true, true, pactmut, pexeerr);

  // GOOD GOOD
  AssertTrue(t, don.ChooseDonate(rec, rnGen))
  AssertInt32Equal(t, don.PlayRound(rec, cost, benefit, rnGen), benefit-cost+2*cost);
  AssertInt32Equal(t, don.payout, 0)
  AssertInt8Equal(t, don.numGames, 1)
  AssertRepEqual(t, don.rep, GOOD)
  AssertInt32Equal(t, rec.payout, 4)
  AssertInt8Equal(t, rec.numGames, 1)
  AssertRepEqual(t, rec.rep, GOOD)

  // GOOD BAD
  rec.rep = BAD
  AssertFalse(t, don.ChooseDonate(rec, rnGen))
  AssertInt32Equal(t, don.PlayRound(rec, cost, benefit, rnGen), 2*cost);
  AssertInt32Equal(t, don.payout, 1)
  AssertInt8Equal(t, don.numGames, 2)
  AssertRepEqual(t, don.rep, GOOD)
  AssertInt32Equal(t, rec.payout, 5)
  AssertInt8Equal(t, rec.numGames, 2)
  AssertRepEqual(t, rec.rep, BAD)

  // BAD BAD
  don.rep = BAD
  AssertTrue(t, don.ChooseDonate(rec, rnGen))
  AssertInt32Equal(t, don.PlayRound(rec, cost, benefit, rnGen), benefit-cost+2*cost);
  AssertInt32Equal(t, don.payout, 1)
  AssertInt8Equal(t, don.numGames, 3)
  AssertRepEqual(t, don.rep, BAD)
  AssertInt32Equal(t, rec.payout, 9)
  AssertInt8Equal(t, rec.numGames, 3)
  AssertRepEqual(t, rec.rep, BAD)

  // BAD GOOD
  rec.rep = GOOD
  AssertRepEqual(t, don.rep, BAD)
  AssertTrue(t, don.ChooseDonate(rec, rnGen))
  AssertInt32Equal(t, don.PlayRound(rec, cost, benefit, rnGen), benefit-cost+2*cost);
  AssertInt32Equal(t, don.payout, 1)
  AssertInt8Equal(t, don.numGames, 4)
  AssertRepEqual(t, don.rep, GOOD)
  AssertInt32Equal(t, rec.payout, 13)
  AssertInt8Equal(t, rec.numGames, 4)
  AssertRepEqual(t, rec.rep, GOOD)

  // reset
  don.Reset()
  AssertInt32Equal(t, don.payout, 0)
  AssertInt8Equal(t, don.numGames, 0)
  rec.Reset()
  AssertInt32Equal(t, rec.payout, 0)
  AssertInt8Equal(t, rec.numGames, 0)
}
