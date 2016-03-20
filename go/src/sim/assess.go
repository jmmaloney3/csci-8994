package sim

import "math"

type AssessModule struct {
  bits [8]Rep
}

func NewAssessModule(r1 Rep, r2 Rep, r3 Rep, r4 Rep, r5 Rep,
                     r6 Rep, r7 Rep, r8 Rep) *AssessModule {
  return &AssessModule { bits: [8]Rep{r1, r2, r3, r4, r5, r6, r7, r8} }
}

func CopyAssessModule(am AssessModule) *AssessModule {
  return NewAssessModule(am.bits[0], am.bits[1], am.bits[2], am.bits[3],
                         am.bits[4], am.bits[5], am.bits[6], am.bits[7])
}

// return true of the two modules have the same bits
func (self *AssessModule) SameBits(am *AssessModule) bool {
  if (self == am) { return true }
  for i := 0; i < 8; i++ {
    if (self.bits[i] != am.bits[i]) {
      return false
    }
  }
  return true
}

func (self *AssessModule) GetBits() int {
  rval := int(0)
  for i := 0; i < 8; i++ {
    rval += self.GetBit(i) * int(math.Pow(2,float64(8-i)))
  }
  return rval
}

func (self *AssessModule) GetBit(i int) int {
  if (self.bits[i] == GOOD) {
    return 1
  } else {
    return 0
  }
}

func (self *AssessModule) AssignRep(donor Rep, recip Rep, act Act) Rep {
  if (donor == GOOD) {
    if (recip == GOOD) {
      if (act == DONATE) {
        return self.bits[0]
      } else { // action is REFUSE
        return self.bits[1]
      }
    } else { // recipient rep is BAD
      if (act == DONATE) {
        return self.bits[2]
      } else { // action is REFUSE
        return self.bits[3]
      }
    }
  } else { // donor rep is BAD
    if (recip == GOOD) {
      if (act == DONATE) {
        return self.bits[4]
      } else { // action is REFUSE
        return self.bits[5]
      }
    } else { // recipient rep is BAD
      if (act == DONATE) {
        return self.bits[6]
      } else { // action is REFUSE
        return self.bits[7]
      }
    }
  }
}