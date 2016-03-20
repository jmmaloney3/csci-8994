package sim

import "testing"

func TestAssess(u *testing.T) {
  // stern judging assessment module
  am := NewAssessModule(GOOD, BAD, BAD, GOOD, GOOD, BAD, BAD, GOOD)

  AssertRepEqual(u, am.AssignRep(GOOD, GOOD, DONATE), GOOD)
  AssertRepEqual(u, am.AssignRep(GOOD, GOOD, REFUSE), BAD)
  AssertRepEqual(u, am.AssignRep(GOOD, BAD,  DONATE), BAD)
  AssertRepEqual(u, am.AssignRep(GOOD, BAD,  REFUSE), GOOD)
  AssertRepEqual(u, am.AssignRep(BAD,  GOOD, DONATE), GOOD)
  AssertRepEqual(u, am.AssignRep(BAD,  GOOD, REFUSE), BAD)
  AssertRepEqual(u, am.AssignRep(BAD,  BAD,  DONATE), BAD)
  AssertRepEqual(u, am.AssignRep(BAD,  BAD,  REFUSE), GOOD)
}
