package connected_component

import "graphs"

func ConnectedComponents(
    g graphs.VertexIncidenceGraph,
    color_map graphs.ColorMap,
    component_map graphs.IndexMap) int {
    
    count := -1  // start counting components at zero
    graphs.DepthFirstSearch(g, color_map, nil, graphs.NilVertex, graphs.DfsVisitor{
        OnStartVertex:    func (v graphs.VertexID) error { count++; return nil },
        OnDiscoverVertex: func (v graphs.VertexID) error { component_map.SetIndex(v, count); return nil },
    })
    return count
}

func StrongComponents(
    g graphs.VertexIncidenceGraph,
    color_map graphs.ColorMap,
    component_map graphs.IndexMap,
    root_map graphs.VertexMap,
    discover_time_map graphs.IndexMap) int {
    
    stack := []graphs.VertexID{}
    count := -1
    dfs_time := 0
    
    graphs.DepthFirstSearch(g, color_map, nil, graphs.NilVertex, graphs.DfsVisitor{
        OnDiscoverVertex: func (v graphs.VertexID) error {
            root_map.SetVertex(v, v)
            component_map.SetIndex(v, -1)
            discover_time_map.SetIndex(v, dfs_time)
            dfs_time++
            stack = append(stack, v)
            return nil
        },
        OnFinishVertex: func (v graphs.VertexID) error {
            for e := g.OutEdgeBegin(v); e.EdgeID.IsValid(); e = g.OutEdgeNext(e) {
                if component_map.GetIndex(e.Target) < 0 {
                    root_v := root_map.GetVertex(v)
                    root_w := root_map.GetVertex(e.Target)
                    if discover_time_map.GetIndex(root_v) >= discover_time_map.GetIndex(root_w) {
                        root_map.SetVertex(v, root_w)
                    }
                }
            }
            if root_map.GetVertex(v) == v {
                i := len(stack) - 1
                for ; i >= 0; i-- {
                    w := stack[i]
                    stack[i] = graphs.NilVertex
                    component_map.SetIndex(w, count)
                    if w == v { break }
                }
                stack = stack[0:i]
                count++
            }
            return nil
        },
    })
    return count
}

func BuildComponentList(
    g graphs.VertexList,
    component_map graphs.IndexMap,
    component_counts int) [][]graphs.VertexID {

    component_list := make([][]graphs.VertexID, component_counts)
    for v := g.VertexBegin(); v.IsValid(); v = g.VertexNext(v) {
        cid := component_map.GetIndex(v)
        component_list[cid] = append(component_list[cid], v)
    }
    return component_list
}
