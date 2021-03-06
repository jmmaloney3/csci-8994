package simbase

import "math/rand"
import "fmt"

type ActionModule struct {
  bits [4]bool
  pexeerr float32 // mu_e - execution error - fail to take action
}

func NewActionModule(b1 bool, b2 bool, b3 bool, b4 bool, pexeerr float32) *ActionModule {
  return &ActionModule { bits: [4]bool{b1, b2, b3, b4},
                         pexeerr: pexeerr }
}

func (am *ActionModule) Copy() *ActionModule {
  return NewActionModule(am.bits[0], am.bits[1], am.bits[2], am.bits[3], am.pexeerr)
}

func (self *ActionModule) ChooseAction(actionTaker Rep, other Rep, rnGen *rand.Rand) bool {
  var rval bool
  if (actionTaker == GOOD) {
    if (other == GOOD) {
      rval = self.bits[0]
    } else {
      rval = self.bits[1]
    }
  } else {
    if (other == GOOD) {
      rval = self.bits[2]
    } else {
      rval = self.bits[3]
    }
  }
  // check for execution error
  if (RandPercent(rnGen) < float64(self.pexeerr)) {
    rval = !rval
  }
  // return the action
  return rval
}

// return true of the two modules have the same bits
func (self *ActionModule) SameBits(am *ActionModule) bool {
  if (self == am) { return true }
  for i := 0; i < 4; i++ {
    if (self.bits[i] != am.bits[i]) {
      return false
    }
  }
  return true
}

func (self *ActionModule) GetBit(i int) int {
  if (self.bits[i]) {
    return 1
  } else {
    return 0
  }
}

func (self *ActionModule) WriteSimParams() {
  fmt.Printf("  \"pexeerr\":%.5f\n", self.pexeerr)
}
