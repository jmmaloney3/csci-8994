package simgpgg

import "math/rand"
import "goraph"
import "fmt"
import "io"

// A simulation engine for simulating the public goods games
// played among agents occupying the nodes of a graph.
type SimEngine struct {
  graph goraph.Graph // the graph that holds the agents
  numAgents int32    // total number of agents being simulated
  avgdeg int32       // average degree of the graph (z)
  agents []*Agent    // list of agents
  Nc int32           // current number of cooperators
  Nd int32           // current number of defectors
  numGens int32      // number of generations to be simulated
  mult int32         // contribution multiplier (r)
  cost int32         // contribution to game
  betae float64      // the selection strength for strategy updates
  betaa float64      // the selection strength for structure updates
  W float64          // relative frequency of structure updates
  rnGen *rand.Rand   // hold a RN generator
}

// Make a new SimEngine with the specified parameters
func NewSimEngine(numAgents int32, numGens int32, avgdeg int32, mult int32,
                  cost int32, W float64, betae float64, betaa float64) *SimEngine {
  // initialize simengine
  graph := NewRegularRing(numAgents, avgdeg)
  rnGen := NewRandNumGen()
  // create the agents
  agents := make([]*Agent, numAgents)
  // create agents
  Nd := int32(0)
  Nc := int32(0)
  for i := int32(0); i < numAgents; i++ {
    cooperate := RandBool(rnGen)
    agents[i] = NewAgent(cooperate)
    if (cooperate) {
      Nc += 1
    } else {
      Nd += 1
    }
  }

  return &SimEngine { numAgents: numAgents, numGens: numGens, avgdeg: avgdeg,
                      mult: mult, cost: cost, W: W, betae: betae, betaa: betaa,
                      rnGen: rnGen, graph: graph, agents: agents,
                      Nc: Nc, Nd: Nd }
}

func (self *SimEngine) RunSim(psWriter io.Writer, dhWriter io.Writer) int32 {
  // write header to population stats files
  self.WritePStatsHeader(psWriter)

  // calculate probability that a structure update occurs
  stratUpdProb := float64(1)/(float64(1) + self.W)

  // loop until one strategy is eliminated or the max num of gens is reached
  var g int32
  for g = int32(0); (!self.SimComplete(g)); g++ {
    // randomly select an agent
    x := goraph.Vertex(RandInt(self.rnGen, int64(self.numAgents)))
    // get the neighbors of x
    Nx := self.graph.Neighbors(x)
    // randomly select a neighbr of x
    y := goraph.Vertex(Nx[RandInt(self.rnGen, int64(len(Nx)))])
    // get the neighbors of y
    Ny := self.graph.Neighbors(y)
    // create combined list of agents without duplicates
    sponsors := make([]goraph.Vertex,2)
    sponsors[0] = x
    sponsors[1] = y
    // need accurate payout information for agents x and y
    // -- set their payouts equal to zro
    for i := 0; i < len(sponsors); i++ {
      self.agents[sponsors[i]].payouts = float64(0)
    }
    // add neighbors to list of game sponsors
    // -- don't need accurate payout information for these agents
    sponsors = append(sponsors, Nx...)
    sponsors = append(sponsors, Ny...)
    sponsors = removeDuplicates(sponsors)
    // play the games
    for i := 0; i < len(sponsors); i++ {
      sponsor := sponsors[i]
      players := append(self.graph.Neighbors(sponsor), sponsor)
      self.PlayGame(players)
    }

    if (RandProb(self.rnGen) <= stratUpdProb) {
      // update agent strategy - if appropriate
      self.UpdateStrategy(x, y)
    } else {
      // update network structure - if appropriate
      self.UpdateStructure(x, y)
    }
    // write out population stats
    self.WritePStats(psWriter, g)
  }

  self.DegreeHistogramData(dhWriter)

  // return number of generations completed
  return g
}

func (self *SimEngine) SimComplete(genNum int32) bool {
  return (genNum >= self.numGens) || (self.Nc >= self.numAgents) || (self.Nd >= self.numAgents)
}

// calculate the public goods payout for the specified set of players
func (self *SimEngine) CalcPayouts(players []goraph.Vertex) (Pc, Pd float64) {
  Nc := int32(0)
  Nd := int32(0)
  // count up cooperators and defectors
  for i := 0; i < len(players); i++ {
    if (self.agents[players[i]].cooperate) {
      Nc += 1
    } else {
      Nd += 1
    }
  }
  // calculate payouts
  Pd = float64(self.mult*self.cost*Nc)/float64(Nc+Nd)
  Pc = Pd - float64(self.cost)
  return Pc, Pd
}

// play a public goods game with the specified set of agents
func (self *SimEngine) PlayGame(players []goraph.Vertex) {
  // calculate payouts
  Pc, Pd := self.CalcPayouts(players)
  // distribute payouts
  for i := 0; i < len(players); i++ {
    player := self.agents[players[i]]
    if (player.cooperate) {
      player.payouts += Pc
    } else {
      player.payouts += Pd
    }
  }
}

