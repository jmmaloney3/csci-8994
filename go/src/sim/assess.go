package sim

type AssessModule struct {
  bits [8]Rep
}

func NewAssessModule(r1 Rep, r2 Rep, r3 Rep, r4 Rep, r5 Rep,
                     r6 Rep, r7 Rep, r8 Rep) *AssessModule {
  return &AssessModule { bits: [8]Rep{r1, r2, r3, r4, r5, r6, r7, r8} }
}

func CopyAssessModule(am AssessModule) *AssessModule {
    copiedBits := [8]Rep{am.bits[0], am.bits[1], am.bits[2], am.bits[3],
                        am.bits[4], am.bits[5], am.bits[6], am.bits[7]}
    return &AssessModule { bits: copiedBits }
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
