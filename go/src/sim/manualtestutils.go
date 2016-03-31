package sim

import "fmt"

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
