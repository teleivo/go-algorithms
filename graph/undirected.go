// Package graph provides graph data structures and algorithm implementations.
// The APIs are following the APIs defined in the beautiful resource
// https://algs4.cs.princeton.edu/40graphs/ authored by Robert Sedgewick and
// Kevin Wayne. Some APIs will likely diverge since the Princeton algorithms
// course is written in Java. Sample test data is also taken from this
// resource.
//
// Parallel edges and self-loops are allowed.
package graph

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Undirected represents an undirected graph.
type Undirected struct {
	v   int
	e   int
	adj [][]int
}

// NewUndirected constructs an undirected graph holding v vertices. The graph
// does not contain any edges initially. Edges can be added with vertex ids are
// in the range of [0, V).
func NewUndirected(v int) (*Undirected, error) {
	if v < 0 {
		return nil, fmt.Errorf("number of vertices must be positive, V=%d", v)
	}
	adj := make([][]int, v)
	return &Undirected{v: v, adj: adj}, nil
}

// NewUndirectedFromReader constructs an undirected graph from the given reader.  The
// input needs to have the number of vertices (V) on its first line. Followed
// by the number of edges (E) on the second line. Every line thereafter represents
// an undirected edge containing both vertices separated by a space. A vertex
// is represented by a number in the range of [0, V).
//
// A short example of the format is:
// 4
// 2
// 0 1
// 2 3
//
// The above represents an undirected graph with 4 vertices and 2 edges.
// The first edge connects vertex 0 and 1, while the second connects vertex 2
// and 3.
func NewUndirectedFromReader(r io.Reader) (*Undirected, error) {
	s := bufio.NewScanner(r)
	if !s.Scan() {
		return nil, errors.New("first line missing, must contain number of vertices")
	}
	t := s.Text()
	v, err := strconv.Atoi(t)
	if err != nil {
		return nil, fmt.Errorf("failed to parse number of vertices, invalid token %q", t)
	}
	if v < 0 {
		return nil, fmt.Errorf("number of vertices must be positive, V=%d", v)
	}
	if !s.Scan() {
		return nil, errors.New("second line missing, must contain number of edges")
	}
	t = s.Text()
	e, err := strconv.Atoi(t)
	if err != nil {
		return nil, fmt.Errorf("failed to parse number of vertices, invalid token %q", t)
	}
	if e < 0 {
		return nil, fmt.Errorf("number of edges must be positive, V=%d", v)
	}
	adj := make([][]int, v)
	u := &Undirected{v: v, adj: adj}
	var edges int
	for s.Scan() {
		t = s.Text()
		n := strings.Split(t, " ")
		if len(n) != 2 {
			return nil, fmt.Errorf("edge must have two vertices, invalid line %q", t)
		}
		v, err := strconv.Atoi(n[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse vertex %q, invalid line %q", n[0], t)
		}
		if err := u.validateVertex(v); err != nil {
			return nil, err
		}
		w, err := strconv.Atoi(n[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse vertex %q, invalid line %q", n[1], t)
		}
		if err := u.validateVertex(w); err != nil {
			return nil, err
		}
		u.AddEdge(v, w)
		if err != nil {
			return nil, err
		}
		edges++
	}
	if edges != e {
		return nil, fmt.Errorf("number of edges %d is not equal to edges in list %d", u.e, edges)
	}
	return u, nil
}

// V returns the number of vertices of the graph.
func (u *Undirected) V() int {
	return u.v
}

// E returns the number of edges of the graph.
func (u *Undirected) E() int {
	return u.e
}

func (u *Undirected) validateVertex(v int) error {
	if v < 0 || v > u.V()-1 {
		return fmt.Errorf("vertex id must be within 0 and %d, invalid vertex %d", u.V()-1, v)
	}
	return nil
}

// AddEdge adds the given undirected edge to the graph.
func (u *Undirected) AddEdge(v, w int) {
	if err := u.validateVertex(v); err != nil {
		panic(err)
	}
	if err := u.validateVertex(w); err != nil {
		panic(err)
	}
	u.adj[v] = append(u.adj[v], w)
	u.adj[w] = append(u.adj[w], v)
	u.e++
}

// Adj returns the vertices adjacent to the given vertex.
func (u *Undirected) Adj(v int) []int {
	if err := u.validateVertex(v); err != nil {
		panic(err)
	}
	return u.adj[v]
}

func (u *Undirected) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d vertices, %d edges\n", u.v, u.e)
	for v := 0; v < u.v; v++ {
		fmt.Fprintf(&sb, "%d: ", v)
		for _, w := range u.Adj(v) {
			sb.WriteString(strconv.Itoa(w))
			sb.WriteString(" ")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// DepthFirstPaths finds a path from a source vertex to every other vertex in
// an undirected graph using depth-first search.
type DepthFirstPaths struct {
	g      *Undirected
	s      int
	marked []bool
	edgeTo []int
}

func NewDepthFirstPaths(g *Undirected, s int) (*DepthFirstPaths, error) {
	if s < 0 || s >= g.V() {
		return nil, fmt.Errorf("source vertex must be within range [0,%d), invalid vertex %d", g.V(), s)
	}
	dfp := &DepthFirstPaths{
		g:      g,
		s:      s,
		marked: make([]bool, g.V()),
		edgeTo: make([]int, g.V()),
	}
	dfp.marked[s] = true
	dfp.dfs(s)
	return dfp, nil
}

func (dfp *DepthFirstPaths) dfs(v int) {
	for _, w := range dfp.g.Adj(v) {
		if !dfp.marked[w] {
			dfp.edgeTo[w] = v
			dfp.marked[w] = true
			dfp.dfs(w)
		}
	}
}

// HasPathTo returns true if there is a path from the source vertex to the
// given vertex. Otherwise it will return false.
func (dfp *DepthFirstPaths) HasPathTo(v int) bool {
	if err := dfp.g.validateVertex(v); err != nil {
		panic(err)
	}
	return dfp.marked[v]
}

// PathTo returns the path from the source vertex to the given vertex. It
// returns nil if there is no such path.
func (dfp *DepthFirstPaths) PathTo(v int) []int {
	if !dfp.HasPathTo(v) {
		return nil
	}
	var p []int
	for {
		p = append(p, v)
		if v == dfp.s {
			break
		}
		v = dfp.edgeTo[v]
	}
	for i, j := 0, len(p)-1; i < j; i, j = i+1, j-1 {
		p[i], p[j] = p[j], p[i]
	}
	return p
}
