package simgpgg

import "testing"
import "testutil"
import "goraph"
import "math"

func NewTestSimEngine() *SimEngine {
  numAgents := int32(7)
  avgdeg := int32(4)
  numGens := int32(5)
  mult := int32(3)
  cost := int32(1)
  W := float64(0)
  betae := math.Inf(+1)
  betaa := math.Inf(+1)
  return NewSimEngine(numAgents, numGens, avgdeg, mult, cost, W, betae, betaa)
}

func TestPlayGame(u *testing.T) {
  simeng := NewTestSimEngine()
  vertices := simeng.graph.Vertices()

  // calculate payouts
  Pc, Pd := simeng.CalcPayouts(vertices)

  // play a game with all agents
  simeng.PlayGame(vertices)

  // make sure that each agent got the correct payout
  var agent *Agent
  var po float64
  for _, v := range vertices {
    agent = simeng.agents[v]
    if (agent.cooperate) {
      po = Pc
    } else {
      po = Pd
    }
    testutil.AssertFloat64Equal(u, agent.payouts, po)
  }
}

func TestUpdateStrategy(u *testing.T) {
  simeng := NewTestSimEngine()

  // select an agent x
  x := simeng.graph.Vertices()[0]
  // set agent x's payout to 99
  simeng.agents[x].payouts = 99
  // make sure x is a cooperator
  if (!simeng.agents[x].cooperate) {
    simeng.agents[x].cooperate = true
    simeng.Nc += 1
    simeng.Nd -= 1
  }
  // select one of its neighbors
  Nx := simeng.graph.Neighbors(x)
  y := Nx[0]
  // set y's payout equal to 100
  simeng.agents[y].payouts = 100
  // make sure y is a defectors
  if (simeng.agents[y].cooperate) {
    simeng.agents[y].cooperate = false
    simeng.Nc -= 1
    simeng.Nd += 1
  }

  nc := simeng.Nc
  nd := simeng.Nd

  // update structure and make sure that y is emoved from x's neighbors
  simeng.UpdateStrategy(x, y)

  testutil.AssertFalse(u, simeng.agents[x].cooperate)
  testutil.AssertFalse(u, simeng.agents[y].cooperate)
  testutil.AssertInt32Equal(u, simeng.Nc, nc-1)
  testutil.AssertInt32Equal(u, simeng.Nd, nd+1)
}

func TestUpdateStructure(u *testing.T) {
  simeng := NewTestSimEngine()

  // select an agent x
  x := simeng.graph.Vertices()[0]
  // set agent x's payout to 100
  simeng.agents[x].payouts = 100
  // select one of its neighbors
  Nx := simeng.graph.Neighbors(x)
  y := Nx[0]
  // set y's payout equal to 99
  simeng.agents[y].payouts = 99
  // make sure y is a defectors
  simeng.agents[y].cooperate = false

  // update structure and make sure that y is emoved from x's neighbors
  simeng.UpdateStructure(x, y)

  // get the list of neighbors
  Nx  = simeng.graph.Neighbors(x)
  Ny := simeng.graph.Neighbors(y)

  testutil.AssertFalse(u, goraph.VertexSlice(Nx).Contains(y))
  testutil.AssertFalse(u, goraph.VertexSlice(Ny).Contains(x))
}

