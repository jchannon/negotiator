package negotiator

// WeightedValue is a value and associate weight between 0.0 and 1.0
type weightedValue struct {
	Value  string
	Weight float64
}

// ByWeight implements sort.Interface for []WeightedValue based
//on the Weight field. The data will be returned sorted decending
type byWeight []weightedValue

func (a byWeight) Len() int           { return len(a) }
func (a byWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byWeight) Less(i, j int) bool { return a[i].Weight > a[j].Weight }
