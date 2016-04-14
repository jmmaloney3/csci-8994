package sim

import "fmt"
import "math"
import "math/rand"

func AssignRepManualTest() {
  fmt.Println("AssessModule.AssignRep")
  errorRate := float32(0.5)

  // stern judging assessment module
  am := NewAssessModule(GOOD, BAD, BAD, GOOD, GOOD, BAD, BAD, GOOD, errorRate)

  rnGen := NewRandNumGen()

  // print out some assign rep results
  errors := 0
  N := 100
  for i := 0; i < N; i++ {
    rep := am.AssignRep(GOOD, GOOD, DONATE, rnGen)
    if (rep == BAD) { errors++ }
    //fmt.Printf("rep: %v\n", rep)
  }
  fmt.Printf("  expected error rate: %6.4f\n", errorRate)
  fmt.Printf("  actual error rate:   %6.4f\n", float64(errors)/float64(N))
}

func ChooseDonateManualTest() {
  fmt.Println("ActionModule.ChooseDonate")
  errorRate := float32(0.5)

  // CO action module
  am := NewActionModule(true, false, true, false, errorRate)

  rnGen := NewRandNumGen()

  // print out some choose donate results
  errors := 0
  N := 100
  for i := 0; i < N; i++ {
    choice := am.ChooseDonate(GOOD, GOOD, rnGen)
    if (!choice) { errors++ }
    //fmt.Printf("choice: %t\n", choice)
  }
  fmt.Printf("  expected error rate: %6.4f\n", errorRate)
  fmt.Printf("  actual error rate:   %6.4f\n", float64(errors)/float64(N))
}

func CloneWithMutationsManualTest() {
  fmt.Println("ActionModule.CloneWithMutations")
  pactmut := float64(0.5)

  // CO action module
  am := NewActionModule(true, false, true, false, PEXEERR)

  rnGen := NewRandNumGen()

  mutations := 0
  N := 100
  for i := 0; i < N; i++ {
    clone := am.CloneWithMutations(pactmut, rnGen)
    // count mutations
    for j := 0; j < 4; j++ {
      if (am.bits[j] != clone.bits[j]) { mutations++ }
    }
  }
  fmt.Printf("  expected mutation rate: %6.4f\n", pactmut)
  fmt.Printf("  actual error rate:   %6.4f\n", float64(mutations)/float64(N*4))
}

func AssignRolesManualTest() {
  fmt.Println("Tribe.AssignRoles")
  rnGen := NewRandNumGen()
  const numAgents = 2

  // tribe with numAgents agents
  t := NewTribe(numAgents, PASSERR, PACTMUT, PEXEERR, rnGen)

  // print out some randomly selected parent agents
  var selected [numAgents]int
  N := 100
  for i := 0; i < N; i++ {
    a,_ := t.AssignRoles(t.agents[0], t.agents[1], rnGen)
    for j := 0; j < numAgents; j++ {
      if (a == t.agents[j]) {
        selected[j]++
      }
    }
    // fmt.Printf("parent payout: %d\n", a.payout)
  }

  for i := 0; i < numAgents; i++ {
    fmt.Printf("  expected rate for agent %d: %6.4f\n", i, 0.5)
    fmt.Printf("  actual rate for agent   %d: %6.4f\n", i, float64(selected[i])/float64(N))
  }
}