func TestUpdateStructure2(u *testing.T) {
  simeng := NewTestSimEngine()
  vertices := simeng.graph.Vertices()

  // select an agent x
  x := vertices[0]
  // set agent x's payout to 100
  simeng.agents[x].payouts = 100
  // select an agent y
  y := vertices[1]
  // set y's payout equal to 99
  simeng.agents[y].payouts = 99
  // make sure y is a defectors
  simeng.agents[y].cooperate = false

  // --------------------------------------
  // test case when x is y's only neighbor

  // remove all x's neighbors
  Nx := simeng.graph.Neighbors(x)
  for i := 0; i < len(Nx); i++ {
    simeng.graph.RemoveEdge(x, Nx[i])
  }
  // remove all y's neighbors
  Ny := simeng.graph.Neighbors(y)
  for i := 0; i < len(Ny); i++ {
    simeng.graph.RemoveEdge(y, Ny[i])
  }
  Nx = simeng.graph.Neighbors(x)
  Ny = simeng.graph.Neighbors(y)

  Nx = simeng.graph.Neighbors(x)
  testutil.AssertIntEqual(u, len(Nx), 0)
  Ny = simeng.graph.Neighbors(y)
  testutil.AssertIntEqual(u, len(Ny), 0)

  // make x and y neighbors - the only neighbor they have
  simeng.graph.AddEdge(x, y)

  Nx = simeng.graph.Neighbors(x)
  testutil.AssertIntEqual(u, len(Nx), 1)
  goraph.AssertVertexEqual(u, Nx[0], y)

  Ny = simeng.graph.Neighbors(y)
  testutil.AssertIntEqual(u, len(Ny), 1)
  goraph.AssertVertexEqual(u, Ny[0], x)

  // update strategy - no changes should be made
  simeng.UpdateStructure(x, y)

  Nx = simeng.graph.Neighbors(x)
  testutil.AssertIntEqual(u, len(Nx), 1)
  goraph.AssertVertexEqual(u, Nx[0], y)

  Ny = simeng.graph.Neighbors(y)
  testutil.AssertIntEqual(u, len(Ny), 1)
  goraph.AssertVertexEqual(u, Ny[0], x)

  // --------------------------------------------------------------
  // test case when all of y's neighbors are alreagy x's neighbors

  // make sure x and y have the same neighbors except themselves
  var v, loner goraph.Vertex
  lonerAdded := false
  for i := 0; i < len(vertices); i++ {
    v = vertices[i]
    if ((v != x) && (v != y)) {
      if (!lonerAdded) {
        // first vertex tha is not x or y is the loner
        loner = v
        lonerAdded = true
      } else {
        // rest of vertices become neighbors of x and y
        simeng.graph.AddEdge(x, v)
        simeng.graph.AddEdge(y, v)
      }
    }
  }

  Nx = simeng.graph.Neighbors(x)
  testutil.AssertIntEqual(u, len(Nx), len(vertices)-2)
  testutil.AssertTrue(u, goraph.VertexSlice(Nx).Contains(y))

  Ny = simeng.graph.Neighbors(y)
  testutil.AssertIntEqual(u, len(Ny), len(vertices)-2)
  testutil.AssertTrue(u, goraph.VertexSlice(Ny).Contains(x))

  // update strategy - x should be paired with the loner agent
  simeng.UpdateStructure(x, y)

  Nx = simeng.graph.Neighbors(x)
  testutil.AssertIntEqual(u, len(Nx), len(vertices)-2)
  testutil.AssertFalse(u, goraph.VertexSlice(Nx).Contains(y))
  testutil.AssertTrue(u, goraph.VertexSlice(Nx).Contains(loner))

  Ny = simeng.graph.Neighbors(y)
  testutil.AssertIntEqual(u, len(Ny), len(vertices)-3)
  testutil.AssertFalse(u, goraph.VertexSlice(Ny).Contains(x))

  // --------------------------------------------------------------
  // test case when there are no vertices available to switch to

  // make x and y neighbors with all vertices
  // -- make y neighbors with loner
  simeng.graph.AddEdge(y, loner)
  // -- make x and y neighbors
  simeng.graph.AddEdge(y, x)

  Nx = simeng.graph.Neighbors(x)
  testutil.AssertIntEqual(u, len(Nx), len(vertices)-1)
  testutil.AssertTrue(u, goraph.VertexSlice(Nx).Contains(y))
  testutil.AssertTrue(u, goraph.VertexSlice(Nx).Contains(loner))

  Ny = simeng.graph.Neighbors(y)
  testutil.AssertIntEqual(u, len(Ny), len(vertices)-1)
  testutil.AssertTrue(u, goraph.VertexSlice(Ny).Contains(x))
  testutil.AssertTrue(u, goraph.VertexSlice(Ny).Contains(loner))

  // update strategy - nothing should change
  simeng.UpdateStructure(x, y)

  Nx = simeng.graph.Neighbors(x)
  testutil.AssertIntEqual(u, len(Nx), len(vertices)-1)
  testutil.AssertTrue(u, goraph.VertexSlice(Nx).Contains(y))
  testutil.AssertTrue(u, goraph.VertexSlice(Nx).Contains(loner))

  Ny = simeng.graph.Neighbors(y)
  testutil.AssertIntEqual(u, len(Ny), len(vertices)-1)
  testutil.AssertTrue(u, goraph.VertexSlice(Ny).Contains(x))
  testutil.AssertTrue(u, goraph.VertexSlice(Ny).Contains(loner))
}
