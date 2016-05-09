/**
 * A modified version of goraph from https://github.com/echlebek/goraph
 * with bug fixes for undirected graphs
 */
package goraph

import (
	"sort"
)

// Graph is implemented by all of the graph types. All of the graph
// algorithms use this data type instead of the concrete types.
type Graph interface {
	// AddVertex creates an returns a new vertex in the graph.
	AddVertex() Vertex

	// RemoveVertex permanently removes a vertex from the graph.
	RemoveVertex(v Vertex)

	// AddEdge adds an edge between u and v. If the graph is directional,
	// then the edge will go from u to v.
	AddEdge(u, v Vertex)

	// RemoveEdge removes the edge between u and v.
	RemoveEdge(u, v Vertex)

	// Vertices returns a slice of the graph's vertices.
	Vertices() []Vertex

	// Edges returns a slice of the graph's edges.
	Edges() []Edge

	// Neighbours returns a slice of the vertices that neighbour v.
	Neighbors(v Vertex) []Vertex
}

var (
	_ Graph = &AdjacencyList{}
)

// Vertex represents a node in the graph. Users should create
// new Vertex values with AddVertex.
type Vertex int

// Edge represents an edge between two vertices.
// In a directed graph, the edge is from U to V.
type Edge struct{ U, V Vertex }

// AdjacencyList implements an undirected graph using an adjacency list.
type AdjacencyList struct {
	edges      map[Vertex][]Vertex
	nextVertex Vertex
}

// NewAdjacencyList creates an empty graph.
func NewAdjacencyList() *AdjacencyList {
	return &AdjacencyList{edges: make(map[Vertex][]Vertex)}
}

func (g *AdjacencyList) AddVertex() Vertex {
	v := g.nextVertex
	g.edges[v] = make([]Vertex, 0)
	g.nextVertex++
	return v
}

func (g *AdjacencyList) RemoveVertex(v Vertex) {
	delete(g.edges, v)
	for vtx, vertices := range g.edges {
		for idx, candidate := range vertices {
			if candidate == v {
				g.edges[vtx] = append(vertices[:idx], vertices[idx+1:len(vertices)]...)
			}
		}
	}
}

func (g *AdjacencyList) AddEdge(u, v Vertex) {
	g.addHalfEdge(u,v)
	g.addHalfEdge(v,u)
}

func (g *AdjacencyList) addHalfEdge(u, v Vertex) {
	edges := g.edges[u]
	g.edges[u] = append(edges, v)
}

func (g *AdjacencyList) RemoveEdge(u, v Vertex) {
	g.removeHalfEdge(u,v)
	g.removeHalfEdge(v,u)
}

func (g *AdjacencyList) removeHalfEdge(u, v Vertex) {
	vertices, ok := g.edges[u]
	if !ok {
		return
	}

	for idx, vtx := range vertices {
		if vtx == v {
			// Remove the edge
			g.edges[u] = append(vertices[:idx], vertices[idx+1:len(vertices)]...)
			break
		}
	}
}

// VertexSlice is a convenience type for sorting vertices by ID.
type VertexSlice []Vertex

func (p VertexSlice) Len() int           { return len(p) }
func (p VertexSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p VertexSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p VertexSlice) Sort()              { sort.Sort(p) }

// EdgeSlice is a convenience type for sorting edges by ID.
type EdgeSlice []Edge

func (p EdgeSlice) Len() int { return len(p) }
func (p EdgeSlice) Less(i, j int) bool {
	if p[i].U == p[j].U {
		return p[i].V < p[j].V
	} else {
		return p[i].U < p[j].U
	}
}
func (p EdgeSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p EdgeSlice) Sort()         { sort.Sort(p) }

func (g *AdjacencyList) Vertices() []Vertex {
	vertices := make(VertexSlice, len(g.edges))
	var i int
	for k := range g.edges {
		vertices[i] = k
		i++
	}
	return vertices
}

func (g *AdjacencyList) Edges() []Edge {
	var edges []Edge
	for k, neighbors := range g.edges {
		for _, n := range neighbors {
			// to prevent duplicates, only add the edge if (k < n)
			if (k < n) {
				edges = append(edges, Edge{k, n})
			}
		}
	}
	return edges
}

func (g *AdjacencyList) Neighbors(v Vertex) []Vertex {
	return g.edges[v]
}
