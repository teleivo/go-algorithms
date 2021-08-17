package graph_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/teleivo/go-algorithms/graph"
)

func TestNewUndirected(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		g, err := graph.NewUndirected(2)

		if err != nil {
			t.Fatalf("NewUndirected(%d) returned an unexpected error: %s", 2, err)
		}
		if want, got := 2, g.V(); got != 2 {
			t.Errorf("NewUndirected(%d).V() got %d, want %d", want, got, want)
		}
	})

	t.Run("NumberOfVerticesIsNegative", func(t *testing.T) {
		_, err := graph.NewUndirected(-1)

		if err == nil {
			t.Fatalf("NewUndirected(%d) expected an error but got none", -1)
		}
	})
}

func TestNewUndirectedFromReader(t *testing.T) {
	tt := []struct {
		name string
		in   string
		v    int
		e    int
		adj  [][]int
	}{
		{
			name: "GraphWithOneEdge",
			in:   "2\n1\n0 1",
			v:    2,
			e:    1,
			adj:  [][]int{0: {1}, 1: {0}},
		},
	}

	for _, tt := range tt {
		t.Run(tt.name, func(t *testing.T) {
			g, err := graph.NewUndirectedFromReader(strings.NewReader(tt.in))

			if err != nil {
				t.Fatalf("NewUndirectedFromReader(%q) returned an unexpected error: %s", tt.in, err)
			}
			if got := g.V(); got != tt.v {
				t.Errorf("NewUndirectedFromReader(%q).V() got %d, want %d", tt.in, got, tt.v)
			}
			if got := g.E(); got != tt.e {
				t.Errorf("NewUndirectedFromReader(%q).E() got %d, want %d", tt.in, got, tt.e)
			}
			for v := range tt.adj {
				got := g.Adj(v)
				if diff := cmp.Diff(tt.adj[v], got); diff != "" {
					t.Errorf("Adj(%d) mismatch (-want +got):\n%s", v, diff)
				}
			}
		})
	}

	tt = []struct {
		name string
		in   string
		v    int
		e    int
		adj  [][]int
	}{
		{
			name: "NumberOfVerticesMissing",
			in:   "",
		},
		{
			name: "VertexNotANumber",
			in:   "a",
		},
		{
			name: "VertexNegative",
			in:   "-1",
		},
		{
			name: "NumberOfEdgesMissing",
			in:   "1\n",
		},
		{
			name: "NumberOfEdgesIsNotANumber",
			in:   "1\na",
		},
		{
			name: "NumberOfEdgesIsNegative",
			in:   "1\n-1",
		},
		{
			name: "NumberOfEdgesNotEqualToEdgeList",
			in:   "2\n1",
		},
		{
			name: "FirstVertexIsOutOfRange",
			in:   "2\n1\n3 1",
		},
		{
			name: "FirstVertexIsNotANumber",
			in:   "2\n1\na 1",
		},
		{
			name: "SecondVertexIsOutOfRange",
			in:   "2\n1\n0 3",
		},
		{
			name: "SecondVertexIsNotANumber",
			in:   "2\n1\n0 a",
		},
		{
			name: "EdgeMissesAVertex",
			in:   "2\n1\n3",
		},
	}

	for _, tt := range tt {
		t.Run(tt.name, func(t *testing.T) {
			_, err := graph.NewUndirectedFromReader(strings.NewReader(tt.in))

			if err == nil {
				t.Fatalf("NewUndirectedFromReader(%q) expected an error but got none", tt.in)
			}
		})
	}
}

