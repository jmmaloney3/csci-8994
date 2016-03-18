package sim

type ActionModule struct {
  bits [4]bool
}

func NewActionModule() *ActionModule {
  // return &ActionModule { bits: [4]bool{true, false, true, false} }
  return &ActionModule { bits: [4]bool{RandBool(), RandBool(), RandBool(), RandBool()} }
}

func (self *ActionModule) ChooseDonate(donor *Agent, recipient *Agent) bool {
  if (donor.rep == GOOD) {
    if (recipient.rep == GOOD) {
      return self.bits[0]
    } else {
      return self.bits[1]
    }
  } else {
    if (recipient.rep == GOOD) {
      return self.bits[2]
    } else {
      return self.bits[3]
    }
  }
}
