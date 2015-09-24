package ipv

type DistanceSorter []IPVServer

func (a DistanceSorter) Len() int           { return len(a) }
func (a DistanceSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DistanceSorter) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

type LatencySorter []IPVServer

func (a LatencySorter) Len() int           { return len(a) }
func (a LatencySorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LatencySorter) Less(i, j int) bool { return a[i].Distance < a[j].Distance }