package main

import "flag"
import "simgpgg"

// default parameter values
const (
 GENS = 10      // default number of generations per simulation
 GENS_F = "g"   // flag for GENS parameter
 AGENTS = 10    // default number of agents per tribe
 AGENTS_F = "a" // flag for AGENTS parameter
 Z = 4          // default average degree of the graph (z)
 Z_F = "z"      // flag for Z parameter
 MULT = 3       // default contribution multiplier (r)
 MULT_F = "r"   // flag for MULT parameter
 COST = 1       // default contribution made by cooperators
 COST_F = "c"   // flag for COST parameter
 BETAE = 10     // selection strength for strategy updates
 BETAE_F = "betae"   // flag for BETAE parameter
 BETAA = 10     // selection strength for structure updates
 BETAA_F = "betaa"   // flag for BETAA parameter
 W = 0          // ratio of time scales for strategy and structure updates
 W_F = "w"      // flag for W parameter
)

/*
Run the simulation with the specified arguments.

Arguments:
  gens    - number of generations

Author: John Maloney
*/
func main() {
  // parse command line arguments
  numGens   := flag.Int(GENS_F, GENS, "number of generations to simulate")
  numAgents := flag.Int(AGENTS_F, AGENTS, "number of agents")
  avgdeg    := flag.Int(Z_F, Z, "average degree of the graph (z)")
  mult      := flag.Int(MULT_F, MULT, "contribution multiplier (r)")
  cost      := flag.Int(COST_F, COST, "cost to contribute")
  betae     := flag.Float64(BETAE_F, BETAE, "selection strength for strategy updates")
  betaa     := flag.Float64(BETAA_F, BETAA, "selection strength for structure updates")
  w         := flag.Float64(W_F, W, "ratio of time scales for strategy and structure updates")
  flag.Parse()

  // create the sim engine
  simeng := simgpgg.NewSimEngine(int32(*numAgents), int32(*numGens), int32(*avgdeg),
                                 int32(*mult), int32(*cost), *w, *betae, *betaa)

  // run the simulation
  simeng.RunSim()
}
