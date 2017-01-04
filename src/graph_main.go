package main

import "errors"
import "fmt"
import "graphs"
import "graphs/alg/shortest_path"
import gs "graphs/adjacency_list_weighted"

type concretGraph struct {
    gs.BidirectionalGraph
    names []string
}

func (g *concretGraph) NumVertices() int { return len(g.BidirectionalGraph.IncidenceGraph) }
func (g *concretGraph) VertexName(v graphs.VertexID) string {
    if v.IsNil() {
        return "nil"
    }
    return g.names[int(v)]
}
// Implement WeightedGraphs
func (g *concretGraph) Weight(e graphs.EdgeIter) float64 {
    return g.BidirectionalGraph.IncidenceGraph[int(e.Source)].OutEdges[int(e.EdgeID)].Weight
}

var count = 0

func VisitVertex(event_name string, g *concretGraph, v graphs.VertexID) error {
    fmt.Printf("%.2d : %20s  %s\n", count, event_name, g.VertexName(v))
    count++
    if count > 100 { return errors.New("too many iterations") }
    return nil
}

func VisitEdge(event_name string, g *concretGraph, e graphs.EdgeIter) error {
    fmt.Printf("%.2d : %20s  %s -> %s\n", count, event_name, g.VertexName(e.Source), g.VertexName(e.Target))
    count++
    if count > 100 { return errors.New("too many iterations") }
    return nil
}

func bfs_visit(g *concretGraph, color_map graphs.ColorMap) {
    fmt.Println("\n\nBFS")
    count = 0
    err := graphs.BreadthFirstSearch(g, color_map, nil, []graphs.VertexID{ 0 }, graphs.BfsVisitor{
        OnInitializeVertex:    func (v graphs.VertexID) error { return VisitVertex("OnInitializeVertex", g, v) },
        OnExamineVertex:       func (v graphs.VertexID) error { return VisitVertex("OnExamineVertex", g, v) },
        OnExamineEdge:         func (e graphs.EdgeIter) error { return VisitEdge  ("OnExamineEdge", g, e) },
        OnTreeEdge:            func (e graphs.EdgeIter) error { return VisitEdge  ("OnTreeEdge", g, e) },
        OnDiscoverVertex:      func (v graphs.VertexID) error { return VisitVertex("OnDiscoverVertex", g, v) },
        OnNonTreeEdge:         func (e graphs.EdgeIter) error { return VisitEdge  ("OnNonTreeEdge", g, e) },
        OnGrayTargetEdge:      func (e graphs.EdgeIter) error { return VisitEdge  ("OnGrayTargetEdge", g, e) },
        OnBlackTargeEdge:      func (e graphs.EdgeIter) error { return VisitEdge  ("OnBlackTargeEdge", g, e) },
        OnFinishVertex:        func (v graphs.VertexID) error { return VisitVertex("OnFinishVertex", g, v) },
    })
    if err != nil { fmt.Println(err.Error()) }
}

func dfs_visit(g *concretGraph, color_map graphs.ColorMap) {
    fmt.Println("\n\nDFS")
    count = 0
    err := graphs.DepthFirstSearch(g, color_map, nil, 0, graphs.DfsVisitor{
        OnInitializeVertex:    func (v graphs.VertexID) error { return VisitVertex("OnInitializeVertex", g, v) },
        OnStartVertex:         func (v graphs.VertexID) error { return VisitVertex("OnStartVertex", g, v) },
        OnExamineEdge:         func (e graphs.EdgeIter) error { return VisitEdge  ("OnExamineEdge", g, e) },
        OnTreeEdge:            func (e graphs.EdgeIter) error { return VisitEdge  ("OnTreeEdge", g, e) },
        OnDiscoverVertex:      func (v graphs.VertexID) error { return VisitVertex("OnDiscoverVertex", g, v) },
        OnBackEdge:            func (e graphs.EdgeIter) error { return VisitEdge  ("OnBackEdge", g, e) },
        OnForwardOrCrossEdge:  func (e graphs.EdgeIter) error { return VisitEdge  ("OnForwardOrCrossEdge", g, e) },
        OnFinishNonTreeEdge:   func (e graphs.EdgeIter) error { return VisitEdge  ("OnFinishNonTreeEdge", g, e) },
        OnFinishVertex:        func (v graphs.VertexID) error { return VisitVertex("OnFinishVertex", g, v) },
        OnFinishTreeEdge:      func (e graphs.EdgeIter) error { return VisitEdge  ("OnFinishTreeEdge", g, e) },
    })
    if err != nil { fmt.Println(err.Error()) }
}

