package ipv
import (
	"sort"
)

var Distance = func(c1, c2 *IPVServer) bool {
return c1.Distance < c2.Distance
}
var Latency = func(c1, c2 *IPVServer) bool {
return c1.Latency < c2.Latency
}
var Capacity = func(c1, c2 *IPVServer) bool {
return c1.Properties.Capacity < c2.Properties.Capacity
}


type lessFunc func(p1, p2 *IPVServer) bool

type multiSorter struct {
	servers []IPVServer
	less    []lessFunc
}

func (ms *multiSorter) Sort(servers []IPVServer) {
	ms.servers = servers
	sort.Sort(ms)
}

func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

func (ms *multiSorter) Len() int {
	return len(ms.servers)
}

func (ms *multiSorter) Swap(i, j int) {
	ms.servers[i], ms.servers[j] = ms.servers[j], ms.servers[i]
}

func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.servers[i], &ms.servers[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	return ms.less[k](p, q)
}

