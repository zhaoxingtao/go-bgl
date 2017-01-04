package adjacency_list_weighted

import "graphs"

type BidirectionalGraph struct {
    IncidenceGraph
    VertexInEdges [][]graphs.VertexID
}

func MakeBidirectionalGraph(g IncidenceGraph) BidirectionalGraph {
    bg := BidirectionalGraph{ g, make([][]graphs.VertexID, len(g)) }
    for e := g.EdgeBegin(); e.EdgeID.IsValid(); e = g.EdgeNext(e) {
        bg.VertexInEdges[int(e.Target)] = append(bg.VertexInEdges[int(e.Target)], e.Source)
    }
    return bg
}

// Implement the BidirectionalGraph interface
func (g *BidirectionalGraph) InEdgeBegin(v graphs.VertexID) graphs.EdgeIter {
    if g.IncidenceGraph[int(v)].Deleted {
        return graphs.NilEdgeIter()
    }
    inEdges := g.VertexInEdges[int(v)]
    for i, n := graphs.EdgeID(0), graphs.EdgeID(len(inEdges)); i < n; i++ {
        if inEdges[i].IsValid() && !g.IncidenceGraph[int(inEdges[i])].Deleted {
            return graphs.EdgeIter{ Source: v, Target: inEdges[i], EdgeID: i }
        }
    }
    return graphs.EdgeIter{ Source: v, Target: graphs.NilVertex, EdgeID: graphs.NilEdgeID }
}
func (g *BidirectionalGraph) InEdgeNext(e graphs.EdgeIter) graphs.EdgeIter {
    if g.IncidenceGraph[int(e.Source)].Deleted {
        return graphs.EdgeIter{ Source: e.Source, Target: graphs.NilVertex, EdgeID: graphs.NilEdgeID }
    }
    inEdges := g.VertexInEdges[int(e.Source)]
    for i, n := e.EdgeID + 1, graphs.EdgeID(len(inEdges)); i < n; i++ {
        if inEdges[i].IsValid() && !g.IncidenceGraph[int(inEdges[i])].Deleted {
            return graphs.EdgeIter{ Source: e.Source, Target: inEdges[i], EdgeID: i }
        }
    }
    return graphs.EdgeIter{ Source: e.Source, Target: graphs.NilVertex, EdgeID: graphs.NilEdgeID }
}

// Implement Mutable interface
func (g *BidirectionalGraph) AddVertex() graphs.VertexID {
    v := g.IncidenceGraph.AddVertex()
    g.VertexInEdges = append(g.VertexInEdges, nil)
    return v
}

func (g *BidirectionalGraph) ClearVertex(v graphs.VertexID) {
    g.IncidenceGraph.ClearVertex(v)
    inEdges := g.VertexInEdges[int(v)];
    for e := 0; e < len(inEdges); e++ {
        inEdges[e] = graphs.NilVertex;
    }
    g.VertexInEdges[int(v)] = inEdges[0:0]

    for i, n := 0, len(g.VertexInEdges); i < n; i++ {
        if g.IncidenceGraph[i].Deleted {
            continue
        }
        inEdges := g.VertexInEdges[i]
        for e := 0; e < len(inEdges); e++ {
            if inEdges[e] == v {
                inEdges[e] = graphs.NilVertex
            }
        }
    }
}

func (g *BidirectionalGraph) AddEdge(u graphs.VertexID, v graphs.VertexID) (e graphs.EdgeIter, inserted bool) {
    e, inserted = g.IncidenceGraph.AddEdge(u, v)
    if inserted {
        g.VertexInEdges[int(v)] = append(g.VertexInEdges[int(v)], u)
    }
    return e, inserted
}

func (g *BidirectionalGraph) RemoveEdgeBetween(u graphs.VertexID, v graphs.VertexID) {
    g.IncidenceGraph.RemoveEdgeBetween(u, v)
    inEdges := g.VertexInEdges[v]
    for i := 0; i < len(inEdges); i++ {
        if inEdges[i] == u {
            inEdges[i] = graphs.NilVertex
        }
    }
}

func (g *BidirectionalGraph) RemoveEdge(e graphs.EdgeIter) {
    if e.EdgeID.IsValid() {
        return
    }
    g.IncidenceGraph.RemoveEdge(e)

    inEdges := g.VertexInEdges[int(e.Target)]
    for i := 0; i < len(inEdges); i++ {
        if inEdges[i] == e.Source {
            inEdges[i] = graphs.NilVertex
            // we do not know which outEdge is mapping to which inEdge if
            // there are multiple edges between two vertices. So here we
            // delete the first inEdge.
            return
        }
    }
}

func (g *BidirectionalGraph) RemoveInEdge(e graphs.EdgeIter) {
    if e.EdgeID.IsValid() {
        g.VertexInEdges[int(e.Source)][int(e.EdgeID)] = graphs.NilVertex
    }
}

func (g *BidirectionalGraph) RemoveInEdgeIf(u graphs.VertexID, p func (e graphs.EdgeIter) bool) {
    for e := g.InEdgeBegin(u); e.EdgeID.IsValid(); e = g.InEdgeNext(e) {
        if p(e) {
            g.VertexInEdges[int(e.Source)][int(e.EdgeID)] = graphs.NilVertex
        }
    }
}

func (g *BidirectionalGraph) Compress(vmover func (src graphs.VertexID, dst graphs.VertexID),
                                      emover func (src graphs.EdgeIter, dst graphs.EdgeIter)) []graphs.VertexID {
    vmap := g.IncidenceGraph.Compress(vmover, emover)

    j := 0
    for i, n := 0, len(g.VertexInEdges); i < n; i++ {
        if t := vmap[i]; t.IsValid() {
            g.VertexInEdges[int(t)] = g.VertexInEdges[i]
            if j < int(t) {
                j = int(t)
            }
        }
    }
    g.VertexInEdges = g.VertexInEdges[0:j]

    for i, n := 0, len(g.VertexInEdges); i < n; i++ {
        inEdges := g.VertexInEdges[i]
        j := 0;
        for k := 0; k < len(inEdges); k++ {
            if vk := inEdges[k]; vk.IsValid() {
                inEdges[j] = vk;
                j++
            }
        }
        g.VertexInEdges[i] = g.VertexInEdges[i][0:j]
    }
    return vmap
}
