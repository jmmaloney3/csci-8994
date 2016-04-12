package sim

import "testing"

func TestCloneWithMutation(u *testing.T) {
  pactmut := float64(0.0)
  rnGen := NewRandNumGen()
  // CO action module
  // set mutation rate to zero
  co := NewActionModule(true, false, true, false, PEXEERR)

  clone := co.CloneWithMutations(pactmut, rnGen)
  // mutation rate is zero so modules will be the same
  AssertActModEqual(u, co, clone)

  // set mutation rate to one
  pactmut = float64(1.0)

  clone = co.CloneWithMutations(pactmut, rnGen)
  // mutation rate is one so modules will be opposites
  AssertActModOpposite(u, co, clone)
}

func TestChooseDonate(u *testing.T) {
  // CO action module
  am := NewActionModule(true, false, true, false, PEXEERR)

  rnGen := NewRandNumGen()
  AssertTrue(u, am.ChooseDonate(GOOD, GOOD, rnGen))
  AssertFalse(u, am.ChooseDonate(GOOD, BAD, rnGen))
  AssertTrue(u, am.ChooseDonate(BAD, GOOD, rnGen))
  AssertFalse(u, am.ChooseDonate(BAD, BAD, rnGen))
}

func TestActionGetBit(u *testing.T) {
  // CO action module
  am := NewActionModule(true, false, true, false, PEXEERR)

  AssertIntEqual(u, am.GetBit(0), 1)
  AssertIntEqual(u, am.GetBit(1), 0)
  AssertIntEqual(u, am.GetBit(2), 1)
  AssertIntEqual(u, am.GetBit(3), 0)
}
