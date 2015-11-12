package negotiator

import (
	"sort"
	"strconv"
	"strings"
)

// Accept is an http accept
type Accept struct {
	Header string
}

// WeightedValue is a value and associate weight between 0.0 and 1.0
type WeightedValue struct {
	Value  string
	Weight float64
}

// ByWeight implements sort.Interface for []WeightedValue based
//on the Weight field
type ByWeight []WeightedValue

func (a ByWeight) Len() int           { return len(a) }
func (a ByWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByWeight) Less(i, j int) bool { return a[i].Weight > a[j].Weight }

// MediaRanges returns prioritized media ranges
func (accept *Accept) MediaRanges() []WeightedValue {
	var retVals []WeightedValue
	mrs := strings.Split(accept.Header, ",")

	for _, mr := range mrs {
		mrAndAcceptParam := strings.Split(mr, ";")
		//if no quality assigned then give 1.0
		if len(mrAndAcceptParam) == 1 {
			wv := new(WeightedValue)
			wv.Value = strings.TrimSpace(mrAndAcceptParam[0])
			wv.Weight = 1.0
			retVals = append(retVals, *wv)
			continue
		}

		wv := new(WeightedValue)
		wv.Value = strings.TrimSpace(mrAndAcceptParam[0])

		var weight float64
		var err error
		for index := 1; index < len(mrAndAcceptParam); index++ {
			if strings.Contains(mrAndAcceptParam[index], "q=") {
				weight, err = strconv.ParseFloat(strings.SplitAfter(mrAndAcceptParam[index], "q=")[1], 64)
				if err != nil {
					weight = 1.0
				}
			} else {
				wv.Value = strings.Join([]string{wv.Value, mrAndAcceptParam[index]}, ";")
			}

		}

		wv.Weight = weight
		retVals = append(retVals, *wv)
	}

	//If no Accept header field is present, then it is assumed that the client
	//accepts all media types. If an Accept header field is present, and if the
	//server cannot send a response which is acceptable according to the combined
	//Accept field value, then the server SHOULD send a 406 (not acceptable)
	//response.

	sort.Sort(ByWeight(retVals))

	return retVals
}
