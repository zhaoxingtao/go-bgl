package graphs

type EdgeStack []EdgeIter
func (q *EdgeStack) Clear()          { *q = (*q)[0:0] }
func (q *EdgeStack) IsEmpty() bool   { return len(*q) == 0 }
func (q *EdgeStack) Push(e EdgeIter) { *q = append(*q, e) }
func (q *EdgeStack) Pop() EdgeIter {
    e := (*q)[len(*q)-1];
    *q = (*q)[0:len(*q)-1];
    return e
}

type DfsVisitor struct {
    OnInitializeVertex   func (v VertexID) error
    OnStartVertex        func (v VertexID) error
    OnExamineEdge        func (e EdgeIter) error
    OnTreeEdge           func (e EdgeIter) error
    OnDiscoverVertex     func (v VertexID) error
    OnBackEdge           func (e EdgeIter) error
    OnForwardOrCrossEdge func (e EdgeIter) error
    OnFinishNonTreeEdge  func (e EdgeIter) error
    OnFinishVertex       func (v VertexID) error
    OnFinishTreeEdge     func (e EdgeIter) error
}

func DepthFirstSearch(
    g VertexIncidenceGraph,
    cmap ColorMap,
    pmap VertexMap,
    source VertexID,
    visitor DfsVisitor) error {

    for v := g.VertexBegin(); v.IsValid(); v = g.VertexNext(v) {
        cmap.SetColor(v, WHITE)
        if pmap != nil { pmap.SetVertex(v, v) }
        if visitor.OnInitializeVertex != nil {
            if err := visitor.OnInitializeVertex(v); err != nil { return err }
        }
    }

    if source.IsValid() {
        if visitor.OnStartVertex != nil {
            if err := visitor.OnStartVertex(source); err != nil { return err }
        }
        if err := DepthFirstVisit(g, cmap, pmap, source, visitor); err != nil { return err }
    }

    for v := g.VertexBegin(); v.IsValid(); v = g.VertexNext(v) {
        if cmap.GetColor(v) == WHITE {
            if visitor.OnStartVertex != nil {
                if err := visitor.OnStartVertex(v); err != nil { return err }
            }
            if err := DepthFirstVisit(g, cmap, pmap, v, visitor); err != nil { return err }
        }
    }
    
    return nil
}

func DepthFirstVisit(
    g IncidenceGraph,
    cmap ColorMap,
    pmap VertexMap,
    u VertexID,
    visitor DfsVisitor) error {

    stack := EdgeStack{}

    cmap.SetColor(u, GRAY)
    if visitor.OnDiscoverVertex != nil {
        if err := visitor.OnDiscoverVertex(u); err != nil { return err }
    }
    
    stack.Push(g.OutEdgeBegin(u))
    
    for !stack.IsEmpty() {
        e := stack.Pop()
        for e.EdgeID.IsValid() {
            if visitor.OnExamineEdge != nil {
                if err := visitor.OnExamineEdge(e); err != nil { return err }
            }
            u = e.Target
            color := cmap.GetColor(u)
            if color == WHITE {
                if visitor.OnTreeEdge != nil {
                    if err := visitor.OnTreeEdge(e); err != nil { return err }
                }
                cmap.SetColor(u, GRAY)
                if pmap != nil { pmap.SetVertex(u, e.Source) }
                stack.Push(e)
                e = g.OutEdgeBegin(u)
                if visitor.OnDiscoverVertex != nil {
                    if err := visitor.OnDiscoverVertex(u); err != nil { return err }
                }
            } else {
                if color == GRAY {
                    if visitor.OnBackEdge != nil {
                        if err := visitor.OnBackEdge(e); err != nil { return err }
                    }
                } else {
                    if visitor.OnForwardOrCrossEdge != nil {
                        if err := visitor.OnForwardOrCrossEdge(e); err != nil { return err }
                    }
                }
                if visitor.OnFinishNonTreeEdge != nil {
                    if err := visitor.OnFinishNonTreeEdge(e); err != nil { return err }
                }
                e = g.OutEdgeNext(e)
            }
        }  // end loop depth
        u = e.Source
        cmap.SetColor(u, BLACK)
        if visitor.OnFinishVertex != nil {
            if err := visitor.OnFinishVertex(u); err != nil { return err }
        }
        if !stack.IsEmpty() {
            e = stack.Pop()
            if visitor.OnFinishTreeEdge != nil {
                if err := visitor.OnFinishTreeEdge(e); err != nil { return err }
            }
            stack.Push(g.OutEdgeNext(e))
        }
    }  // end loop stack
    
    return nil
}
