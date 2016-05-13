package simgpgg

import "goraph"
import "math"
import "math/rand"

// Remove the specified vertex from the slice
func RemoveVertexFromSlice(vertices []goraph.Vertex, v goraph.Vertex) []goraph.Vertex {
  // create a copy that will serve as the return value
  clone := make([]goraph.Vertex, len(vertices))
  copy(clone, vertices)

  // find vertex v and remove it
  for idx, candidate := range clone {
    if candidate == v {
      return append(clone[:idx], clone[idx+1:len(clone)]...)
		}
	}
  // slice doesn't contain v
  return clone
}

// Remove the specified vertices from the slice
func RemoveVerticesFromSlice(vertices []goraph.Vertex, toRemove []goraph.Vertex) []goraph.Vertex {
  // create a new slice that will serve as the return value
  var newSlice []goraph.Vertex

  // sort the two lists of vertices
  goraph.VertexSlice(vertices).Sort()
  goraph.VertexSlice(toRemove).Sort()

  canIdx := 0
  remIdx := 0
  var candidate, v goraph.Vertex
  for ;((remIdx < len(toRemove)) && (canIdx < len(vertices))); {
    v = toRemove[remIdx]
    candidate = vertices[canIdx]
    if (candidate < v) {
      // add the candidate to the new slice
      newSlice = append(newSlice, candidate)
      // get next candidate
      canIdx++
    } else if (candidate == v) {
      // don't add candidate to the new slice
      // get next candidate
      canIdx++
      // get next vertex to be removed
      remIdx++
    } else { // (candidate > v)
      // get next vertex to be removed
      remIdx++
    }
  }

  // add any remaining vertices to the new slice
  if (canIdx < len(vertices)) {
    newSlice = append(newSlice, vertices[canIdx:]...)
  }

  // return the slice with vertices removed
  return newSlice
}

// Create a regular ring graph with N nodes each with degree K
func NewRegularRing(N, K int32) *goraph.AdjacencyList {
  // create the nodes for the graph
  graph := goraph.NewAdjacencyList()
  for i := int32(0); i < N; i++ {
    graph.AddVertex()
  }
  // add the edges to the graph
  ub := float64(K)/float64(2)
  div := float64(N) - float64(1) - ub
  for i := int32(0); i < N; i++ {
    for j := i+1; j < N; j++ {
      diff := math.Abs(float64(i-j))
      mod := math.Mod(diff, div)
      if ((0 < mod) && (mod <= ub)) {
        graph.AddEdge(goraph.Vertex(i), goraph.Vertex(j))
      }
    }
  }
  return graph
}

// Create a homogeneous random graph by rewiring the edges of a regular ring
func NewHomoRandom(N, K int32, rnGen *rand.Rand) *goraph.AdjacencyList {
  // create a regular ring
  graph := NewRegularRing(N, K)
  // swapped he edges randomly
  edges := graph.Edges()
  var i,j int
  var ei, ej goraph.Edge
  var neiu, neju goraph.VertexSlice
  for ;len(edges)>0; {
    i = int(RandInt(rnGen, int64(len(edges))))
    j = int(RandInt(rnGen, int64(len(edges))))
    if (i != j) {
      ei = edges[i]
      ej = edges[j]
      // make sure circular links are not created
      if ((ei.U != ej.V) && (ej.U != ei.V)) {
        // make sure a duplicate edge will not be inserted
        neiu = graph.Neighbors(ei.U)
        neiu = graph.Neighbors(ej.U)
        if (!neiu.Contains(ej.V) && !neju.Contains(ei.V)) {
          // remove ei and ej
          graph.RemoveEdge(ei.U, ei.V)
          graph.RemoveEdge(ej.U, ej.V)
          // add swapped edges
          graph.AddEdge(ei.U, ej.V)
          graph.AddEdge(ej.U, ei.V)
        }
      }
    }
    // remove the edges from the list since they have been swapped
    edges = goraph.EdgeSlice(edges).Remove(i)
    edges = goraph.EdgeSlice(edges).Remove(j)
  }
  // return the resulting graph with randomized edges
  return graph
}

