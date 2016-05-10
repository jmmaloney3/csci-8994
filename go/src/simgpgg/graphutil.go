package simgpgg

import "goraph"
import "math"

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
