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
  flag.Parse()

  // set up the output file
  ofile, err := os.Create(*fname)
  if (err != nil) { panic (err) }
  defer ofile.Close()
  writer := bufio.NewWriter(ofile)

  // run simulation
  start := time.Now()
  var s *sim.SimEngine = sim.NewSimEngine(*numTribes,*numAgents)
  for g := 0; g < *gens; g++ {
    var p = s.PlayRounds(int32(*cost),int32(*benefit))
    // fmt.Printf("total payout for generation %5d: %7d\n", g, p)
    s.EvolveTribes()
    s.Reset()
    WriteStats(writer, g, s.GetStats(), p)
  }
  end := time.Now()

  writer.Flush()
  PrintStats(s.GetStats())

  fmt.Println("completed in ", end.Sub(start))
}

func PrintStats(s [8]int) {
  fmt.Printf("%d %d %d %d %d %d %d %d\n",
             s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7])
}

func WriteStats(w io.Writer, gen int, s [8]int, p int32) {
  fmt.Fprintf(w, "%d, %d, %d, %d, %d, %d, %d, %d, %d, %d\n",
                 gen, s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7], p)
}