func dijkstra(g *concretGraph, color_map graphs.ColorMap) {
    fmt.Println("\n\ndijkstra")
    pmap := make(graphs.VertexMap, g.NumVertices())
    dmap := make(graphs.DistanceMap, g.NumVertices())
    imap := make(graphs.IndexMap, g.NumVertices())

    count = 0
    err := shortest_path.DijkstraShortestPaths(g, color_map, pmap, dmap, imap, []graphs.VertexID{ 0 }, shortest_path.DijkstraVisitor{
        OnInitializeVertex:    func (v graphs.VertexID) error { return VisitVertex("OnInitializeVertex", g, v) },
        OnExamineVertex:       func (v graphs.VertexID) error { return VisitVertex("OnExamineVertex", g, v) },
        OnExamineEdge:         func (e graphs.EdgeIter) error { return VisitEdge  ("OnExamineEdge", g, e) },
        OnDiscoverVertex:      func (v graphs.VertexID) error { return VisitVertex("OnDiscoverVertex", g, v) },
        OnEdgeRelaxed:         func (e graphs.EdgeIter) error { return VisitEdge  ("OnEdgeRelaxed", g, e) },
        OnEdgeNonRelaxed:      func (e graphs.EdgeIter) error { return VisitEdge  ("OnEdgeNonRelaxed", g, e) },
        OnFinishVertex:        func (v graphs.VertexID) error { return VisitVertex("OnFinishVertex", g, v) },
    })
    if err != nil { fmt.Println(err.Error()) }
    fmt.Println("pmap: ", pmap)
    fmt.Println("dmap: ", dmap)
}

func dijkstraNoColorMap(g *concretGraph) {
    fmt.Println("\n\ndijkstraNoColorMap")
    pmap := make(graphs.VertexMap, g.NumVertices())
    dmap := make(graphs.DistanceMap, g.NumVertices())
    imap := make(graphs.IndexMap, g.NumVertices())

    count = 0
    err := shortest_path.DijkstraShortestPathsNoColorMap(g, pmap, dmap, imap, []graphs.VertexID{ 0 }, shortest_path.DijkstraVisitor{
        OnInitializeVertex:    func (v graphs.VertexID) error { return VisitVertex("OnInitializeVertex", g, v) },
        OnExamineVertex:       func (v graphs.VertexID) error { return VisitVertex("OnExamineVertex", g, v) },
        OnExamineEdge:         func (e graphs.EdgeIter) error { return VisitEdge  ("OnExamineEdge", g, e) },
        OnDiscoverVertex:      func (v graphs.VertexID) error { return VisitVertex("OnDiscoverVertex", g, v) },
        OnEdgeRelaxed:         func (e graphs.EdgeIter) error { return VisitEdge  ("OnEdgeRelaxed", g, e) },
        OnEdgeNonRelaxed:      func (e graphs.EdgeIter) error { return VisitEdge  ("OnEdgeNonRelaxed", g, e) },
        OnFinishVertex:        func (v graphs.VertexID) error { return VisitVertex("OnFinishVertex", g, v) },
    })
    if err != nil { fmt.Println(err.Error()) }
    fmt.Println("pmap: ", pmap)
    fmt.Println("dmap: ", dmap)
}

