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
