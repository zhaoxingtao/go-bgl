package graphs

import "fmt"

const debug = true

type VertexID int
type EdgeID int

const (
    NilVertex VertexID = -1
    NilEdgeID EdgeID = -1
)

func (v VertexID) IsNil() bool { return v < 0 }
func (v VertexID) IsValid() bool { return v >= 0 }
func (e EdgeID) IsNil() bool { return e < 0 }
func (e EdgeID) IsValid() bool { return e >= 0 }

type EdgeIter struct {
    Source VertexID
    Target VertexID
    EdgeID EdgeID
}
func NilEdgeIter() EdgeIter { return EdgeIter{ NilVertex, NilVertex, NilEdgeID } }

type IncidenceGraph interface {
    // Even when nilEdge, the EdgeIter.Source should be specified.
    OutEdgeBegin(v VertexID) EdgeIter
    OutEdgeNext(e EdgeIter) EdgeIter
}

type BidirectionalGraph interface {
    IncidenceGraph
    // Even when nilEdge, the EdgeIter.Source should be specified.
    InEdgeBegin(v VertexID) EdgeIter
    InEdgeNext(e EdgeIter) EdgeIter
}

type VertexList interface {
    VertexBegin() VertexID
    VertexNext(v VertexID) VertexID
}

type EdgeListGraph interface {
    EdgeBegin() EdgeIter
    EdgeNext(e EdgeIter) EdgeIter
}

type VertexIncidenceGraph interface {
    VertexList
    IncidenceGraph
}

type VertexEdgeListGraph interface {
    VertexList
    EdgeListGraph
}

type WeightedGraph interface {
    Weight(e EdgeIter) float64
}

type IncidenceWeightedGraph interface {
    IncidenceGraph
    WeightedGraph
}

type VertexIncidenceWeightedGraph interface {
    VertexList
    IncidenceWeightedGraph
}

type EdgeListWeightedGraph interface {
    EdgeListGraph
    WeightedGraph
}

type VertexEdgeListWeightedGraph interface {
    VertexEdgeListGraph
    WeightedGraph
}

type AdjacencyMatrixGraph interface {
    Edge(v VertexID, u VertexID) EdgeIter
}

type MutableGraph interface {
    AddVertex() VertexID
    ClearVertex(v VertexID)
    RemoveVertex(v VertexID)
    
    AddEdge(u VertexID, v VertexID) (e EdgeIter, inserted bool)
    RemoveEdgeBetween(u VertexID, v VertexID)
    RemoveEdge(e EdgeIter)
}

type MutableIncidenceGraph interface {
    MutableGraph
    RemoveOutEdge(e EdgeIter)
    RemoveOutEdgeIf(u VertexID, p func (e EdgeIter) bool)
}

type MutableBidirectionalGraph interface {
    MutableIncidenceGraph
    RemoveInEdge(e EdgeIter) error
    RemoveInEdgeIf(u VertexID, p func (e EdgeIter) bool)
}

type NegativeEdgeWeight EdgeIter
func (e NegativeEdgeWeight) Error() string {
    return fmt.Sprintf("edge(s: %v; t: %v; mega: %v) has negative weight.", e.Source, e.Target, e.EdgeID)
}

type NotDAG struct{}
func (e NotDAG) Error() string {
    return "Graph is not a Directed Acyclic Graph."
}

type NoSuchVertexError VertexID
func (e NoSuchVertexError) Error() string {
    return fmt.Sprintf("vertex (%v) does not exist.", VertexID(e))
}

type NoSuchEdgeError EdgeIter
func (e NoSuchEdgeError) Error() string {
    return fmt.Sprintf("edge(s: %v; t: %v; id: %v) does not exist.", e.Source, e.Target, e.EdgeID)
}

