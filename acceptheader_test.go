package negotiator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMediaRanges_parses_single(t *testing.T) {
	a := new(accept)
	a.Header = "application/json"
	mr := a.ParseMediaRanges()

	assert.Equal(t, 1, len(mr))
	assert.Equal(t, "application/json", mr[0].Value)
	assert.Equal(t, TypeSubtypeMediaRangeWeight, mr[0].Weight)
}

func TestParseMediaRanges_preserves_case_of_mediaRange(t *testing.T) {
	a := new(accept)
	a.Header = "application/CEA"
	mr := a.ParseMediaRanges()

	assert.Equal(t, 1, len(mr))
	assert.Equal(t, "application/CEA", mr[0].Value)
}

func TestParseMediaRanges_defaults_quality_if_not_explicit(t *testing.T) {
	a := new(accept)
	a.Header = "text/plain"
	mr := a.ParseMediaRanges()
	assert.Equal(t, 1, len(mr))
	assert.Equal(t, TypeSubtypeMediaRangeWeight, mr[0].Weight)
}

func TestParseMediaRanges_should_parse_quality(t *testing.T) {
	a := new(accept)
	a.Header = "application/json;q=0.9"
	mr := a.ParseMediaRanges()

	assert.Equal(t, 1, len(mr))
	assert.Equal(t, "application/json", mr[0].Value)
	assert.Equal(t, 0.9, mr[0].Weight)
}

func TestParseMediaRanges_should_parse_multi_qualities(t *testing.T) {
	a := new(accept)
	a.Header = "application/xml;q=1, application/json;q=0.9"
	mr := a.ParseMediaRanges()

	assert.Equal(t, 2, len(mr))

	assert.Equal(t, "application/xml", mr[0].Value)
	assert.Equal(t, 1.0, mr[0].Weight)

	assert.Equal(t, "application/json", mr[1].Value)
	assert.Equal(t, 0.9, mr[1].Weight)
}

func TestParseMediaRanges_reorders_by_quality_decending(t *testing.T) {
	a := new(accept)
	a.Header = "application/json;q=0.8, application/xml"
	mr := a.ParseMediaRanges()

	assert.Equal(t, 2, len(mr))

	assert.Equal(t, "application/xml", mr[0].Value)
	assert.Equal(t, TypeSubtypeMediaRangeWeight, mr[0].Weight)

	assert.Equal(t, "application/json", mr[1].Value)
	assert.Equal(t, 0.8, mr[1].Weight)
}

func TestMediaRanges_should_ignore_invalid_quality(t *testing.T) {
	a := new(accept)
	a.Header = "text/html;q=blah"
	mr := a.ParseMediaRanges()

	assert.Equal(t, 1, len(mr))
	assert.Equal(t, "text/html", mr[0].Value)
	assert.Equal(t, ParameteredMediaRangeWeight, mr[0].Weight)
}

func TestMediaRanges_should_not_remove_accept_extension(t *testing.T) {
	a := new(accept)
	a.Header = "text/html;q=0.5;a=1;b=2"
	mr := a.ParseMediaRanges()
	assert.Equal(t, 1, len(mr))
	assert.Equal(t, "text/html;a=1;b=2", mr[0].Value)
	assert.Equal(t, 0.5, mr[0].Weight)
}

func TestMediaRanges_should_handle_precedence(t *testing.T) {
	a := new(accept)
	a.Header = "text/*, text/html, text/html;level=1, */*"
	mr := a.ParseMediaRanges()
	assert.Equal(t, "text/html;level=1", mr[0].Value)
	assert.Equal(t, "text/html", mr[1].Value)
	assert.Equal(t, "text/*", mr[2].Value)
	assert.Equal(t, "*/*", mr[3].Value)
}

func TestMediaRanges_should_handle_precedence2(t *testing.T) {
	a := new(accept)
	a.Header = "text/*;q=0.3, text/html;q=0.7, text/html;level=1, text/html;level=2;q=0.4, */*;q=0.5"
	mr := a.ParseMediaRanges()

	assert.Equal(t, 5, len(mr))

	assert.Equal(t, "text/html;level=1", mr[0].Value)
	assert.Equal(t, 1.0, mr[0].Weight)

	assert.Equal(t, "text/html", mr[1].Value)
	assert.Equal(t, 0.7, mr[1].Weight)

	assert.Equal(t, "*/*", mr[2].Value)
	assert.Equal(t, 0.5, mr[2].Weight)

	assert.Equal(t, "text/html;level=2", mr[3].Value)
	assert.Equal(t, 0.4, mr[3].Weight)

	assert.Equal(t, "text/*", mr[4].Value)
	assert.Equal(t, 0.3, mr[4].Weight)
}

func TestMediaRanges_should_handle_precedence3(t *testing.T) {
	a := new(accept)
	// from http://tools.ietf.org/html/rfc7231#section-5.3.2
	a.Header = "text/*, text/plain, text/plain;format=flowed, */*"
	mr := a.ParseMediaRanges()

	assert.Equal(t, 4, len(mr))

	assert.Equal(t, "text/plain;format=flowed", mr[0].Value)
	assert.Equal(t, 1.0, mr[0].Weight)

	assert.Equal(t, "text/plain", mr[1].Value)
	assert.Equal(t, 0.9, mr[1].Weight)

	assert.Equal(t, "text/*", mr[2].Value)
	assert.Equal(t, 0.8, mr[2].Weight)

	assert.Equal(t, "*/*", mr[3].Value)
	assert.Equal(t, 0.7, mr[3].Weight)
}

func TestMediaRanges_should_handle_precedence4(t *testing.T) {
	a := new(accept)
	// from http://tools.ietf.org/html/rfc7231#section-5.3.1
	// and http://tools.ietf.org/html/rfc7231#section-5.3.2
	a.Header = "text/* ; q=0.3, text/html ; Q=0.7, text/html;level=1, text/html;level=2; q=0.4, */*; q=0.5"
	mr := a.ParseMediaRanges()

	assert.Equal(t, 5, len(mr))

	assert.Equal(t, "text/html;level=1", mr[0].Value)
	assert.Equal(t, 1.0, mr[0].Weight)

	assert.Equal(t, "text/html", mr[1].Value)
	assert.Equal(t, 0.7, mr[1].Weight)

	assert.Equal(t, "*/*", mr[2].Value)
	assert.Equal(t, 0.5, mr[2].Weight)

	assert.Equal(t, "text/html;level=2", mr[3].Value)
	assert.Equal(t, 0.4, mr[3].Weight)

	assert.Equal(t, "text/*", mr[4].Value)
	assert.Equal(t, 0.3, mr[4].Weight)
}