func TestUndirectedString(t *testing.T) {

	tt := []struct {
		name  string
		graph *graph.Undirected
		want  string
	}{
		{
			name: "GraphWithoutEdges",
			graph: func() *graph.Undirected {
				g, _ := graph.NewUndirected(8)
				return g
			}(),
			want: "8 vertices, 0 edges\n0: \n1: \n2: \n3: \n4: \n5: \n6: \n7: \n",
		},
		{
			name: "GraphWithEdges",
			graph: func() *graph.Undirected {
				g, _ := graph.NewUndirected(8)
				g.AddEdge(0, 1)
				g.AddEdge(0, 6)
				g.AddEdge(0, 7)
				g.AddEdge(1, 6)
				g.AddEdge(1, 2)
				g.AddEdge(1, 4)
				return g
			}(),
			want: "8 vertices, 6 edges\n0: 1 6 7 \n1: 0 6 2 4 \n2: 1 \n3: \n4: 1 \n5: \n6: 0 1 \n7: 0 \n",
		},
	}

	for _, tt := range tt {
		if diff := cmp.Diff(tt.want, tt.graph.String()); diff != "" {
			t.Errorf("String() mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestDepthFirstPaths(t *testing.T) {
	tt := []struct {
		name        string
		graph       *graph.Undirected
		source      int
		hasPathTo   []int
		hasNoPathTo []int
		pathTo      map[int][]int
	}{
		{
			name: "GraphWithTwoComponents",
			graph: func() *graph.Undirected {
				g, _ := graph.NewUndirected(4)
				g.AddEdge(0, 1)
				g.AddEdge(2, 3)
				return g
			}(),
			source:      0,
			hasPathTo:   []int{0, 1},
			hasNoPathTo: []int{2, 3},
			pathTo: map[int][]int{
				0: {0},
				1: {0, 1},
			},
		},
		{
			name: "GraphWithSelfLoop",
			graph: func() *graph.Undirected {
				g, _ := graph.NewUndirected(2)
				g.AddEdge(0, 1)
				g.AddEdge(0, 0)
				return g
			}(),
			source:    0,
			hasPathTo: []int{0, 1},
			pathTo: map[int][]int{
				0: {0},
				1: {0, 1},
			},
		},
		{
			name: "GraphWithParallelEdge",
			graph: func() *graph.Undirected {
				g, _ := graph.NewUndirected(2)
				g.AddEdge(0, 1)
				g.AddEdge(1, 0)
				return g
			}(),
			source:    0,
			hasPathTo: []int{0, 1},
			pathTo: map[int][]int{
				0: {0},
				1: {0, 1},
			},
		},
		{
			name: "GraphWithCycles",
			graph: func() *graph.Undirected {
				g, _ := graph.NewUndirected(8)
				g.AddEdge(0, 1)
				g.AddEdge(0, 6)
				g.AddEdge(0, 7)
				g.AddEdge(1, 6)
				g.AddEdge(1, 2)
				g.AddEdge(1, 4)
				g.AddEdge(2, 3)
				g.AddEdge(2, 4)
				g.AddEdge(3, 4)
				g.AddEdge(4, 5)
				return g
			}(),
			source:      0,
			hasPathTo:   []int{0, 1, 2, 3, 4, 5, 6, 7},
			hasNoPathTo: []int{},
			pathTo: map[int][]int{
				0: {0},
				1: {0, 1},
				2: {0, 1, 2},
				3: {0, 1, 2, 3},
				4: {0, 1, 2, 3, 4},
				5: {0, 1, 2, 3, 4, 5},
				6: {0, 1, 6},
				7: {0, 7},
			},
		},
	}

	for _, tt := range tt {
		t.Run(tt.name, func(t *testing.T) {
			dfp, err := graph.NewDepthFirstPaths(tt.graph, tt.source)

			if err != nil {
				t.Fatalf("NewDepthFirstPaths(%v, %d) returned an unexpected error: %s", tt.graph, tt.source, err)
			}

			for _, v := range tt.hasPathTo {
				if !dfp.HasPathTo(v) {
					t.Errorf("HasPathTo(%d) found no path to source %d but should have", v, tt.source)
				}
			}
			for _, v := range tt.hasNoPathTo {
				if dfp.HasPathTo(v) {
					t.Errorf("HasPathTo(%d) found a path to source %d but should not have", v, tt.source)
				}
			}
			for v, want := range tt.pathTo {
				if diff := cmp.Diff(want, dfp.PathTo(v)); diff != "" {
					t.Errorf("PathTo(%d) mismatch (-want +got):\n%s", v, diff)
				}
			}
		})
	}

}

func TestIterativeDepthFirstPaths(t *testing.T) {
	tt := []struct {
		name        string
		graph       *graph.Undirected
		source      int
		hasPathTo   []int
		hasNoPathTo []int
		pathTo      map[int][]int
	}{
		{
			name: "GraphWithTwoComponents",
			graph: func() *graph.Undirected {
				g, _ := graph.NewUndirected(4)
				g.AddEdge(0, 1)
				g.AddEdge(2, 3)
				return g
			}(),
			source:      0,
			hasPathTo:   []int{0, 1},
			hasNoPathTo: []int{2, 3},
			pathTo: map[int][]int{
				0: {0},
				1: {0, 1},
			},
		},
		{
			name: "GraphWithSelfLoop",
			graph: func() *graph.Undirected {
				g, _ := graph.NewUndirected(2)
				g.AddEdge(0, 1)
				g.AddEdge(0, 0)
				return g
			}(),
			source:    0,
			hasPathTo: []int{0, 1},
			pathTo: map[int][]int{
				0: {0},
				1: {0, 1},
			},
		},
		{
			name: "GraphWithParallelEdge",
			graph: func() *graph.Undirected {
				g, _ := graph.NewUndirected(2)
				g.AddEdge(0, 1)
				g.AddEdge(1, 0)
				return g
			}(),
			source:    0,
			hasPathTo: []int{0, 1},
			pathTo: map[int][]int{
				0: {0},
				1: {0, 1},
			},
		},
		{
			name: "GraphWithCycles",
			graph: func() *graph.Undirected {
				g, _ := graph.NewUndirected(8)
				g.AddEdge(0, 1)
				g.AddEdge(0, 6)
				g.AddEdge(0, 7)
				g.AddEdge(1, 6)
				g.AddEdge(1, 2)
				g.AddEdge(1, 4)
				g.AddEdge(2, 3)
				g.AddEdge(2, 4)
				g.AddEdge(3, 4)
				g.AddEdge(4, 5)
				return g
			}(),
			source:      0,
			hasPathTo:   []int{0, 1, 2, 3, 4, 5, 6, 7},
			hasNoPathTo: []int{},
			pathTo: map[int][]int{
				0: {0},
				1: {0, 1},
				2: {0, 1, 2},
				3: {0, 1, 2, 3},
				4: {0, 1, 2, 3, 4},
				5: {0, 1, 2, 3, 4, 5},
				6: {0, 1, 6},
				7: {0, 7},
			},
		},
	}

	for _, tt := range tt {
		t.Run(tt.name, func(t *testing.T) {
			dfp, err := graph.NewIterativeDepthFirstPaths(tt.graph, tt.source)

			if err != nil {
				t.Fatalf("NewIterativeDepthFirstPaths(%v, %d) returned an unexpected error: %s", tt.graph, tt.source, err)
			}

			for _, v := range tt.hasPathTo {
				if !dfp.HasPathTo(v) {
					t.Errorf("HasPathTo(%d) found no path to source %d but should have", v, tt.source)
				}
			}
			for _, v := range tt.hasNoPathTo {
				if dfp.HasPathTo(v) {
					t.Errorf("HasPathTo(%d) found a path to source %d but should not have", v, tt.source)
				}
			}
			for v, want := range tt.pathTo {
				if diff := cmp.Diff(want, dfp.PathTo(v)); diff != "" {
					t.Errorf("PathTo(%d) mismatch (-want +got):\n%s", v, diff)
				}
			}
		})
	}

}
