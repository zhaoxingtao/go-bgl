package graphs

import "strconv"

type Color byte
const (
    WHITE Color = iota
    GRAY
    GREEN
    RED
    BLACK
)

func (c Color) String() string {
    switch c {
    case WHITE: return "WHITE"
    case GRAY:  return "GRAY"
    case GREEN: return "GREEN"
    case RED:   return "RED"
    case BLACK: return "BLACK"
    default:    return strconv.Itoa(int(c))
    }
}

// The following maps should be initialized with length = numVertices before using.

type ColorMap []Color
func (cmap ColorMap) GetColor(v VertexID) Color    { return cmap[int(v)] }
func (cmap ColorMap) SetColor(v VertexID, c Color) { cmap[int(v)] = c    }

type IndexMap []int
func (imap IndexMap) GetIndex(v VertexID) int    { return imap[int(v)] }
func (imap IndexMap) SetIndex(v VertexID, i int) { imap[int(v)] = i    }

type DistanceMap []float64
func (dmap DistanceMap) GetDistance(v VertexID) float64    { return dmap[int(v)] }
func (dmap DistanceMap) SetDistance(v VertexID, d float64) { dmap[int(v)] = d    }

type VertexMap []VertexID
func (vmap VertexMap) GetVertex(v VertexID) VertexID    { return vmap[int(v)] }
func (vmap VertexMap) SetVertex(v VertexID, u VertexID) { vmap[int(v)] = u    }

type EdgeColorMap map[EdgeIter]Color
func (emap EdgeColorMap) GetColor(e EdgeIter) Color { return emap[e] }
func (emap EdgeColorMap) SetColor(e EdgeIter, c Color) {
    if c == WHITE {
        delete(emap, e)
    } else {
        emap[e] = c
    }
}
