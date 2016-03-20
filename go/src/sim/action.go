package sim

type ActionModule struct {
  bits [4]bool
}

func NewActionModule(b1 bool, b2 bool, b3 bool, b4 bool) *ActionModule {
  return &ActionModule { bits: [4]bool{b1, b2, b3, b4} }
}

func (self *ActionModule) ChooseDonate(donor Rep, recip Rep) bool {
  if (donor == GOOD) {
    if (recip == GOOD) {
      return self.bits[0]
    } else {
      return self.bits[1]
    }
  } else {
    if (recip == GOOD) {
      return self.bits[2]
    } else {
      return self.bits[3]
    }
  }
}