// update the strategy of agent x based on the payouts
func (self *SimEngine) UpdateStrategy(x goraph.Vertex, y goraph.Vertex) {
  // get the agents
  agentx := self.agents[x]
  agenty := self.agents[y]

  // calculate the probability that x's strategy will be updated
  Pe := Fermi(self.betae, float64(agenty.payouts), float64(agentx.payouts))

  // update x's strategy if appropriate
  if ((RandProb(self.rnGen) < Pe) && (agentx.cooperate != agenty.cooperate)) {
    if (agentx.cooperate) {
      self.Nc -= 1
      self.Nd += 1
    } else {
      self.Nd -= 1
      self.Nc += 1
    }
    agentx.cooperate = agenty.cooperate
  }
}

// update the structure of the network based on the payouts
func (self *SimEngine) UpdateStructure(x goraph.Vertex, y goraph.Vertex) {
  fmt.Println("UPDATING STRUCTURE")
  // check to see if y is a cooperator
  agenty := self.agents[y]
  if (agenty.cooperate) {
    // link to cooperator is satisfactory - no network update required
    return
  }
  agentx := self.agents[x]

  // get the list of y's neighbors minus agent x
  Ny := self.graph.Neighbors(y)
  Ny = RemoveVertexFromSlice(Ny, x)

  // if x is y's last neighbor then structure update cannot be done
  if (len(Ny) <= 0) {
    return
  }

  // calculate the probability that x's link to y will be updated
  Pa := Fermi(self.betaa, float64(agentx.payouts), float64(agenty.payouts))

  // switch x's link with y if appropriate
  if (RandProb(self.rnGen) <= Pa) {
    // remove all shared neighbors
    Nx := self.graph.Neighbors(x)
    Ny = RemoveVerticesFromSlice(Ny, Nx)

    // vx's new neighbor
    var newy goraph.Vertex

    if (len(Ny) > 0) { // y has some neighbors that are not x's neighbors
      // select new neighbor from y's neighbors
      newy = Ny[RandInt(self.rnGen, int64(len(Ny)))]
    } else { // y doesn't have any neighbors that are not also x's neighbors
      // get all vertices except x, y and x's neighbors
      // -- Note that Nx includes y
      available := RemoveVerticesFromSlice(self.graph.Vertices(), append(Nx, x))

      // chose an newy if some vertices are available
      if (len(available) > 0) { // some nodes are available to choose from
        newy = available[RandInt(self.rnGen, int64(len(available)))]
      } else { // no nodes are available to choose from
        return
      }
    }

    // replace y with newy
    self.graph.RemoveEdge(x, y)
    self.graph.AddEdge(x, newy)
  }
}

// remove the duplicates from the slice
// see https://play.golang.org/p/HJEm6qy5wb
func removeDuplicates(a []goraph.Vertex) []goraph.Vertex {
	result := []goraph.Vertex{}
	seen := map[goraph.Vertex]goraph.Vertex{}
	for _, val := range a {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = val
		}
	}
	return result
}

func (self *SimEngine) String() string {
  s := "  {"
  s = fmt.Sprintf("%s\n  \"%v\":%d,", s, "ngens", self.numGens)
  s = fmt.Sprintf("%s\n  \"%v\":%d,", s, "nagents", self.numAgents)
  s = fmt.Sprintf("%s\n  \"%v\":%d,", s, "z", self.avgdeg)
  s = fmt.Sprintf("%s\n  \"%v\":%d,", s, "r", self.mult)
  s = fmt.Sprintf("%s\n  \"%v\":%d,", s, "cost", self.cost)
  s = fmt.Sprintf("%s\n  \"%v\":%.5f,", s, "betae", self.betae)
  s = fmt.Sprintf("%s\n  \"%v\":%.5f,", s, "betaa", self.betaa)
  s = fmt.Sprintf("%s\n  \"%v\":%.5f", s, "W", self.W)
  s = fmt.Sprintf("%s\n  }", s)
  return s
}

// write the header for the population statistics file
func (self *SimEngine) WritePStatsHeader(w io.Writer) {
  // write out headers
  fmt.Fprintf(w, "%s,%s,%s\n", "g", "Pc", "Pd")
}

// write population statistics for current gen to pstats file
func (self *SimEngine) WritePStats(w io.Writer, gen int32) {
  pNc := (float64(self.Nc)/float64(self.numAgents))
  pNd := (float64(self.Nd)/float64(self.numAgents))
  fmt.Fprintf(w,"%d,%5.3f,%5.3f\n", gen, pNc, pNd)
}

func (self *SimEngine) DegreeHistogramData(w io.Writer) {
  // write header
  fmt.Fprintf(w, "%s,%s,%s\n", "id", "S", "K")
  // write data
  var strat string
  for _, v := range self.graph.Vertices() {
    if (self.agents[v].cooperate) {
      strat = "C"
    } else {
      strat = "D"
    }
    fmt.Fprintf(w, "%v,%s,%d\n", v, strat, self.graph.Degree(v))
  }
}
