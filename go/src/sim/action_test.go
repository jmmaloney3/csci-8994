package sim

import "testing"

func TestAction(u *testing.T) {
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
