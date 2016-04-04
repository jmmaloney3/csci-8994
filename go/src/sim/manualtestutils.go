package sim

import "fmt"
import "math"

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
  am := NewActionModule(true, false, true, false, PACTMUT, errorRate)

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
  pactmut := float32(0.5)

  // CO action module
  am := NewActionModule(true, false, true, false, pactmut, PEXEERR)

  rnGen := NewRandNumGen()

  mutations := 0
  N := 100
  for i := 0; i < N; i++ {
    clone := am.CloneWithMutations(rnGen)
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
  s := NewDefaultSimEngine(numTribes, numAgents, true)
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
  pactmut := float32(0) // float32(PACTMUT)
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
  co := NewActionModule(true, false, true, false, pactmut, pexeerr)
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
  passerr := float32(0)
  pactmut := float32(0) // float32(PACTMUT)
  pexeerr := float32(0)
  passmut := PASSMUT+0.005 // float32(0)
  pcon    := PCON // float64(0.1)
  beta    := BETA // math.Inf(int(1))
  eta     := ETA // float64(0.5)
  pmig    := PMIG // float64(1) // since all agents are using the same AM, this doesn't matter
  numGens := 10000
  cost := int32(1)
  benefit := int32(3)

  // create parameter map
  var params = make(map[string]float64)

  // populate arg maps
  params[PASSE_F] = float64(passerr)
  params[PACTM_F] = float64(pactmut)
  params[PEXEE_F] = float64(pexeerr)
  params[PCON_F]  = pcon
  params[BETA_F]  = beta
  params[ETA_F]   = eta
  params[PMIG_F]  = pmig
  params[PASSM_F] = float64(passmut)

  // create the simengine
  s := NewSimEngine(numTribes, numAgents, params, true)

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
  var newTribes []*Tribe
  for g := 0; g < numGens; g++ {
    newTribes = s.PlayRounds(cost, benefit)
    // print current state
    var t *Tribe
    for i := 0; i < numTribes; i++ {
      t = s.tribes[i]
      fmt.Printf("%5d [%08b] ", t.totalPayouts, t.assessMod.GetBits())
    }
    fmt.Println()
    // evolve tribes to next generation
    s.EvolveTribes(newTribes)
  }
  max, min := s.MaxMinPayouts(cost, benefit)
  fmt.Printf("max: %d  min: %d\n", max, min)

}
