package main

import "flag"
import "time"
import "os"
import "path"
import "bufio"
import "fmt"
import "simgpgg"

// default parameter values
const (
 SIMS = 10      // default number of simulations to conduct
 SIMS_F = "s"   // flag for SIMS parameter
 GENS = 10      // default number of generations per simulation
 GENS_F = "g"   // flag for GENS parameter
 AGENTS = 10    // default number of agents per tribe
 AGENTS_F = "a" // flag for AGENTS parameter
 Z = 4          // default average degree of the graph (z)
 Z_F = "z"      // flag for Z parameter
 GTYPE = 0      // default graph type (0 = regular ring)
 GTYPE_F = "gtype" // flag for GTYPE parameter
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
 DNAME = "gpggdata"
 DNAME_F = "d"
 OWDIR = false
 OWDIR_F = "f"
)

/*
Run the simulation with the specified arguments.

Arguments:
  gens    - number of generations

Author: John Maloney
*/
func main() {
  // parse command line arguments
  numSims   := flag.Int(SIMS_F, SIMS, "number of simulations to conduct")
  numGens   := flag.Int(GENS_F, GENS, "number of generations to simulate")
  numAgents := flag.Int(AGENTS_F, AGENTS, "number of agents")
  avgdeg    := flag.Int(Z_F, Z, "average degree of the graph (z)")
  gtype     := flag.Int(GTYPE_F, GTYPE, "type of graph to use")
  mult      := flag.Int(MULT_F, MULT, "contribution multiplier (r)")
  cost      := flag.Int(COST_F, COST, "cost to contribute")
  betae     := flag.Float64(BETAE_F, BETAE, "selection strength for strategy updates")
  betaa     := flag.Float64(BETAA_F, BETAA, "selection strength for structure updates")
  w         := flag.Float64(W_F, W, "ratio of time scales for strategy and structure updates")
  dname     := flag.String(DNAME_F, DNAME, "directory to write stats")
  owDir     := flag.Bool(OWDIR_F, OWDIR, "overwrite data if directory exists")
  flag.Parse()

  // set up the output director for the experiment
  var err error
  err = os.MkdirAll(*dname, os.ModePerm)
  if (err != nil) {
    if (os.IsExist(err)) {
      if (!(*owDir)) {
        fmt.Fprintf(os.Stderr, "ERROR: directory exists: %v\n", *dname)
        // don't overwrite data - exit program
        return
      }
    } else {
      panic (err)
    }
  }

  // make stdout a valid JSON array
  fmt.Println("[")

  // perform the requested number of simulations
  for s := 0; s < *numSims; s++ {
    // set up the output directory and files for the simulation
    var err error
    // -- create directory
    simdname := path.Join(*dname, fmt.Sprintf("%s%d", "sim", s))
    err = os.Mkdir(simdname, os.ModePerm)
    if (err != nil) {
      if (os.IsExist(err)) {
        if (!(*owDir)) {
          fmt.Fprintf(os.Stderr, "ERROR: directory exists: %v\n", simdname)
          // don't overwrite data - exit program
          return
        }
      } else {
        panic (err)
      }
    }

    // -- file for population statistics (strategy percentages)
    var psfile *os.File
    psfname := path.Join(simdname, "pstat.csv")
    psfile, err = os.Create(psfname)
    if (err != nil) { panic (err) }
    defer psfile.Close()
    psWriter := bufio.NewWriter(psfile)
    // -- file for degree histogram
    var dhfile *os.File
    dhfname := path.Join(simdname, "dhist.csv")
    dhfile, err = os.Create(dhfname)
    if (err != nil) { panic (err) }
    defer dhfile.Close()
    dhWriter := bufio.NewWriter(dhfile)

    start := time.Now()

    // create the sim engine
    simeng := simgpgg.NewSimEngine(int32(*numAgents), int32(*numGens), int32(*gtype), int32(*avgdeg),
                                   int32(*mult), int32(*cost), *w, *betae, *betaa)

    // output simulation parameters to stdout
    fmt.Println("{")
    fmt.Printf("  \"params\":\n")
    fmt.Printf("%v", simeng)
    fmt.Print(",\n")

    // run the simulation
    gens := simeng.RunSim(psWriter, dhWriter)

    end := time.Now()

    psWriter.Flush()
    dhWriter.Flush()

    // write sim results
    fmt.Printf("  \"results\":\n")
    fmt.Printf("  {\n")
    fmt.Printf("  \"psfile\":\"%s\",\n", psfname)
    fmt.Printf("  \"dhfile\":\"%s\",\n", psfname)
    fmt.Printf("  \"ngens-completed\":%d,\n", gens)
    fmt.Printf("  \"runtime\":\"%v\"\n",end.Sub(start))
    fmt.Printf("  }\n")
    if (s+1 >= *numSims) {
      fmt.Println("}")
    } else {
      fmt.Println("},")
    }
  }

  // make stdout a valid JSON array
  fmt.Println("]")
}
