package negotiator

// WeightedValue is a value and associate weight between 0.0 and 1.0
type WeightedValue struct {
	Value  string
	Weight float64
}

// ByWeight implements sort.Interface for []WeightedValue based
//on the Weight field. The data will be returned sorted decending
type ByWeight []WeightedValue

func (a ByWeight) Len() int           { return len(a) }
func (a ByWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByWeight) Less(i, j int) bool { return a[i].Weight > a[j].Weight }
