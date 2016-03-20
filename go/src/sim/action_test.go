package sim

import "testing"

func TestAction(u *testing.T) {
  // CO action module
  am := NewActionModule(true, false, true, false)

  AssertTrue(u, am.ChooseDonate(GOOD, GOOD))
  AssertFalse(u, am.ChooseDonate(GOOD, BAD))
  AssertTrue(u, am.ChooseDonate(BAD, GOOD))
  AssertFalse(u, am.ChooseDonate(BAD, BAD))
}
