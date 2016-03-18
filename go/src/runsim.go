package main

import "sim"
import "fmt"
import "flag"
import "time"

/*
Run the simulation with the specified arguments.

Arguments:
  tribes  - number of tribes
  agents  - number of agents per tribe
  cost    - cost c to donate
  benefit - benefit b received from donation

Author: John Maloney
*/
func main() {
  // parse command line arguments
  numTribes := flag.Int("t", 64, "number f tribes")
  numAgents := flag.Int("a", 64, "number of agents")
  cost := flag.Int("c", 1, "cost c to donate")
  benefit := flag.Int("b", 3, "benefit b received from donation")
  gens := flag.Int("g", 10, "number of generations to simulate")
  flag.Parse()

  // run simulation
  start := time.Now()
  var s *sim.SimEngine = sim.NewSimEngine(*numTribes,*numAgents)
  for g := 0; g < *gens; g++ {
    var p = s.PlayRounds(int32(*cost),int32(*benefit))
    fmt.Println("total payout for generation", g, ": ", p)
    s.CreateNextGen()
    s.Reset()
  }
  end := time.Now()

  fmt.Println("completed in ", end.Sub(start))
}
