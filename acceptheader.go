package negotiator

import (
	"sort"
	"strconv"
	"strings"
)

const (
	// ParameteredMediaRangeWeight is the default weight of a media range with an
	// accept-param
	ParameteredMediaRangeWeight float64 = 1.0 //e.g text/html;level=1
	// TypeSubtypeMediaRangeWeight is the default weight of a media range with
	// type and subtype defined
	TypeSubtypeMediaRangeWeight float64 = 0.9 //e.g text/html
	// TypeStarMediaRangeWeight is the default weight of a media range with a type
	// defined but * for subtype
	TypeStarMediaRangeWeight float64 = 0.8 //e.g text/*
	// StarStarMediaRangeWeight is the default weight of a media range with any
	// type or any subtype defined
	StarStarMediaRangeWeight float64 = 0.7 //e.g */*
)

// Accept is an http accept
type accept struct {
	Header string
}

// MediaRanges returns prioritized media ranges
func (accept *accept) ParseMediaRanges() []weightedValue {
	var retVals []weightedValue
	mrs := strings.Split(accept.Header, ",")

	for _, mr := range mrs {
		mrAndAcceptParam := strings.Split(mr, ";")
		//if no accept-param
		if len(mrAndAcceptParam) == 1 {
			retVals = append(retVals, handleMediaRangeNoAcceptParams(mrAndAcceptParam[0]))
			continue
		}

		retVals = append(retVals, handleMediaRangeWithAcceptParams(mrAndAcceptParam[0], mrAndAcceptParam[1:]))
	}

	//If no Accept header field is present, then it is assumed that the client
	//accepts all media types. If an Accept header field is present, and if the
	//server cannot send a response which is acceptable according to the combined
	//Accept field value, then the server SHOULD send a 406 (not acceptable)
	//response.
	sort.Sort(byWeight(retVals))

	return retVals
}

func handleMediaRangeWithAcceptParams(mediaRange string, acceptParams []string) weightedValue {
	wv := new(weightedValue)
	wv.Value = strings.TrimSpace(mediaRange)
	wv.Weight = ParameteredMediaRangeWeight

	for index := 0; index < len(acceptParams); index++ {
		ap := strings.ToLower(acceptParams[index])
		if isQualityAcceptParam(ap) {
			wv.Weight = parseQuality(ap)
		} else {
			wv.Value = strings.Join([]string{wv.Value, acceptParams[index]}, ";")
		}
	}
	return *wv
}

func isQualityAcceptParam(acceptParam string) bool {
	return strings.Contains(acceptParam, "q=")
}

func parseQuality(acceptParam string) float64 {
	weight, err := strconv.ParseFloat(strings.SplitAfter(acceptParam, "q=")[1], 64)
	if err != nil {
		weight = 1.0
	}
	return weight
}

func handleMediaRangeNoAcceptParams(mediaRange string) weightedValue {
	wv := new(weightedValue)
	wv.Value = strings.TrimSpace(mediaRange)
	wv.Weight = 0.0

	typeSubtype := strings.Split(wv.Value, "/")
	if len(typeSubtype) == 2 {
		switch {
		//a type of * with a non-star subtype is invalid, so if the type is
		//star the assume that the subtype is too
		case typeSubtype[0] == "*": //&& typeSubtype[1] == "*":
			wv.Weight = StarStarMediaRangeWeight
			break
		case typeSubtype[1] == "*":
			wv.Weight = TypeStarMediaRangeWeight
			break
		case typeSubtype[1] != "*":
			wv.Weight = TypeSubtypeMediaRangeWeight
			break
		}
	} //else invalid media range the weight remains 0.0

	return *wv
}