func SelectParentManualTest() {
  fmt.Println("Tribe.SelectParent")
  rnGen := NewRandNumGen()
  const numAgents = 3

  // tribe with numAgents agents
  t := NewTribe(numAgents, PASSERR, PACTMUT, PEXEERR, rnGen)

  // set agent's payouts
  t.agents[0].payout = 10
  t.agents[1].payout = 20
  t.agents[2].payout = 30

  // update the tribe total payouts
  for i := 0; i < numAgents; i++ {
    t.totalPayouts += t.agents[i].payout
  }

  // print out some randomly selected parent agents
  var selected [numAgents]int
  N := 100
  for i := 0; i < N; i++ {
    a := t.SelectParent(rnGen)
    for j := 0; j < numAgents; j++ {
      if (a == t.agents[j]) {
        selected[j]++
      }
    }
    // fmt.Printf("parent payout: %d\n", a.payout)
  }

  for i := 0; i < numAgents; i++ {
    fmt.Printf("  expected rate for agent %d: %6.4f\n", i, float64(t.agents[i].payout)/float64(t.totalPayouts))
    fmt.Printf("  actual rate for agent   %d: %6.4f\n", i, float64(selected[i])/float64(N))
  }
}

func ConflictManualTest() {
  fmt.Println("SimEngine.Conflict")
  rnGen := NewRandNumGen()
  const numTribes = 2
  const numAgents = 2
  const beta = .6

  // simengine with numTribes tribes and numAgents agents per tribes
  s := NewDefaultSimEngine(numTribes, numAgents, true, true)
  s.beta = beta

  // set tribe payouts
  t1_idx := 0
  t1 := s.tribes[t1_idx]
  t1.totalPayouts = 10
  t2_idx := 1
  t2 := s.tribes[t2_idx]
  t2.totalPayouts = 20

  // calculate conflict threshold
  diff := t2.AvgPayout() - t1.AvgPayout()
  p  := math.Pow(float64(1) + math.Exp(diff*(-beta)), float64(-1))

  t2Wins := 0
  N := 100
  for i := 0; i < N; i++ {
    w_idx, _ := s.Conflict(t1_idx, t2_idx, rnGen)
    if (w_idx == t2_idx) { t2Wins++ }
  }

  fmt.Printf("  expected win rate: %6.4f\n", p)
  fmt.Printf("  actual win rate:   %6.4f\n", float64(t2Wins)/float64(N))
}

func SingleTribeSim() {
  numAgents := 5
  passerr := float32(0)
  pactmut := float64(0) // float32(PACTMUT)
  pexeerr := float32(0)
  rnGen := NewRandNumGen()
  numGens := 20
  cost := int32(1)
  benefit := int32(3)

  // create the tribe
  t := NewTribe(numAgents, passerr, pactmut, pexeerr, rnGen)

  // set the tribe assessment modue to stern-judging
  sj := NewAssessModule(GOOD, BAD, BAD, GOOD, GOOD, BAD, BAD, GOOD, passerr)
  t.assessMod = sj
  fmt.Printf("assess module: [%d][%08b]\n", sj.GetBits(), sj.GetBits())

  // print put CO action module for reference
  co := NewActionModule(true, false, true, false, pexeerr)
  fmt.Println(co.bits)
  fmt.Printf("CO: [%d][%04b]\n", co.GetBits(), co.GetBits())

  // configure one agent to use CO
  t.agents[0].actMod = co

  // configure all agents to use CO
  /*
  for i := 0; i < numAgents; i++ {
    t.agents[i].actMod = co
  }
  */

  // run generations
  for g := 0; g < numGens; g++ {
    t.PlayRounds(cost, benefit, rnGen)
    // print current state
    var a *Agent
    for i := 0; i < numAgents; i++ {
      a = t.agents[i]
      fmt.Printf("%4d %1d [%04b] ", a.payout, a.rep, a.actMod.GetBits())
    }
    fmt.Println()
    // get next generation
    t = t.CreateNextGen(rnGen)
  }
}

