package simgpgg

import "testing"
import "testutil"
import "goraph"

func TestRemoveVerticesFromSlice(u *testing.T) {
  var vertices []goraph.Vertex
  vertices = append(vertices, 1)
  vertices = append(vertices, 5)
  vertices = append(vertices, 9)
  vertices = append(vertices, 13)
  vertices = append(vertices, 17)

  var toRemove []goraph.Vertex
  toRemove = append(toRemove, 2)
  toRemove = append(toRemove, 5)
  toRemove = append(toRemove, 10)
  toRemove = append(toRemove, 11)
  toRemove = append(toRemove, 12)
  toRemove = append(toRemove, 13)
  toRemove = append(toRemove, 14)

  newSlice := RemoveVerticesFromSlice(vertices, toRemove)

  // make sure new slice has correct elements
  testutil.AssertIntEqual(u, len(newSlice), 3)
  goraph.AssertVertexEqual(u, newSlice[0], 1)
  goraph.AssertVertexEqual(u, newSlice[1], 9)
  goraph.AssertVertexEqual(u, newSlice[2], 17)

  // make sure vertices has not changed
  testutil.AssertIntEqual(u, len(vertices), 5)
  goraph.AssertVertexEqual(u, vertices[0], 1)
  goraph.AssertVertexEqual(u, vertices[1], 5)
  goraph.AssertVertexEqual(u, vertices[2], 9)
  goraph.AssertVertexEqual(u, vertices[3], 13)
  goraph.AssertVertexEqual(u, vertices[4], 17)
}

func TestRemoveVertexFromSlice(u *testing.T) {
  var vertices []goraph.Vertex
  vertices = append(vertices, 1)
  vertices = append(vertices, 5)
  vertices = append(vertices, 9)
  vertices = append(vertices, 13)

  v := goraph.Vertex(9)

  newSlice := vertices
  for i := 0; i < 2; i++ {
    newSlice = RemoveVertexFromSlice(newSlice, v)

    // make sure new slice has correct elements
    testutil.AssertIntEqual(u, len(newSlice), 3)
    goraph.AssertVertexEqual(u, newSlice[0], 1)
    goraph.AssertVertexEqual(u, newSlice[1], 5)
    goraph.AssertVertexEqual(u, newSlice[2], 13)

    // make sure vertices is unchanged
    testutil.AssertIntEqual(u, len(vertices), 4)
    goraph.AssertVertexEqual(u, vertices[0], 1)
    goraph.AssertVertexEqual(u, vertices[1], 5)
    goraph.AssertVertexEqual(u, vertices[2], 9)
    goraph.AssertVertexEqual(u, vertices[3], 13)
  }
}

func TestRemoveEdgeFromSlice(u *testing.T) {
  graph := NewRegularRing(4,2)

  edges := graph.Edges()

  remEdge := edges[3]


  edges = goraph.EdgeSlice(edges).Remove(3)

  found := goraph.EdgeSlice(edges).Contains(remEdge)
  testutil.AssertFalse(u, found)
}

func TestGraphGen(u *testing.T) {
  rnGen := NewRandNumGen()
  // TODO: actually test that these graphs are correctly constructed
  NewRegularRing(64,4)
  NewHomoRandom(64,4, rnGen)
  NewScaleFreeNet(64, 4, 4, rnGen)
  NewUniScaleFreeNet(64, 4, 4, rnGen)
  NewSmallWorldNet(64, 4, float64(0), rnGen)
  NewSmallWorldNet(64, 4, float64(0.5), rnGen)
  NewSmallWorldNet(64, 4, float64(1.0), rnGen)
}
