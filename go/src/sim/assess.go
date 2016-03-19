package sim

type AssessModule struct {
  bits [8]bool
}

func NewAssessModule(b1 bool, b2 bool, b3 bool, b4 bool,
  b5 bool, b6 bool, b7 bool, b8 bool) *AssessModule {
  // return &ActionModule { bits: [4]bool{true, false, true, false} }
  return &AssessModule { bits: [8]bool{b1, b2, b3, b4, b5, b6, b7, b8} }
}

func (self *AssessModule) AssignRep(donor *Agent, recip *Agent, act Act) bool {
  if (donor.rep == GOOD) {
    if (recip.rep == GOOD) {
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
    if (recip.rep == GOOD) {
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
