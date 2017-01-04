package shortest_path

import "graphs"
import "math"

type AStarVisitor struct {
    OnInitializeVertex func (v graphs.VertexID) error
    OnExamineVertex    func (v graphs.VertexID) error
    OnExamineEdge      func (e graphs.EdgeIter) error
    OnDiscoverVertex   func (v graphs.VertexID) error
    OnEdgeRelaxed      func (e graphs.EdgeIter) error
    OnEdgeNonRelaxed   func (e graphs.EdgeIter) error
    OnBlackTargetEdge  func (e graphs.EdgeIter) error
    OnFinishVertex     func (v graphs.VertexID) error
}

func AStarSearch(
    g graphs.VertexIncidenceWeightedGraph,
    cmap graphs.ColorMap,
    pmap graphs.VertexMap,
    dmap graphs.DistanceMap,
    cost graphs.DistanceMap,
    imap graphs.IndexMap,
    source graphs.VertexID,
    heuristic func(v graphs.VertexID) float64,
    visitor AStarVisitor) error {

    d_inf := math.Inf(1)
    for v := g.VertexBegin(); v.IsValid(); v = g.VertexNext(v) {
        if visitor.OnInitializeVertex != nil {
            if err := visitor.OnInitializeVertex(v); err != nil { return err }
        }
        cmap.SetColor(v, graphs.WHITE)
        if pmap != nil { pmap.SetVertex(v, v) }
        dmap.SetDistance(v, d_inf)
        cost.SetDistance(v, d_inf)
    }

    return AStarSearchNoInit(g, cmap, pmap, dmap, cost, imap, source, heuristic, visitor)
}

func AStarSearchNoInit(
    g graphs.IncidenceWeightedGraph,
    cmap graphs.ColorMap,
    pmap graphs.VertexMap,
    dmap graphs.DistanceMap,
    cost graphs.DistanceMap,
    imap graphs.IndexMap,
    source graphs.VertexID,
    heuristic func (v graphs.VertexID) float64,
    visitor AStarVisitor) error {

    queue := NewVertexPriorityQueue(imap, cost)

    cmap.SetColor(source, graphs.GRAY)
    dmap.SetDistance(source, 0.0)
    cost.SetDistance(source, heuristic(source))
    queue.Push(source)
    if visitor.OnDiscoverVertex != nil {
        if err := visitor.OnDiscoverVertex(source); err != nil { return err }
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
            d_v2 := w_uv + d_u
            if d_v2 < d_v {
                if visitor.OnEdgeRelaxed != nil {
                    if err := visitor.OnEdgeRelaxed(e); err != nil { return err }
                }
                dmap.SetDistance(e.Target, d_v2)
                cost.SetDistance(e.Target, d_v2 + heuristic(e.Target))
                if pmap != nil { pmap.SetVertex(e.Target, u) }
                color := cmap.GetColor(e.Target)
                if color == graphs.WHITE {
                    cmap.SetColor(e.Target, graphs.GRAY)
                    queue.Push(e.Target)
                    if visitor.OnDiscoverVertex != nil {
                        if err := visitor.OnDiscoverVertex(e.Target); err != nil { return err }
                    }
                } else if color == graphs.BLACK {
                    cmap.SetColor(e.Target, graphs.GRAY)
                    queue.Push(e.Target)
                    if visitor.OnBlackTargetEdge != nil {
                        if err := visitor.OnBlackTargetEdge(e); err != nil { return err }
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
        cmap.SetColor(u, graphs.BLACK)
        if visitor.OnFinishVertex != nil {
            if err := visitor.OnFinishVertex(u); err != nil { return err }
        }
    }
    
    return nil
}

func AStarSearchTree(
    g graphs.VertexIncidenceWeightedGraph,
    pmap graphs.VertexMap,
    dmap graphs.DistanceMap,
    cost graphs.DistanceMap,
    imap graphs.IndexMap,
    source graphs.VertexID,
    heuristic func(v graphs.VertexID) float64,
    visitor AStarVisitor) error {

    d_inf := math.Inf(1)
    for v := g.VertexBegin(); v.IsValid(); v = g.VertexNext(v) {
        if visitor.OnInitializeVertex != nil {
            if err := visitor.OnInitializeVertex(v); err != nil { return err }
        }
        if pmap != nil { pmap.SetVertex(v, v) }
        dmap.SetDistance(v, d_inf)
        cost.SetDistance(v, d_inf)
        imap.SetIndex(v, 0);
    }

    return AStarSearchNoInitTree(g, pmap, dmap, cost, imap, source, heuristic, visitor)
}

func AStarSearchNoInitTree(
    g graphs.IncidenceWeightedGraph,
    pmap graphs.VertexMap,
    dmap graphs.DistanceMap,
    cost graphs.DistanceMap,
    imap graphs.IndexMap,
    source graphs.VertexID,
    heuristic func(v graphs.VertexID) float64,
    visitor AStarVisitor) error {

    queue := NewVertexPriorityQueue(imap, cost)

    dmap.SetDistance(source, 0.0)
    cost.SetDistance(source, heuristic(source))
    queue.Push(source)
    if visitor.OnDiscoverVertex != nil {
        if err := visitor.OnDiscoverVertex(source); err != nil { return err }
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
            d_v2 := w_uv + d_u
            if d_v2 < d_v {
                if visitor.OnEdgeRelaxed != nil {
                    if err := visitor.OnEdgeRelaxed(e); err != nil { return err }
                }
                dmap.SetDistance(e.Target, d_v2)
                cost.SetDistance(e.Target, d_v2 + heuristic(e.Target))
                if pmap != nil { pmap.SetVertex(e.Target, u) }
                if visitor.OnDiscoverVertex != nil {
                    if err := visitor.OnDiscoverVertex(e.Target); err != nil { return err }
                }
                queue.InsertOrUpdate(e.Target)
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
