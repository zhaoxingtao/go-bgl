package graphs

type VertexQueue []VertexID
func (q *VertexQueue) Clear()          { *q = (*q)[0:0] }
func (q *VertexQueue) IsEmpty() bool   { return len(*q) == 0 }
func (q *VertexQueue) Push(v VertexID) { *q = append(*q, v) }
func (q *VertexQueue) Pop() VertexID   {
    v := (*q)[0];
    *q = (*q)[1:];
    return v
}

type BfsVisitor struct {
    OnInitializeVertex func (v VertexID) error
    OnExamineVertex    func (v VertexID) error
    OnExamineEdge      func (e EdgeIter) error
    OnTreeEdge         func (e EdgeIter) error
    OnDiscoverVertex   func (v VertexID) error
    OnNonTreeEdge      func (e EdgeIter) error
    OnGrayTargetEdge   func (e EdgeIter) error
    OnBlackTargeEdge   func (e EdgeIter) error
    OnFinishVertex     func (v VertexID) error
}

func BreadthFirstSearch(
    g VertexIncidenceGraph,
    cmap ColorMap,
    pmap VertexMap,
    sources []VertexID,
    visitor BfsVisitor) error {

    for v := g.VertexBegin(); v.IsValid(); v = g.VertexNext(v) {
        if visitor.OnInitializeVertex != nil {
            if err := visitor.OnInitializeVertex(v); err != nil { return err }
        }
        cmap.SetColor(v, WHITE)
        if pmap != nil { pmap.SetVertex(v, v) }
    }

    if err := BreadthFirstVisit(g, cmap, pmap, sources, visitor); err != nil { return err }
    return nil
}

func BreadthFirstVisit(
    g IncidenceGraph,
    cmap ColorMap,
    pmap VertexMap,
    sources []VertexID,
    visitor BfsVisitor) error {

    queue := VertexQueue{}

    for _, v := range sources {
        if visitor.OnDiscoverVertex != nil {
            if err := visitor.OnDiscoverVertex(v); err != nil { return err }
        }
        queue.Push(v)
    }
    
    for !queue.IsEmpty() {
        v := queue.Pop()
        if visitor.OnExamineVertex != nil {
            if err := visitor.OnExamineVertex(v); err != nil { return err }
        }
        for e := g.OutEdgeBegin(v); e.EdgeID.IsValid(); e = g.OutEdgeNext(e) {
            if visitor.OnExamineEdge != nil {
                if err := visitor.OnExamineEdge(e); err != nil { return err }
            }
            color := cmap.GetColor(e.Target)
            if color == WHITE {
                if visitor.OnTreeEdge != nil {
                    if err := visitor.OnTreeEdge(e); err != nil { return err }
                }
                cmap.SetColor(e.Target, GRAY)
                if pmap != nil { pmap.SetVertex(e.Target, v) }
                if visitor.OnDiscoverVertex != nil {
                    if err := visitor.OnDiscoverVertex(e.Target); err != nil { return err }
                }
                queue.Push(e.Target)
            } else {
                if visitor.OnNonTreeEdge != nil {
                    if err := visitor.OnNonTreeEdge(e); err != nil { return err }
                }
                if color == GRAY {
                    if visitor.OnGrayTargetEdge != nil {
                        if err := visitor.OnGrayTargetEdge(e); err != nil { return err }
                    }
                } else {
                    if visitor.OnBlackTargeEdge != nil {
                        if err := visitor.OnBlackTargeEdge(e); err != nil { return err }
                    }
                }
            }
        }
        cmap.SetColor(v, BLACK)
        if visitor.OnFinishVertex != nil {
            if err := visitor.OnFinishVertex(v); err != nil { return err }
        }
    }
    
    return nil
}
