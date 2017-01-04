package topological_sort

import "graphs"

func TopologicalSort(
    g graphs.VertexIncidenceGraph,
    color_map graphs.ColorMap,
    outputer func (v graphs.VertexID) error) error {
        
    return graphs.DepthFirstSearch(g, color_map, nil, graphs.NilVertex, graphs.DfsVisitor{
        OnBackEdge: func (e graphs.EdgeIter) error { return graphs.NotDAG{} },
        OnFinishVertex: outputer,
    })
}