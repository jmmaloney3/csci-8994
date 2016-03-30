package sim

import "testing"

func TestAssess(u *testing.T) {
  rnGen := NewRandNumGen()

  // stern judging assessment module
  am := NewAssessModule(GOOD, BAD, BAD, GOOD, GOOD, BAD, BAD, GOOD, 0)

  AssertRepEqual(u, am.AssignRep(GOOD, GOOD, DONATE, rnGen), GOOD)
  AssertRepEqual(u, am.AssignRep(GOOD, GOOD, REFUSE, rnGen), BAD)
  AssertRepEqual(u, am.AssignRep(GOOD, BAD,  DONATE, rnGen), BAD)
  AssertRepEqual(u, am.AssignRep(GOOD, BAD,  REFUSE, rnGen), GOOD)
  AssertRepEqual(u, am.AssignRep(BAD,  GOOD, DONATE, rnGen), GOOD)
  AssertRepEqual(u, am.AssignRep(BAD,  GOOD, REFUSE, rnGen), BAD)
  AssertRepEqual(u, am.AssignRep(BAD,  BAD,  DONATE, rnGen), BAD)
  AssertRepEqual(u, am.AssignRep(BAD,  BAD,  REFUSE, rnGen), GOOD)
}

func TestCopy(u *testing.T) {
  // stern judging assessment module
  am1 := NewAssessModule(GOOD, BAD, BAD, GOOD, GOOD, BAD, BAD, GOOD, PASSERR)
  am2 := am1.Copy()

  AssertFalse(u, am1 == am2)
  AssertRepEqual(u, am1.bits[0], am2.bits[0])
  AssertRepEqual(u, am1.bits[1], am2.bits[1])
  AssertRepEqual(u, am1.bits[2], am2.bits[2])
  AssertRepEqual(u, am1.bits[3], am2.bits[3])
  AssertRepEqual(u, am1.bits[4], am2.bits[4])
  AssertRepEqual(u, am1.bits[5], am2.bits[5])
  AssertRepEqual(u, am1.bits[6], am2.bits[6])
  AssertRepEqual(u, am1.bits[7], am2.bits[7])

  am2.bits[0] = BAD
  am2.bits[1] = GOOD
  am2.bits[2] = GOOD
  am2.bits[3] = BAD
  am2.bits[4] = BAD
  am2.bits[5] = GOOD
  am2.bits[6] = GOOD
  am2.bits[7] = BAD

  AssertRepNotEqual(u, am1.bits[0], am2.bits[0])
  AssertRepNotEqual(u, am1.bits[1], am2.bits[1])
  AssertRepNotEqual(u, am1.bits[2], am2.bits[2])
  AssertRepNotEqual(u, am1.bits[3], am2.bits[3])
  AssertRepNotEqual(u, am1.bits[4], am2.bits[4])
  AssertRepNotEqual(u, am1.bits[5], am2.bits[5])
  AssertRepNotEqual(u, am1.bits[6], am2.bits[6])
  AssertRepNotEqual(u, am1.bits[7], am2.bits[7])
}

func TestGetBit(u *testing.T) {
  // stern judging assessment module
  am := NewAssessModule(GOOD, BAD, BAD, GOOD, GOOD, BAD, BAD, GOOD, PASSERR)

  AssertIntEqual(u, am.GetBit(0), 1)
  AssertIntEqual(u, am.GetBit(1), 0)
  AssertIntEqual(u, am.GetBit(2), 0)
  AssertIntEqual(u, am.GetBit(3), 1)
  AssertIntEqual(u, am.GetBit(4), 1)
  AssertIntEqual(u, am.GetBit(5), 0)
  AssertIntEqual(u, am.GetBit(6), 0)
  AssertIntEqual(u, am.GetBit(7), 1)  
}
