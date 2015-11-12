package negotiator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMediaRanges_parses_single(t *testing.T) {
	a := new(Accept)
	a.Header = "application/json"
	mr := a.ParseMediaRanges()

	assert.Equal(t, 1, len(mr))
	assert.Equal(t, "application/json", mr[0].Value)
	assert.Equal(t, 1.0, mr[0].Weight)
}

func TestParseMediaRanges_defaults_quality_if_not_explicit(t *testing.T) {
	a := new(Accept)
	a.Header = "text/plain"
	mr := a.ParseMediaRanges()
	assert.Equal(t, 1, len(mr))
	assert.Equal(t, 1.0, mr[0].Weight)
}

func TestParseMediaRanges_should_parse_quality(t *testing.T) {
	a := new(Accept)
	a.Header = "application/json;q=0.9"
	mr := a.ParseMediaRanges()

	assert.Equal(t, 1, len(mr))
	assert.Equal(t, "application/json", mr[0].Value)
	assert.Equal(t, 0.9, mr[0].Weight)
}

func TestParseMediaRanges_should_parse_multi_qualities(t *testing.T) {
	a := new(Accept)
	a.Header = "application/xml;q=1, application/json;q=0.9"
	mr := a.ParseMediaRanges()

	assert.Equal(t, 2, len(mr))

	assert.Equal(t, "application/xml", mr[0].Value)
	assert.Equal(t, 1.0, mr[0].Weight)

	assert.Equal(t, "application/json", mr[1].Value)
	assert.Equal(t, 0.9, mr[1].Weight)
}

func TestParseMediaRanges_reorders_by_quality_decending(t *testing.T) {
	a := new(Accept)
	a.Header = "application/json;q=0.9, application/xml"
	mr := a.ParseMediaRanges()

	assert.Equal(t, 2, len(mr))

	assert.Equal(t, "application/xml", mr[0].Value)
	assert.Equal(t, 1.0, mr[0].Weight)

	assert.Equal(t, "application/json", mr[1].Value)
	assert.Equal(t, 0.9, mr[1].Weight)
}

func TestMediaRanges_should_ignore_invalid_quality(t *testing.T) {
	a := new(Accept)
	a.Header = "text/html;q=blah"
	mr := a.ParseMediaRanges()

	assert.Equal(t, 1, len(mr))
	assert.Equal(t, "text/html", mr[0].Value)
	assert.Equal(t, 1.0, mr[0].Weight)
}

func TestMediaRanges_should_not_remove_accept_extension(t *testing.T) {
	a := new(Accept)
	a.Header = "text/html;q=0.5;a=1;b=2"
	mr := a.ParseMediaRanges()
	assert.Equal(t, 1, len(mr))
	assert.Equal(t, "text/html;a=1;b=2", mr[0].Value)
	assert.Equal(t, 0.5, mr[0].Weight)
}
