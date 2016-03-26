package main

import "sim"
import "fmt"
import "flag"
import "time"
import "os"
import "bufio"
import "io"

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
  fname := flag.String("f", "stats.csv", "file to collect stats")
  useMP := flag.Bool("mp", false, "whether to use multiprocessing")
  beta := flag.Float64("beta", 1.2, "selection strength")
  flag.Parse()

  // set up the output file
  ofile, err := os.Create(*fname)
  if (err != nil) { panic (err) }
  defer ofile.Close()
  writer := bufio.NewWriter(ofile)
  WriteHeader(writer)

  // run simulation
  start := time.Now()
  var s *sim.SimEngine = sim.NewSimEngine(*numTribes,*numAgents, *useMP)
  s.Beta = *beta

  // output simulation parameters
  WriteSimParams(s, gens, cost, benefit, fname)

  // calculate max possible payout per generation
  max, min := s.MaxMinPayouts(int32(*cost),int32(*benefit))
  // execute simulation
  for g := 0; g < *gens; g++ {
    var p = s.PlayRounds(int32(*cost),int32(*benefit))
    // fmt.Printf("total payout for generation %5d: %7d\n", g, p)
    s.EvolveTribes()
    s.Reset()
    n, a := s.GetStats()
    WriteStats(writer, g, *numTribes, *numAgents, n, a, p, min, max)
  }
  end := time.Now()

  writer.Flush()

  fmt.Println("completed in ", end.Sub(start))
}

func WriteHeader(w io.Writer) {
  fmt.Fprintf(w, "gen,t,a,n0,n1,n2,n3,n4,n5,n6,n7,a0,a1,a2,a3,po,minpo,maxpo\n")
}
func WriteStats(w io.Writer, gen int, numTribes int, numAgents int,
                n [8]int, a [4]int, p int32, min int32, max int32) {
  fmt.Fprintf(w, "%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d\n",
                 gen, numTribes, numAgents,
                 n[0], n[1], n[2], n[3], n[4], n[5], n[6], n[7],
                 a[0], a[1], a[2], a[3],
                 p, min, max)
}

func WriteSimParams(s *sim.SimEngine, gens *int, cost *int, benefit *int, fname *string) {
  // output simulation parameters
  fmt.Println("IR simulation parameters:")
  fmt.Printf("  num gens:     %8d\n", *gens)
  fmt.Printf("  cost:         %8d\n", *cost)
  fmt.Printf("  benefit:      %8d\n", *benefit)
  fmt.Printf("  out file:     %s\n", *fname)
  s.WriteSimParams()
}
