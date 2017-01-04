package shortest_path

import "graphs"
import "math"

type DijkstraVisitor struct {
    OnInitializeVertex func (v graphs.VertexID) error
    OnExamineVertex    func (v graphs.VertexID) error
    OnExamineEdge      func (e graphs.EdgeIter) error
    OnDiscoverVertex   func (v graphs.VertexID) error
    OnEdgeRelaxed      func (e graphs.EdgeIter) error
    OnEdgeNonRelaxed   func (e graphs.EdgeIter) error
    OnFinishVertex     func (v graphs.VertexID) error
}

func DijkstraShortestPaths(
    g graphs.VertexIncidenceWeightedGraph,
    cmap graphs.ColorMap,
    pmap graphs.VertexMap,
    dmap graphs.DistanceMap,
    imap graphs.IndexMap,
    sources []graphs.VertexID,
    visitor DijkstraVisitor) error {

    d_inf := math.Inf(1)
    for v := g.VertexBegin(); v.IsValid(); v = g.VertexNext(v) {
        if visitor.OnInitializeVertex != nil {
            if err := visitor.OnInitializeVertex(v); err != nil { return err }
        }
        cmap.SetColor(v, graphs.WHITE)
        if pmap != nil { pmap.SetVertex(v, v) }
        dmap.SetDistance(v, d_inf)
        imap.SetIndex(v, 0);
    }

    return DijkstraShortestPathsNoInit(g, cmap, pmap, dmap, imap, sources, visitor)
}

func DijkstraShortestPathsNoInit(
    g graphs.IncidenceWeightedGraph,
    cmap graphs.ColorMap,
    pmap graphs.VertexMap,
    dmap graphs.DistanceMap,
    imap graphs.IndexMap,
    sources []graphs.VertexID,
    visitor DijkstraVisitor) error {

    queue := NewVertexPriorityQueue(imap, dmap)

    for _, v := range sources {
        cmap.SetColor(v, graphs.GRAY)
        dmap.SetDistance(v, 0.0)
        queue.Push(v)
        if visitor.OnDiscoverVertex != nil {
            if err := visitor.OnDiscoverVertex(v); err != nil { return err }
        }
    }
    
    for queue.Len() > 0 {
        u := queue.Pop()
        if visitor.OnExamineVertex != nil {
            if err := visitor.OnExamineVertex(u); err != nil { return err }
        }
        d_u  := dmap.GetDistance(u)
        for e := g.OutEdgeBegin(u); e.EdgeID.IsValid(); e = g.OutEdgeNext(e) {
            if visitor.OnExamineEdge != nil {
                if err := visitor.OnExamineEdge(e); err != nil { return err }
            }
            color := cmap.GetColor(e.Target)
            if color == graphs.WHITE {
                cmap.SetColor(e.Target, graphs.GRAY)
                if pmap != nil { pmap.SetVertex(e.Target, u) }
                
                w_uv := g.Weight(e)
                if w_uv < 0.0 { return graphs.NegativeEdgeWeight(e) }
                dmap.SetDistance(e.Target, w_uv + d_u)
                if visitor.OnEdgeRelaxed != nil {
                    if err := visitor.OnEdgeRelaxed(e); err != nil { return err }
                }

                if visitor.OnDiscoverVertex != nil {
                    if err := visitor.OnDiscoverVertex(e.Target); err != nil { return err }
                }
                queue.Push(e.Target)
            } else if color == graphs.GRAY {
                w_uv := g.Weight(e)
                if w_uv < 0.0 { return graphs.NegativeEdgeWeight(e) }
                d_v := dmap.GetDistance(e.Target)
                if w_uv + d_u < d_v {
                    dmap.SetDistance(e.Target, w_uv + d_u)
                    if pmap != nil { pmap.SetVertex(e.Target, u) }
                    if visitor.OnEdgeRelaxed != nil {
                        if err := visitor.OnEdgeRelaxed(e); err != nil { return err }
                    }
                    queue.InsertOrUpdate(e.Target)
                } else {
                    if visitor.OnEdgeNonRelaxed != nil {
                        if err := visitor.OnEdgeNonRelaxed(e); err != nil { return err }
                    }
                }
            }   
        }
        cmap.SetColor(u, graphs.BLACK)
        if visitor.OnFinishVertex != nil {
            if err := visitor.OnFinishVertex(u); err != nil { return err }
        }
    }
    
    return nil
}

func DijkstraShortestPathsNoColorMap(
    g graphs.VertexIncidenceWeightedGraph,
    pmap graphs.VertexMap,
    dmap graphs.DistanceMap,
    imap graphs.IndexMap,
    sources []graphs.VertexID,
    visitor DijkstraVisitor) error {

    d_inf := math.Inf(1)
    for v := g.VertexBegin(); v.IsValid(); v = g.VertexNext(v) {
        if visitor.OnInitializeVertex != nil {
            if err := visitor.OnInitializeVertex(v); err != nil { return err }
        }
        if pmap != nil { pmap.SetVertex(v, v) }
        dmap.SetDistance(v, d_inf)
        imap.SetIndex(v, 0);
    }

    return DijkstraShortestPathsNoInitNoColorMap(g, pmap, dmap, imap, sources, visitor)
}

func DijkstraShortestPathsNoInitNoColorMap(
    g graphs.IncidenceWeightedGraph,
    pmap graphs.VertexMap,
    dmap graphs.DistanceMap,
    imap graphs.IndexMap,
    sources []graphs.VertexID,
    visitor DijkstraVisitor) error {

    d_inf := math.Inf(1)
    queue := NewVertexPriorityQueue(imap, dmap)

    for _, v := range sources {
        dmap.SetDistance(v, 0.0)
        queue.Push(v)
        if visitor.OnDiscoverVertex != nil {
            if err := visitor.OnDiscoverVertex(v); err != nil { return err }
        }
    }
    
    for queue.Len() > 0 {
        u := queue.Pop()
        if visitor.OnExamineVertex != nil {
            if err := visitor.OnExamineVertex(u); err != nil { return err }
        }
        d_u  := dmap.GetDistance(u)
        for e := g.OutEdgeBegin(u); e.EdgeID.IsValid(); e = g.OutEdgeNext(e) {
            if visitor.OnExamineEdge != nil {
                if err := visitor.OnExamineEdge(e); err != nil { return err }
            }
            w_uv := g.Weight(e)
            if w_uv < 0.0 { return graphs.NegativeEdgeWeight(e) }
            d_v  := dmap.GetDistance(e.Target)
            if w_uv + d_u < d_v {
                dmap.SetDistance(e.Target, w_uv + d_u)
                if pmap != nil { pmap.SetVertex(e.Target, u) }
                if visitor.OnEdgeRelaxed != nil {
                    if err := visitor.OnEdgeRelaxed(e); err != nil { return err }
                }
                if d_v == d_inf {
                    queue.Push(e.Target)
                    if visitor.OnDiscoverVertex != nil {
                        if err := visitor.OnDiscoverVertex(e.Target); err != nil { return err }
                    }
                } else {
                    queue.InsertOrUpdate(e.Target)
                }
            } else {
                if visitor.OnEdgeNonRelaxed != nil {
                    if err := visitor.OnEdgeNonRelaxed(e); err != nil { return err }
                }
            }
        }
        if visitor.OnFinishVertex != nil {
            if err := visitor.OnFinishVertex(u); err != nil { return err }
        }
    }
    
    return nil
}
