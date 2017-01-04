package shortest_path

import "errors"
import "graphs"
import "math"

type BellmanVisitor struct {
    OnExamineEdge       func (e graphs.EdgeIter) error
    OnEdgeRelaxed       func (e graphs.EdgeIter) error
    OnEdgeNonRelaxed    func (e graphs.EdgeIter) error
    OnEdgeMinimized     func (e graphs.EdgeIter) error
    OnEdgeNotMinimized  func (e graphs.EdgeIter) error
}

func BellmanFordShortestPathsNoinit(
    g graphs.EdgeListWeightedGraph,
    N int, 
    pmap graphs.VertexMap,
    dmap graphs.DistanceMap,
    visitor BellmanVisitor) error {
        
    for k := 0; k < N; k++ {
        atLeastOneEdgeRelaxed := false
        for e := g.EdgeBegin(); e.EdgeID.IsValid(); e = g.EdgeNext(e) {
            if visitor.OnExamineEdge != nil {
                if err := visitor.OnExamineEdge(e); err != nil { return err }
            }
            d_u := dmap.GetDistance(e.Source)
            d_v := dmap.GetDistance(e.Target)
            w_uv := g.Weight(e)
            if w_uv + d_u < d_v {
                atLeastOneEdgeRelaxed = true
                dmap.SetDistance(e.Target, w_uv + d_u)
                if pmap != nil { pmap.SetVertex(e.Target, e.Source) }
                if visitor.OnEdgeRelaxed != nil {
                    if err := visitor.OnEdgeRelaxed(e); err != nil { return err }
                }
            } else {
                if visitor.OnEdgeNonRelaxed != nil {
                    if err := visitor.OnEdgeNonRelaxed(e); err != nil { return err }
                }
            }
        }
        if !atLeastOneEdgeRelaxed {
            break
        }
    }

    allMinimized := true
    for e := g.EdgeBegin(); e.EdgeID.IsValid(); e = g.EdgeNext(e) {
        d_u := dmap.GetDistance(e.Source)
        d_v := dmap.GetDistance(e.Target)
        w_uv := g.Weight(e)
        if w_uv + d_u < d_v {
            allMinimized = false
            if visitor.OnEdgeNotMinimized != nil {
                if err := visitor.OnEdgeNotMinimized(e); err != nil { return err }
            }
        } else {
            if visitor.OnEdgeMinimized != nil {
                if err := visitor.OnEdgeMinimized(e); err != nil { return err }
            }
        }
    }
    
    if !allMinimized {
        return errors.New("NotAllEdgesMinimized")
    }
    return nil;
}

func BellmanFordShortestPaths(
    g graphs.VertexEdgeListWeightedGraph,
    s graphs.VertexID,
    pmap graphs.VertexMap,
    dmap graphs.DistanceMap,
    visitor BellmanVisitor) error {
        
    N := 0
    for v := g.VertexBegin(); v.IsValid(); v = g.VertexNext(v) {
        dmap.SetDistance(v, math.Inf(1))
        pmap.SetVertex(v, v)
        N++
    }
    
    dmap.SetDistance(s, 0)
    
    return BellmanFordShortestPathsNoinit(g, N, pmap, dmap, visitor)
}