func astar(g *concretGraph, color_map graphs.ColorMap) {
    fmt.Println("\n\nastar")
    pmap := make(graphs.VertexMap, g.NumVertices())
    dmap := make(graphs.DistanceMap, g.NumVertices())
    cost := make(graphs.DistanceMap, g.NumVertices())
    imap := make(graphs.IndexMap, g.NumVertices())

    count = 0
    goal := graphs.VertexID(7)
    err := shortest_path.AStarSearch(g, color_map, pmap, dmap, cost, imap, 0, func (v graphs.VertexID) float64 {
            if v == goal {
                return 0.0
            }
            return 0.5
        }, shortest_path.AStarVisitor{
        OnInitializeVertex:    func (v graphs.VertexID) error { return VisitVertex("OnInitializeVertex", g, v) },
        OnExamineVertex:       func (v graphs.VertexID) error { return VisitVertex("OnExamineVertex", g, v) },
        OnExamineEdge:         func (e graphs.EdgeIter) error { return VisitEdge  ("OnExamineEdge", g, e) },
        OnDiscoverVertex:      func (v graphs.VertexID) error { return VisitVertex("OnDiscoverVertex", g, v) },
        OnEdgeRelaxed:         func (e graphs.EdgeIter) error { return VisitEdge  ("OnEdgeRelaxed", g, e) },
        OnEdgeNonRelaxed:      func (e graphs.EdgeIter) error { return VisitEdge  ("OnEdgeNonRelaxed", g, e) },
        OnBlackTargetEdge:     func (e graphs.EdgeIter) error { return VisitEdge  ("OnBlackTargeEdge", g, e) },
        OnFinishVertex:        func (v graphs.VertexID) error { return VisitVertex("OnFinishVertex", g, v) },
    })
    if err != nil { fmt.Println(err.Error()) }
    fmt.Println("pmap: ", pmap)
    fmt.Println("dmap: ", dmap)
}

func astar_tree(g *concretGraph) {
    fmt.Println("\n\nastar tree")
    pmap := make(graphs.VertexMap, g.NumVertices())
    dmap := make(graphs.DistanceMap, g.NumVertices())
    cost := make(graphs.DistanceMap, g.NumVertices())
    imap := make(graphs.IndexMap, g.NumVertices())

    count = 0
    goal := graphs.VertexID(7)
    err := shortest_path.AStarSearchTree(g, pmap, dmap, cost, imap, 0, func (v graphs.VertexID) float64 {
            if v == goal {
                return 0.0
            }
            return 0.0
        }, shortest_path.AStarVisitor{
        OnInitializeVertex:    func (v graphs.VertexID) error { return VisitVertex("OnInitializeVertex", g, v) },
        OnExamineVertex:       func (v graphs.VertexID) error { return VisitVertex("OnExamineVertex", g, v) },
        OnExamineEdge:         func (e graphs.EdgeIter) error { return VisitEdge  ("OnExamineEdge", g, e) },
        OnDiscoverVertex:      func (v graphs.VertexID) error { return VisitVertex("OnDiscoverVertex", g, v) },
        OnEdgeRelaxed:         func (e graphs.EdgeIter) error { return VisitEdge  ("OnEdgeRelaxed", g, e) },
        OnEdgeNonRelaxed:      func (e graphs.EdgeIter) error { return VisitEdge  ("OnEdgeNonRelaxed", g, e) },
        OnBlackTargetEdge:     func (e graphs.EdgeIter) error { return VisitEdge  ("OnBlackTargeEdge", g, e) },
        OnFinishVertex:        func (v graphs.VertexID) error { return VisitVertex("OnFinishVertex", g, v) },
    })
    if err != nil { fmt.Println(err.Error()) }
    fmt.Println("pmap: ", pmap)
    fmt.Println("dmap: ", dmap)
}

const (
    A graphs.VertexID = iota
    B
    C
    D
    E
    F
    G
    H
)

func main() {
    g := &concretGraph{
        BidirectionalGraph: gs.MakeBidirectionalGraph(gs.IncidenceGraph{
            { OutEdges: []gs.Edge{ { B, 1.0 }, { C, 1.0 } } },
            { OutEdges: []gs.Edge{ { A, 1.0 }, { D, 1.0 }, { E, 1.0 } } },
            { OutEdges: []gs.Edge{ { F, 3.0 }, { G, 1.0 } } },
            { OutEdges: []gs.Edge{ } },
            { OutEdges: []gs.Edge{ { F, 1.0 }, { H, 2.0 } } },
            { OutEdges: []gs.Edge{ } },
            { OutEdges: []gs.Edge{ { H, 1.0 } } },
            { OutEdges: []gs.Edge{ { A, 1.0 } } },
        }),
        names: []string { "A", "B", "C", "D", "E", "F", "G", "H" },
    }
    color_map := make(graphs.ColorMap, g.NumVertices())

    bfs_visit(g, color_map)
    dfs_visit(g, color_map)
    dijkstra(g, color_map)
    dijkstraNoColorMap(g)
    astar(g, color_map)
    astar_tree(g)
}
