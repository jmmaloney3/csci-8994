package main

import "agent"
import "fmt"

func main() {
  fmt.Println("create simulation")
  var sim agent.SimEngine = agent.MakeSimEngine(2,2)
  var p = sim.PlayRounds(1,3)
  fmt.Println(" payout: ", p)
}