func MultiTribeSim() {
  numTribes := 64
  numAgents := 64
  passerr :=  float32(0) // PASSERR
  pactmut := PACTMUT // float32(0)
  pexeerr :=  float32(0) // PEXEERR
  passmut := PASSMUT // +0.005 // float32(0)
  pcon    := float64(0.5) // PCON // float64(0.1)
  beta    := BETA // math.Inf(int(1))
  eta     := float64(1) // ETA
  pmig    := PMIG // float64(1)
  numGens := 10000
  cost := int32(1)
  benefit := int32(3)

  // create parameter map for floats
  var params = make(map[string]float64)

  // populate arg map for floats
  params[PASSE_F] = float64(passerr)
  params[PACTM_F] = float64(pactmut)
  params[PEXEE_F] = float64(pexeerr)
  params[PCON_F]  = pcon
  params[BETA_F]  = beta
  params[ETA_F]   = eta
  params[PMIG_F]  = pmig
  params[PASSM_F] = float64(passmut)

  singledef := SINGLE_DEF
  passmutall := PASSMUT_ALL
  useMP := true
  useAM := true

  // create parameter map for booleans
  var bparams = make(map[string]bool)

  // populate arg map for booleans
  bparams[SINGLE_DEF_F]  = singledef
  bparams[PASSMUT_ALL_F] = passmutall
  bparams[USEAM_F]       = useAM
  bparams[NOMP_F]        = !useMP

  // create the simengine
  s := NewSimEngine(numTribes, numAgents, params, bparams)

  // calculate max and min payouts
  minPO, maxPO := CalcMinMaxTribalPayouts(numAgents, cost, benefit)
  minMaxDiff := maxPO - minPO

  // convert into overall simulation payouts
  simMinPO := minPO * int32(numTribes)
  simMaxPO := maxPO * int32(numTribes)
  simMinMaxDiff := simMaxPO - simMinPO

  // configure all agents to use the CO action module
  /*
  co := NewActionModule(true, false, true, false, pactmut, pexeerr)
  fmt.Println(co.bits)
  fmt.Printf("CO: [%d][%04b]\n", co.GetBits(), co.GetBits())
  for i := 0; i < numTribes; i++ {
    for j := 0; j < numAgents; j++ {
      s.tribes[i].agents[j].actMod = co
    }
  }
 */
  // run generations
  var assmodCounts = make(map[int]int)
  var assmodNum int
  var actmodCounts = make(map[int]int)
  var actmodNum int
  var newTribes []*Tribe
  // var pMax float32
  var errRate float64
  var error float64

  var evolveCount int
  totalAgents := numTribes * numAgents
  for g := 0; g < numGens; g++ {
    // play rounds
    newTribes = s.PlayRounds(cost, benefit)
    // count up the number of tribes following each assessment module
    var t *Tribe
    var a *Agent
    for i := 0; i < numTribes; i++ {
      t = s.tribes[i]
      // get assessment module number
      assmodNum = t.assessMod.GetBits()
      assmodCounts[assmodNum] = assmodCounts[assmodNum] + 1
      for j := 0; j < numAgents; j++ {
        a = t.agents[j]
        actmodNum = a.actMod.GetBits()
        actmodCounts[actmodNum] = actmodCounts[actmodNum] + 1
      }
    }
    // print the stats for this round
    var assmodCount int
    var actmodCount int
    var ok bool
    var p float32
    // print the payout stats
    //pMax = float32(s.totalPayouts - min)/float32(max - min)
    //fmt.Printf("%5d: total payout: %6d  %% of max: %6.4f  passmut: %6.4f\n", g, s.totalPayouts, pMax, s.passmut)
    // calculate the error rate
    error = float64(simMinMaxDiff) - float64(s.totalPayouts - simMinPO)
    errRate = error / float64(simMinMaxDiff)
    fmt.Printf("%5d: total payout: %6d  err rate: %6.4f  passmut: %6.4f\n", g, s.totalPayouts, errRate, s.passmut)

    // print out the total payout and error rate for this simulation
    //fmt.Printf("%5d: total payout: %6d  err rate: %6.4f\n", g, s.totalPayouts, errRate)

    // print out the error rate and mutation rate for each tribe
    var tErrRate float64
    var tMutRate float64
    for i := 0; i < s.numTribes; i++ {
      t = s.tribes[i]
      tErrRate = (float64(minMaxDiff) - float64(t.totalPayouts - minPO))/float64(minMaxDiff)
      tMutRate = CalcAdaptTribalMutRate(float64(t.totalPayouts), minPO, maxPO)
      fmt.Printf("[%02d:%3.2f(%5.4f)]", i, tErrRate, tMutRate)
      // adapt mutation rate for each agent - if using AM
      if (useAM) {
        for j := 0; j < t.numAgents; j++ {
          t.agents[j].pactmut = tMutRate
        }
      }
    }
    fmt.Println()

    // print assessment modules used
    for i := 0; i < 256; i++ {
      assmodCount, ok = assmodCounts[i]
      if (ok) {
        fmt.Printf(" [%08b]: %03d", i, assmodCount)
      }
    }
    fmt.Println()
    // print action modules used
    for i := 0; i < 16; i++ {
      actmodCount, ok = actmodCounts[i]
      if (ok) {
        p = float32(actmodCount)/float32(totalAgents)
        fmt.Printf(" [%04b]: %05.3f", i, p)
      }
    }
    fmt.Println()
    fmt.Println()

    // set assessment module mutation rate based on most recent results
    // s.passmut = PASSMUT + float32((math.Pow(float64(1) - float64(pMax), float64(4)))*float64(0.002))
    s.passmut = PASSMUT + math.Pow(errRate, float64(4))*float64(0.002)
    //s.passmut = CalcAdaptTribalMutRate(float64(s.totalPayouts), simMinPO, simMaxPO)

    // evolve tribes to next generation
    evolveCount++
    if (evolveCount >= 50) {
      //s.EvolveTribes(newTribes, minPO, maxPO)
      s.EvolveTribes2(newTribes, useAM, minPO, maxPO)
      evolveCount = 0
    }
    s.Reset()
    // clear maps for next round
    assmodCounts = make(map[int]int)
    actmodCounts = make(map[int]int)
  }
  fmt.Printf("min: %d  max: %d\n", simMinPO, simMaxPO)

}

