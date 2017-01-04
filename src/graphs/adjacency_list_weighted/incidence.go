package adjacency_list_weighted

import "graphs"

type Edge struct {
    Target graphs.VertexID
    Weight float64
}

type Vertex struct {
    Deleted  bool
    OutEdges []Edge
}

type IncidenceGraph []Vertex

// Implement the VertexList interface
func (g *IncidenceGraph) VertexBegin() graphs.VertexID {
    for i, n := graphs.VertexID(0), graphs.VertexID(len(*g)); i < n; i++ {
        if !(*g)[i].Deleted { return i }
    }
    return graphs.NilVertex
}
func (g *IncidenceGraph) VertexNext(v graphs.VertexID) graphs.VertexID {
    for i, n := v + 1, graphs.VertexID(len(*g)); i < n; i++ {
        if !(*g)[int(i)].Deleted { return i }
    }
    return graphs.NilVertex
}

// Implement the IncidenceGraph interface
func (g *IncidenceGraph) OutEdgeBegin(v graphs.VertexID) graphs.EdgeIter {
    vertex := (*g)[int(v)]
    if vertex.Deleted {
        return graphs.NilEdgeIter()
    }
    outEdges := vertex.OutEdges
    for i, n := graphs.EdgeID(0), graphs.EdgeID(len(outEdges)); i < n; i++ {
        if outEdges[i].Target.IsValid() && !(*g)[int(outEdges[i].Target)].Deleted {
            return graphs.EdgeIter{ Source: v, Target: outEdges[i].Target, EdgeID: i }
        }
    }
    return graphs.EdgeIter{ Source: v, Target: graphs.NilVertex, EdgeID: graphs.NilEdgeID }
}
func (g *IncidenceGraph) OutEdgeNext(e graphs.EdgeIter) graphs.EdgeIter {
    vertex := (*g)[int(e.Source)]
    if vertex.Deleted {
        return graphs.EdgeIter{ Source: e.Source, Target: graphs.NilVertex, EdgeID: graphs.NilEdgeID }
    }
    outEdges := vertex.OutEdges
    for i, n := e.EdgeID + 1, graphs.EdgeID(len(outEdges)); i < n; i++ {
        if outEdges[i].Target.IsValid() && !(*g)[int(outEdges[i].Target)].Deleted {
            return graphs.EdgeIter{ Source: e.Source, Target: outEdges[i].Target, EdgeID: i }
        }
    }
    return graphs.EdgeIter{ Source: e.Source, Target: graphs.NilVertex, EdgeID: graphs.NilEdgeID }
}

// Implement the EdgeList interface
func (g *IncidenceGraph) EdgeBegin() graphs.EdgeIter {
    v := g.VertexBegin()
    if v.IsNil() { return graphs.NilEdgeIter() }
    e := g.OutEdgeBegin(v)
    if e.Target.IsValid() { return e }
    return g.EdgeNext(e)
}

func (g *IncidenceGraph) EdgeNext(e graphs.EdgeIter) graphs.EdgeIter {
    if e.Target.IsValid() {
        e = g.OutEdgeNext(e)
    }
    for e.Source.IsValid() && e.Target.IsNil() {
        v := g.VertexNext(e.Source);
        if v.IsNil() { return graphs.NilEdgeIter() }
        e = g.OutEdgeBegin(v)
    }
    return e
}

// Implement Mutable interface
func (g *IncidenceGraph) AddVertex() graphs.VertexID {
    *g = append(*g, Vertex{})
    return graphs.VertexID(len(*g))
}

func (g *IncidenceGraph) ClearVertex(v graphs.VertexID) {
    outEdges := (*g)[int(v)].OutEdges;
    for e := 0; e < len(outEdges); e++ {
        outEdges[e].Target = graphs.NilVertex;
    }
    (*g)[int(v)].OutEdges = outEdges[0:0]
    
    for i, n := 0, len(*g); i < n; i++ {
        if (*g)[i].Deleted {
            continue
        }
        outEdges := (*g)[i].OutEdges
        for e := 0; e < len(outEdges); e++ {
            if outEdges[e].Target == v {
                outEdges[e].Target = graphs.NilVertex
            }
        }
    }
}

func (g *IncidenceGraph) RemoveVertex(v graphs.VertexID) {
    (*g)[int(v)].Deleted = true
}
    
