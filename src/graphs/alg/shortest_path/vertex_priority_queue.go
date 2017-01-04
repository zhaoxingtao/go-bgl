package shortest_path

import "graphs"

type vertexPriorityItem struct {
    v graphs.VertexID
    d float64
}

type VertexPriorityQueue struct {
    queue []vertexPriorityItem
    imap  graphs.IndexMap
    dmap  graphs.DistanceMap
}

func NewVertexPriorityQueue(imap graphs.IndexMap, dmap graphs.DistanceMap) VertexPriorityQueue {
    return VertexPriorityQueue{ imap: imap, dmap: dmap }
}

// Len return the depth of the heap.
func (pq *VertexPriorityQueue) Len() int {
	return len(pq.queue)
}

// Push pushes the element x onto the heap. The complexity is
// O(log(n)) where n = pq.Len().
//
func (pq *VertexPriorityQueue) Push(v graphs.VertexID) {
	pq.push(v)
	pq.up(pq.Len()-1)
}

// Pop removes the minimum element (according to Less) from the heap
// and returns it. The complexity is O(log(n)) where n = pq.Len().
// It is equivalent to pq.Remove(0).
//
func (pq *VertexPriorityQueue) Pop() graphs.VertexID {
	n := pq.Len() - 1
	pq.swap(0, n)
	pq.down(0, n)
    return pq.pop()
}

// InsertOrUpdate inserts an element when the element is not in the queue,
// or re-establishes the queue ordering if the element is already in the queue.
//
func (pq *VertexPriorityQueue) InsertOrUpdate(v graphs.VertexID) {
    i := pq.imap.GetIndex(v) - 1
    if i < 0 {
        pq.Push(v)
        return
    }
    pq.queue[i].d = pq.dmap.GetDistance(v)
	pq.down(i, pq.Len())
	pq.up(i)
}

// private methods
func (pq *VertexPriorityQueue) less(i, j int) bool { 
	return pq.queue[i].d < pq.queue[j].d
}

func (pq *VertexPriorityQueue) swap(i, j int) {
    item_i, item_j := pq.queue[i], pq.queue[j]
    pq.imap.SetIndex(item_i.v, j + 1)
    pq.imap.SetIndex(item_j.v, i + 1)
    pq.queue[i], pq.queue[j] = item_j, item_i
}

func (pq *VertexPriorityQueue) push(v graphs.VertexID) {
    n  := pq.Len()
    d  := pq.dmap.GetDistance(v)
    pq.queue = append(pq.queue, vertexPriorityItem{ v: v, d: d })
    pq.imap.SetIndex(v, n + 1)
}

func (pq *VertexPriorityQueue) pop() graphs.VertexID {
	n := pq.Len() - 1
    v := pq.queue[n].v
    pq.queue = pq.queue[0:n]
    pq.imap.SetIndex(v, 0)
    return v
}

func (pq *VertexPriorityQueue) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !pq.less(j, i) {
			break
		}
		pq.swap(i, j)
		j = i
	}
}

func (pq *VertexPriorityQueue) down(i, n int) {
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && !pq.less(j1, j2) {
			j = j2 // = 2*i + 2  // right child
		}
		if !pq.less(j, i) {
			break
		}
		pq.swap(i, j)
		i = j
	}
}
