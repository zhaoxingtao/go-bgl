package graphs

type VertexPropertiesCopier func (src VertexID, dst VertexID) error
type EdgePropertiesCopier   func (src EdgeIter, dst EdgeIter) error

// TransposeAdaptorForCopy still implement MutableGraph interface
// This can be used for transpose a graph
type TransposeAdaptorForCopy struct {
    MutableGraph
}

func (g TransposeAdaptorForCopy) AddEdge(u VertexID, v VertexID) (e EdgeIter, inserted bool) {
    return g.MutableGraph.AddEdge(v, u)
}

func CopyVertexEdgeListGraph(
    in VertexEdgeListGraph,
    out MutableGraph,
    orig2copy VertexMap,
    vcopier VertexPropertiesCopier,
    ecopier EdgePropertiesCopier) error {
        
    for v := in.VertexBegin(); v.IsValid(); v = in.VertexNext(v) {
        newV := out.AddVertex()
        orig2copy.SetVertex(v, newV)
        if vcopier != nil {
            if err := vcopier(v, newV); err != nil { return err }
        }
    }

    for e := in.EdgeBegin(); e.EdgeID.IsValid(); e = in.EdgeNext(e) {
        newE, _ := out.AddEdge(orig2copy.GetVertex(e.Source),
                               orig2copy.GetVertex(e.Target))
        if ecopier != nil {
            if err := ecopier(e, newE); err != nil { return err }
        }
    }
    return nil
}

func CopyVertexIncidenceGraph(
    in VertexIncidenceGraph,
    out MutableGraph,
    orig2copy VertexMap,
    vcopier VertexPropertiesCopier,
    ecopier EdgePropertiesCopier) error {
        
    for v := in.VertexBegin(); v.IsValid(); v = in.VertexNext(v) {
        newV := out.AddVertex()
        orig2copy.SetVertex(v, newV)
        if vcopier != nil {
            if err := vcopier(v, newV); err != nil { return err }
        }
    }
    
    for v := in.VertexBegin(); v.IsValid(); v = in.VertexNext(v) {
        for e := in.OutEdgeBegin(v); e.EdgeID.IsValid(); e = in.OutEdgeNext(e) {
            newE, _ := out.AddEdge(orig2copy.GetVertex(e.Source),
                                   orig2copy.GetVertex(e.Target))
            if ecopier != nil {
                if err := ecopier(e, newE); err != nil { return err }
            }
        }
    }
    return nil
}

func CopyVertexIncidenceUndirectedGraph(
    in VertexIncidenceGraph,
    out MutableGraph,
    orig2copy VertexMap,
    cmap ColorMap,
    vcopier VertexPropertiesCopier,
    ecopier EdgePropertiesCopier) error {
        
    for v := in.VertexBegin(); v.IsValid(); v = in.VertexNext(v) {
        newV := out.AddVertex()
        orig2copy.SetVertex(v, newV)
        cmap.SetColor(v, WHITE)
        if vcopier != nil {
            if err := vcopier(v, newV); err != nil { return err }
        }
    }
    
    for v := in.VertexBegin(); v.IsValid(); v = in.VertexNext(v) {
        for e := in.OutEdgeBegin(v); e.EdgeID.IsValid(); e = in.OutEdgeNext(e) {
            if cmap.GetColor(e.Target) == WHITE {
                newE, _ := out.AddEdge(orig2copy.GetVertex(e.Source),
                                       orig2copy.GetVertex(e.Target))
                if ecopier != nil {
                    if err := ecopier(e, newE); err != nil { return err }
                }
            }
        }
        cmap.SetColor(v, BLACK)
    }
    return nil
}

func CopyComponent(
    in IncidenceGraph,
    out MutableGraph,
    src VertexID,
    cmap ColorMap,
    orig2copy VertexMap,
    vcopier VertexPropertiesCopier,
    ecopier EdgePropertiesCopier) error {

    if src == NilVertex { return nil }
    
    newSrc := out.AddVertex()
    orig2copy.SetVertex(src, newSrc)
    if vcopier != nil {
        if err := vcopier(src, newSrc); err != nil { return err }
    }
    
    err := BreadthFirstVisit(in, cmap, nil, []VertexID{ src }, BfsVisitor{
        OnTreeEdge: func (e EdgeIter) error {
            // For a tree EdgeIter, the target vertex has not been copied yet.
            newV := out.AddVertex()
            orig2copy.SetVertex(e.Target, newV)
            if vcopier != nil {
                if err := vcopier(e.Target, newV); err != nil { return err }
            }
            newE, _ := out.AddEdge(orig2copy.GetVertex(e.Source), newV)
            if ecopier != nil {
                if err := ecopier(e, newE); err != nil { return err }
            }
            return nil
        },
        OnNonTreeEdge: func (e EdgeIter) error {
            // For a non-tree EdgeIter, the target vertex has already been copied.
            newE, _ := out.AddEdge(orig2copy.GetVertex(e.Source),
                                   orig2copy.GetVertex(e.Target))
            if ecopier != nil {
                if err := ecopier(e, newE); err != nil { return err }
            }
            return nil
        },
    })
    return err
}