package simpgg

import "math/rand"
import "fmt"
import "simbase"

/*
  Action module for a PGG agent.  The agent can choose to take one or more  of
  the following actions during a round:
    - Contribute to the common pool
    - Punish non-contributors
  If the agent choose to take neither of these actions then they are considered
  a non-participant.

  Besides these two action, the game can also be expanded to include the following
  additional actions:
    - Reward contributors
    - Punish non-punishers
  These actions may be included in a future version of this simulation.

  An action module for PGG is represented as two sets of 4 bits.  The first set
  of four bits determines whether the agent contributes while the second set of
  bits determines whether the agent punishes.
*/
type ActionModule struct {
  cam *simbase.ActionModule // determines when the agent contributes
  pam *simbase.ActionModule // determines when the agent punishes
  pexeerr float32
}

func NewActionModule(c1 bool, c2 bool, c3 bool, c4 bool,
                     p1 bool, p2 bool, p3 bool, p4 bool, pexeerr float32) *ActionModule {
  return &ActionModule { cam: simbase.NewActionModule(c1, c2, c3, c4, pexeerr),
                         pam: simbase.NewActionModule(p1, p2, p3, p4, pexeerr),
                         pexeerr: pexeerr }
}

func (am *ActionModule) Copy() *ActionModule {
  return &ActionModule { cam: am.cam.Copy(), pam: am.pam.Copy(), pexeerr: am.pexeerr }
}

func (self *ActionModule) ChooseContribute(agent simbase.Rep, group simbase.Rep, rnGen *rand.Rand) bool {
  return self.cam.ChooseAction(agent, group, rnGen)
}

func (self *ActionModule) ChoosePunish(agent simbase.Rep, group simbase.Rep, rnGen *rand.Rand) bool {
  return self.pam.ChooseAction(agent, group, rnGen)
}

// return true of the two modules have the same bits
func (self *ActionModule) SameBits(am *ActionModule) bool {
  return (self.cam.SameBits(am.cam) && self.pam.SameBits(am.pam) )
}

func (self *ActionModule) GetBit(i int) int {
  if (i < 4) {
    return self.cam.GetBit(i)
  } else {
    return self.pam.GetBit(i-4)
  }
}

func (self *ActionModule) WriteSimParams() {
  fmt.Printf("  \"pexeerr\":%.5f\n", self.pexeerr)
}