// Create a scale free network using the Barabasi-Albert algorithm
func NewScaleFreeNet(N, M0, M int32, rnGen *rand.Rand) *goraph.AdjacencyList {
  // M must be less than M0
  if (M0 < M) {
    panic("M0 is less than M")
  }
  // create an array that represents the roulette wheel
  // -- initial length is M0
  // -- capacity is M0 + (N - M0)*M
  wheel := make([]goraph.Vertex, M0 + 2*(N - M0)*M)
  wheelSize := 0

  // create M0 nodes to initially populate the graph
  graph := goraph.NewAdjacencyList()
  for i := int32(0); i < M0; i++ {
    // create node and give it a spot on the wheel
    wheel[wheelSize] = graph.AddVertex()
    wheelSize++
  }

  // add nodes to the graph until it contains N nodes
  // - when a node is added it is connected to M existing nodes
  // - nodes are selected randomly proportional to degree
  // - since the initial M0 nodes have zero degree to begin without, they were
  //   given a free spot on the whell above

  // create a slice to hold the randomly selected vertices
  selected := goraph.VertexSlice(make([]goraph.Vertex, 0, M))

  for i := M0; i < N; i++ {
    // add new node to the graph
    newNode := graph.AddVertex()

    // "clear" the list of selected nodes
    selected = selected[:0]

    // select M nodes without duplicates
    var v goraph.Vertex
    var found bool
    for j := int32(0); j < M; j++ {
      // randomly choose a vertex that isn't a duplicate
      for found = false; !found; {
        // spin the wheel to select a vertex proportional to its degree
        v = wheel[RandInt(rnGen, int64(wheelSize))]
        // check if v has already been selected
        found = !selected.Contains(v)
      }
      // grow the selected list
      selected = selected[:j+1]
      // add the selected vertex to the list
      selected[j] = v
    }

    // link newNode with the selected vertices
    // -- update the wheel to represent new node degrees
    for _, v := range selected {
      graph.AddEdge(newNode, v)
      // add newNode to wheel
      wheel[wheelSize] = newNode
      wheelSize++
      // add v to the wheel again
      wheel[wheelSize] = v
      wheelSize++
    }
  }

  return graph
}

// Create a scale free network using uniform attachment
func NewUniScaleFreeNet(N, M0, M int32, rnGen *rand.Rand) *goraph.AdjacencyList {
  // M must be less than M0
  if (M0 < M) {
    panic("M0 is less than M")
  }
  // create an array that represents the roulette wheel
  wheel := make([]goraph.Vertex, N)
  wheelSize := 0

  // create M0 nodes to initially populate the graph
  graph := goraph.NewAdjacencyList()
  for i := int32(0); i < M0; i++ {
    // create node and give it a spot on the wheel
    wheel[wheelSize] = graph.AddVertex()
    wheelSize++
  }

  // add nodes to the graph until it contains N nodes
  // - when a node is added it is connected to M existing nodes
  // - nodes are selected randomly with uniform probability

  // create a slice to hold the randomly selected vertices
  selected := goraph.VertexSlice(make([]goraph.Vertex, 0, M))

  for i := M0; i < N; i++ {
    // add new node to the graph
    newNode := graph.AddVertex()

    // "clear" the list of selected nodes
    selected = selected[:0]

    // select M nodes without duplicates
    var v goraph.Vertex
    var found bool
    for j := int32(0); j < M; j++ {
      // randomly choose a vertex that isn't a duplicate
      for found = false; !found; {
        // spin the wheel to select a vertex with uniform probability
        v = wheel[RandInt(rnGen, int64(wheelSize))]
        // check if v has already been selected
        found = !selected.Contains(v)
      }
      // grow the selected list
      selected = selected[:j+1]
      // add the selected vertex to the list
      selected[j] = v
    }

    // link newNode with the selected vertices
    for _, v := range selected {
      graph.AddEdge(newNode, v)
    }
    // add the new node to the wheel
    wheel[wheelSize] = newNode
    wheelSize++
  }

  return graph
}

// Create a Watts-Strogatz small world net
func NewSmallWorldNet(N, K int32, p float64, rnGen *rand.Rand) *goraph.AdjacencyList {
  // create a regular ring
  graph := NewRegularRing(N, K)
  // rewire the edges
  ub := float64(K)/float64(2)
  div := float64(N) - float64(1) - ub
  for i := int32(0); i < N; i++ {
    for j := i+1; j < N; j++ {
      diff := math.Abs(float64(i-j))
      mod := math.Mod(diff, div)
      if ((0 < mod) && (mod <= ub)) {
        // with probability p rewite the edge between Vertex(i) and Vertex(j)
        if (RandProb(rnGen) <= p) {
          // randomly select a replacement for Vertex(j)
          replace := RandInt(rnGen, int64(N))
          // make sure the proposed new edge will not be circular
          if (replace != int64(i)) {
            // make sure the proposed new edge is not a duplicate
            ni := goraph.VertexSlice(graph.Neighbors(goraph.Vertex(i)))
            if (!ni.Contains(goraph.Vertex(replace))) {
              // remove current edge
              graph.RemoveEdge(goraph.Vertex(i), goraph.Vertex(j))
              // add the new edge
              graph.AddEdge(goraph.Vertex(i), goraph.Vertex(replace))
            }
          }
        }
      }
    }
  }
  return graph
}
