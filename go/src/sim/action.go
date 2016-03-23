package sim

type ActionModule struct {
  bits [4]bool
  errP float32 // mu_e - execution error - fail to donate
}

func NewActionModule(b1 bool, b2 bool, b3 bool, b4 bool) *ActionModule {
  return &ActionModule { bits: [4]bool{b1, b2, b3, b4},
                         errP: 0.001 }
}

func (self *ActionModule) ChooseDonate(donor Rep, recip Rep) bool {
  var rval bool
  if (donor == GOOD) {
    if (recip == GOOD) {
      rval = self.bits[0]
    } else {
      rval = self.bits[1]
    }
  } else {
    if (recip == GOOD) {
      rval = self.bits[2]
    } else {
      rval = self.bits[3]
    }
  }
  // check for execution error
  if (RandPercent() < float64(self.errP)) {
    rval = !rval
  }
  // return the action
  return rval
}

func (self *ActionModule) GetBit(i int) int {
  if (self.bits[i]) {
    return 1
  } else {
    return 0
  }
}
