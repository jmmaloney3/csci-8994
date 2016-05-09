package simgpgg

import "math"
import "math/rand"
import "fmt"
import "goraph"

// A simulation engine for simulating the public goods games
// played among agents occupying the nodes of a graph.
type SimEngine struct {
  graph goraph.Graph // the graph that holds the agents
  numAgents int32    // total number of agents being simulated
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
func NewSimEngine(numAgents int32, numGens int32, mult int32, cost int32, W float64,
                  betae float64, betaa float64) *SimEngine {
  graph := NewRegularRing(int(numAgents), 6)
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

  return &SimEngine { numAgents: numAgents, numGens: numGens, mult: mult,
                      rnGen: rnGen, graph: graph, agents: agents, cost: cost,
                      Nc: Nc, Nd: Nd, W: W, betae: betae, betaa: betaa }
}

func (self *SimEngine) RunSim() {
  // print out simulation parameters
  fmt.Printf("PGG Graph Simulator:\n")
  fmt.Printf("  Num Agents: %4d\n", self.numAgents)
  fmt.Printf("  Num Gens:   %4d\n", self.numGens)
  fmt.Printf("  Mult (r):   %4d\n", self.mult)

  // calculate probability that a structure update occurs
  stratUpdProb := float64(1)/(float64(1) + self.W)

  // write out headers
  fmt.Printf("%s,%s,%s\n", "g", "Nc", "Nd")

  // loop until one strategy is eliminated or the max num of gens is reached
  for g := int32(0); (!self.SimComplete(g)); g++ {
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
    sponsors = append(sponsors, Nx...)
    sponsors = append(sponsors, Ny...)
    sponsors = removeDuplicates(sponsors)
    // cldar payouts for agents that will play games
    for i := 0; i < len(sponsors); i++ {
      self.agents[sponsors[i]].payouts = float64(0)
    }
    // play the games
    for i := 0; i < len(sponsors); i++ {
      sponsor := sponsors[i]
      players := append(self.graph.Neighbors(sponsor), sponsor)
      self.PlayGame(players)
    }
    // print out payouts
    //fmt.Printf("Player x: %6.4f\n", self.agents[x].payouts)
    //fmt.Printf("Player y: %6.4f\n", self.agents[y].payouts)

    if (RandProb(self.rnGen) < stratUpdProb) {
      // update agent strategy - if appropriate
      self.UpdateStrategy(x, y)
    } else {
      // update network structure - if appropriate
      self.UpdateStructure(x, y)
    }
    // write out stats
    fmt.Printf("%d,%d,%d\n", g, self.Nc, self.Nd)
  }
}

func (self *SimEngine) SimComplete(genNum int32) bool {
  return (genNum >= self.numGens) || (self.Nc >= self.numAgents) || (self.Nd >= self.numAgents)
}

// play a public goods game with the specified set of agents
func (self *SimEngine) PlayGame(players []goraph.Vertex) {
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
  Pd := float64(self.mult*self.cost*Nc)/float64(Nc+Nd)
  Pc := Pd - float64(self.cost)
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
  exp := (-self.betae)*float64(agenty.payouts - agentx.payouts)
  Pe  := float64(1)/(float64(1) + math.Exp(exp))
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
  // check to see if y is a cooperator
  agenty := self.agents[y]
  if (agenty.cooperate) {
    // no network update required
    return
  }
  agentx := self.agents[x]

  // if y only has one neighbor then a structure update cannot be performed
  Ny := self.graph.Neighbors(y)
  if (len(Ny) <= 1) {
    return
  }

  // calculate the probability that x's link to y will be updated
  exp := (-self.betaa)*float64(agentx.payouts - agenty.payouts)
  Pa  := float64(1)/(float64(1) + math.Exp(exp))
  if (RandProb(self.rnGen) < Pa) {
    fmt.Printf("x: %v\n", x)
    fmt.Printf("old Ny: %v\n", Ny)
    fmt.Printf("y: %v\n", y)
    fmt.Printf("old Nx: %v\n", self.graph.Neighbors(x))

    // select new neighbor - make sure its not x - or one of x's neighbors
    newy := x
    for ;newy == x; {
      newy = Ny[RandInt(self.rnGen, int64(len(Ny)))]
    }

    self.graph.RemoveEdge(x, y)
    self.graph.AddEdge(x, newy)
    fmt.Printf("new Ny: %v\n", self.graph.Neighbors(y))
    fmt.Printf("newy: %v\n", newy)
    fmt.Printf("new Nx: %v\n", self.graph.Neighbors(x))
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