func (g *IncidenceGraph) AddEdge(u graphs.VertexID, v graphs.VertexID) (e graphs.EdgeIter, inserted bool) {
    outEdges := (*g)[int(u)].OutEdges
    (*g)[int(u)].OutEdges = append(outEdges, Edge{Target: v})
    return graphs.EdgeIter{ Source: u, Target: v, EdgeID: graphs.EdgeID(len(outEdges)) }, true
}

func (g *IncidenceGraph) RemoveEdgeBetween(u graphs.VertexID, v graphs.VertexID) {
    for e := g.OutEdgeBegin(u); e.EdgeID.IsValid(); e = g.OutEdgeNext(e) {
        if e.Target == v {
            (*g)[int(u)].OutEdges[int(e.EdgeID)].Target = graphs.NilVertex
        }
    }
}

func (g *IncidenceGraph) RemoveEdge(e graphs.EdgeIter) {
    if e.EdgeID.IsValid() {
        (*g)[int(e.Source)].OutEdges[int(e.EdgeID)].Target = graphs.NilVertex
    }
}

func (g *IncidenceGraph) RemoveOutEdge(e graphs.EdgeIter) {
    if e.EdgeID.IsValid() {
        (*g)[int(e.Source)].OutEdges[int(e.EdgeID)].Target = graphs.NilVertex
    }
}

func (g *IncidenceGraph) RemoveOutEdgeIf(u graphs.VertexID, p func (e graphs.EdgeIter) bool) {
    for e := g.OutEdgeBegin(u); e.EdgeID.IsValid(); e = g.OutEdgeNext(e) {
        if p(e) {
            (*g)[int(u)].OutEdges[int(e.EdgeID)].Target = graphs.NilVertex
        }
    }
}

func (g *IncidenceGraph) Compress(vmover func (src graphs.VertexID, dst graphs.VertexID),
                                  emover func (src graphs.EdgeIter, dst graphs.EdgeIter)) []graphs.VertexID {
    vmap := g.compressVertices(vmover)
    g.compressOutEdges(vmap, emover)
    return vmap
}

func (g *IncidenceGraph) compressVertices(vmover func (src graphs.VertexID, dst graphs.VertexID)) []graphs.VertexID {
    vmap := make([]graphs.VertexID, len(*g))
    j := 0;
    for i, n := 0, len(*g); i < n; i++ {
        if (*g)[i].Deleted {
            vmap[i] = graphs.NilVertex
            continue
        }
        if j < i {
            if vmover != nil {
                vmover(graphs.VertexID(i), graphs.VertexID(j))
            }
            (*g)[j] = (*g)[i]
            vmap[i] = graphs.VertexID(j)
            j++
        }
    }
    for i := j; i < len(vmap); i++ {
        if vmover != nil {
            vmover(graphs.VertexID(i), graphs.NilVertex)
        }
        (*g)[i] = Vertex{}
    }
    (*g) = (*g)[0:j]
    return vmap
}

func (g *IncidenceGraph) compressOutEdges(vmap []graphs.VertexID, emover func (src graphs.EdgeIter, dst graphs.EdgeIter)) {
    for i, n := graphs.VertexID(0), graphs.VertexID(len(*g)); i < n; i++ {
        outEdges := (*g)[int(i)].OutEdges
        j := graphs.EdgeID(0);
        for k, m := graphs.EdgeID(0), graphs.EdgeID(len(outEdges)); k < m; k++ {
            vt := outEdges[int(k)].Target
            if vt.IsNil() || vmap[int(vt)].IsNil() {
                continue
            }
            if j < k {
                if emover != nil {
                    emover(graphs.EdgeIter{ Source: i, Target: vt, EdgeID: k },
                           graphs.EdgeIter{ Source: i, Target: vmap[int(vt)], EdgeID: j })
                }
                outEdges[int(j)].Target = vmap[int(vt)]
                j++
            }
        }
        for k, m := j, graphs.EdgeID(len(outEdges)); k < m; k++ {
            vt := outEdges[int(k)].Target
            if emover != nil {
                emover(graphs.EdgeIter{ Source: i, Target: vt, EdgeID: k },
                       graphs.EdgeIter{ Source: i, Target: graphs.NilVertex, EdgeID: graphs.NilEdgeID })
            }
            outEdges[int(k)] = Edge{Target: graphs.NilVertex}
        }
        (*g)[int(i)].OutEdges = outEdges[0:int(j)]
    }
}
