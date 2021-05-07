package binance

import (
	"net/url"
	"strings"
)

// Parameters is required only to skip the regular key sorting of the url.Values type
type Parameters interface {
	Encode() string
	Set(string, ...string)
}

type parameters struct {
	keys   []string
	values [][]string
}

// Encode encodes the values into ``URL encoded'' form
// ("bar=baz&foo=quux") sorted by order of entry.
func (p *parameters) Encode() string {
	var buf strings.Builder
	// encode in reversed order
	for i, k := range p.keys {
		vs := p.values[i]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(v)
		}
	}
	return buf.String()
}

// Set sets the key to value. It replaces any existing
// values.
func (p *parameters) Set(key string, value ...string) {
	if n := pos(p.keys, key); n > -1 {
		p.values[n] = value
	} else {
		p.keys = append(p.keys, key)
		p.values = append(p.values, value)
	}
}

// NewParameters returns a new parameter bag
func NewParameters(initialLength int) Parameters {
	return &parameters{
		keys:   make([]string, 0, initialLength),
		values: make([][]string, 0, initialLength),
	}
}

func pos(list []string, s string) int {
	for i, k := range list {
		if k == s {
			return i
		}
	}
	return -1
}
