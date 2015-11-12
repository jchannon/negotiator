package negotiator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMediaRanges_parses_single(t *testing.T) {
	a := new(Accept)
	a.Header = "application/json"
	retVal := a.MediaRanges()

	assert.Equal(t, 1, len(retVal))
	assert.Equal(t, "application/json", retVal[0].Value)
	assert.Equal(t, 1.0, retVal[0].Weight)
}

func TestMediaRanges_defaults_quality_if_not_explicit(t *testing.T) {
	a := new(Accept)
	a.Header = "text/plain"
	mediaRanges := a.MediaRanges()
	assert.Equal(t, 1, len(mediaRanges))
	assert.Equal(t, 1.0, mediaRanges[0].Weight)
}

func TestMediaRanges_should_parse_quality(t *testing.T) {
	a := new(Accept)
	a.Header = "application/json;q=0.9"
	retVal := a.MediaRanges()

	assert.Equal(t, 1, len(retVal))
	assert.Equal(t, "application/json", retVal[0].Value)
	assert.Equal(t, 0.9, retVal[0].Weight)
}

func TestMediaRanges_should_parse_multi_qualities(t *testing.T) {
	a := new(Accept)
	a.Header = "application/xml;q=1, application/json;q=0.9"
	retVal := a.MediaRanges()

	assert.Equal(t, 2, len(retVal))

	assert.Equal(t, "application/xml", retVal[0].Value)
	assert.Equal(t, 1.0, retVal[0].Weight)

	assert.Equal(t, "application/json", retVal[1].Value)
	assert.Equal(t, 0.9, retVal[1].Weight)
}

func TestMediaRanges_orders_by_quality_decending(t *testing.T) {
	a := new(Accept)
	a.Header = "application/json;q=0.9, application/xml"
	retVal := a.MediaRanges()

	if len(retVal) != 2 {
		t.Errorf("Wanted 2, Got %d", len(retVal))
	}

	if retVal[0].Value != "application/xml" {
		t.Errorf("Wanted application/xml, Got %s", retVal[0].Value)
	}

	if retVal[0].Weight != 1.0 {
		t.Errorf("Wanted 1.0, Got %d", retVal[0].Weight)
	}

	if retVal[1].Value != "application/json" {
		t.Errorf("Wanted application/json, Got %s", retVal[1].Value)
	}

	if retVal[1].Weight != 0.9 {
		t.Errorf("Wanted 0.9, Got %d", retVal[1].Weight)
	}
}

func TestMediaRanges_should_ignore_invalid_quality(t *testing.T) {
	a := new(Accept)
	a.Header = "text/html;q=blah"
	ranges := a.MediaRanges()

	assert.Equal(t, 1, len(ranges))
	assert.Equal(t, "text/html", ranges[0].Value)
	assert.Equal(t, 1.0, ranges[0].Weight)
}

func TestMediaRanges_should_not_remove_accept_extension(t *testing.T) {
	a := new(Accept)
	a.Header = "text/html;q=0.5;a=1;b=2"
	ranges := a.MediaRanges()
	assert.Equal(t, 1, len(ranges))
	assert.Equal(t, "text/html;a=1;b=2", ranges[0].Value)
	assert.Equal(t, 0.5, ranges[0].Weight)
}
