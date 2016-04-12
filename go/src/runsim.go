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
  gens    := flag.Int(sim.GENS_F, sim.NUMGENS, "number of generations to simulate")
  cost    := flag.Int(sim.COST_F, sim.COST, "cost c to donate")
  benefit := flag.Int(sim.BEN_F,  sim.BENEFIT, "benefit b received from donation")
  numTribes := flag.Int(sim.TRIBES_F, sim.NUMTRIBES, "number of tribes")
  numAgents := flag.Int(sim.AGENTS_F, sim.NUMAGENTS, "number of agents")
  beta    := flag.Float64(sim.BETA_F, sim.BETA, "conflict selection strength")
  eta     := flag.Float64(sim.ETA_F, sim.ETA, "bit switch selection strength")
  pcon    := flag.Float64(sim.PCON_F, sim.PCON, "conflict probability")
  pmig    := flag.Float64(sim.PMIG_F, sim.PMIG, "migration probability")
  passmut := flag.Float64(sim.PASSM_F, sim.PASSMUT, "assess module bit mutation probability")
  pactmut := flag.Float64(sim.PACTM_F, sim.PACTMUT, "action module bit mutation probability")
  passerr := flag.Float64(sim.PASSE_F, sim.PASSERR, "assessment error probability")
  pexeerr := flag.Float64(sim.PEXEE_F, sim.PEXEERR, "execution error probability")
  fname   := flag.String(sim.FNAME_F, sim.FNAME, "file to collect stats")
  noMP    := flag.Bool(sim.NOMP_F, sim.NOMP, "turn off multiprocessing")
  noAM    := flag.Bool(sim.NOAM_F, sim.NOAM, "turn off adaptive mutation")
  flag.Parse()

  // create parameter map
  var params = make(map[string]float64)
  params[sim.BETA_F]  = *beta
  params[sim.ETA_F]   = *eta
  params[sim.PCON_F]  = *pcon
  params[sim.PMIG_F]  = *pmig
  params[sim.PASSM_F] = *passmut
  params[sim.PACTM_F] = *pactmut
  params[sim.PASSE_F] = *passerr
  params[sim.PEXEE_F] = *pexeerr

  // set up the output file
  ofile, err := os.Create(*fname)
  if (err != nil) { panic (err) }
  defer ofile.Close()
  writer := bufio.NewWriter(ofile)
  WriteHeader(writer)

  start := time.Now()

  // create simulation
  var s *sim.SimEngine = sim.NewSimEngine(*numTribes, *numAgents, params, !*noAM, !*noMP)

  // output simulation parameters
  fmt.Println("[")
  WriteSimParams(s, *gens, *cost, *benefit, *fname)
  fmt.Println(",")

  // calculate max and min possible payouts per generation
  minPO, maxPO := s.MinMaxTribalPayouts(int32(*cost),int32(*benefit))
  simMinPO := minPO * int32(*numTribes)
  simMaxPO := maxPO * int32(*numTribes)

  // execute simulation
  var p int32
  var nextGen []*sim.Tribe

  for g := 0; g < *gens; g++ {
    nextGen = s.PlayRounds(int32(*cost),int32(*benefit))
    p = s.GetTotalPayouts()
    s.EvolveTribes(nextGen, minPO, maxPO)
    s.Reset()
    n, a, allc, alld := s.GetStats()
    WriteStats(writer, g, *numTribes, *numAgents, n, a, allc, alld, p, simMinPO, simMaxPO)
  }
  end := time.Now()

  writer.Flush()

  fmt.Println("{\n  \"runtime\":", end.Sub(start), "\n}")
  fmt.Println("]")
}

func WriteHeader(w io.Writer) {
  fmt.Fprintf(w, "gen,t,a,n0,n1,n2,n3,n4,n5,n6,n7,a0,a1,a2,a3,allc,alld,po,minpo,maxpo\n")
}
func WriteStats(w io.Writer, gen int, numTribes int, numAgents int,
                n [8]int, a [4]int, allcCnt int, alldCnt int,
                p int32, min int32, max int32) {
  fmt.Fprintf(w, "%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d\n",
                 gen, numTribes, numAgents,
                 n[0], n[1], n[2], n[3], n[4], n[5], n[6], n[7],
                 a[0], a[1], a[2], a[3],
                 allcCnt, alldCnt,
                 p, min, max)
}

func WriteSimParams(s *sim.SimEngine, gens int, cost int, benefit int, fname string) {
  // output simulation parameters
  fmt.Println("{")
  fmt.Printf("  \"simtype\":\"IR\",\n")
  fmt.Printf("  \"ngens\":%d,\n", gens)
  fmt.Printf("  \"cost\":%d,\n", cost)
  fmt.Printf("  \"benefit\":%d,\n", benefit)
  fmt.Printf("  \"ofile\":\"%s\",\n", fname)
  s.WriteSimParams()
  fmt.Println("}")
}
