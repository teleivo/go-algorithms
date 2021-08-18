# graph

## Benchmarks

If you are unfamiliar with benchmarks in Go I recommend you read [How to write benchmarks in Go](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go) first.

Current benchmarks on my machine comparing my recursive and iterative
implementations of depth-first search on an undirected graph.

```
go test -run=XXX -bench=. -benchtime=20s
goos: linux
goarch: amd64
pkg: github.com/teleivo/go-algorithms/graph
cpu: Intel(R) Core(TM) i5-2520M CPU @ 2.50GHz
BenchmarkDepthFirstPathsMediumGraph-4            	 2820112	      9273 ns/op
BenchmarkDepthFirstPathsLargeGraph-4             	      63	 337273120 ns/op
BenchmarkIterativeDepthFirstPathsMediumGraph-4   	  798840	     36035 ns/op
BenchmarkIterativeDepthFirstPathsLargeGraph-4    	      46	 502699563 ns/op
```
