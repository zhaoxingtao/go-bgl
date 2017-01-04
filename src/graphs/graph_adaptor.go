package graphs

type VertexIncidenceToEdgeListAdaptor struct {
    VertexIncidenceGraph
}

func (g VertexIncidenceToEdgeListAdaptor) EdgeBegin() EdgeIter {
    v := g.VertexBegin()
    if v.IsNil() { return NilEdgeIter() }
    e := g.OutEdgeBegin(v)
    if e.Target.IsValid() { return e }
    return g.EdgeNext(e)
}

func (g VertexIncidenceToEdgeListAdaptor) EdgeNext(e EdgeIter) EdgeIter {
    if e.Target.IsValid() {
        e = g.OutEdgeNext(e)
    }
    for e.Source.IsValid() && e.Target.IsNil() {
        v := g.VertexNext(e.Source);
        if v.IsNil() { return NilEdgeIter() }
        e = g.OutEdgeBegin(v)
    }
    return e
}
