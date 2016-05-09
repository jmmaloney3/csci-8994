package simgpgg

import "goraph"
//import "fmt"
import "math"

// Create a regular ring graph with N nodes each with degree K
func NewRegularRing(N, K int) *goraph.AdjacencyList {
  // create the nodes for the graph
  graph := goraph.NewAdjacencyList()
  for i := int(0); i < N; i++ {
    graph.AddVertex()
  }
  // add the edges to the graph
  ub := float64(K)/float64(2)
  div := float64(N) - float64(1) - ub
  for i := int(0); i < N; i++ {
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