// Randomly select an tribe.  The chance that a tribe is selected is
// proportional to its fitness.
func (self *SimEngine) SelectParentTribe(rnGen *rand.Rand) *Tribe {
  ri := int32(RandInt(rnGen, int64(self.totalPayouts)))
  thresh := int32(0);
  var parent *Tribe
  for i := 0; i < self.numTribes; i++ {
    thresh += self.tribes[i].totalPayouts
    if (ri <= thresh) {
      parent = self.tribes[i]
      break
    }
  }
  return parent
}

// Create the next generation by propagating action modules to the next
// generation based on the fitness those modules achieved.
func (self *SimEngine) EvolveTribes2(nextGen []*Tribe, useAM bool, minPO int32, maxPO int32) {
  // propagate successful assessment modules among the tribes
  for i := 0; i < self.numTribes; i++ {
    // select parent tribe from current generation
    parent := self.SelectParentTribe(self.rnGen)
    // switch assessment module of tribe i in next generation
    nextGen[i].assessMod = parent.assessMod.Copy()

    // migrate agents
    self.MigrateAgents(parent, nextGen[i], self.rnGen)

    // mutate assessment module
    // get assessment module bit mutation rate
    var mutRate float64
    if (useAM) {
      mutRate = CalcAdaptTribalMutRate(float64(parent.totalPayouts), minPO, maxPO)
    } else {
      mutRate = self.passmut
    }
    for i := 0; i < 8; i++ {
      if (RandPercent(self.rnGen) < float64(mutRate)) {
        if (nextGen[i].assessMod.bits[i] == GOOD) {
          nextGen[i].assessMod.bits[i] = BAD
        } else {
          nextGen[i].assessMod.bits[i] = GOOD
        }
      }
    }
  }

  // replace the original tribes with the new tribes
  self.tribes = nextGen
}
