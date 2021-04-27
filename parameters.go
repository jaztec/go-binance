package binance

import (
	"net/url"
	"strings"
)

// Parameters is required only to skip the regular key sorting of the url.Values type
type Parameters map[string][]string

// Encode encodes the values into ``URL encoded'' form
// ("bar=baz&foo=quux") sorted by key.
func (v Parameters) Encode() string {
	if v == nil {
		return ""
	}
	// reverse the map
	rev := make(Parameters, len(v))
	for k, val := range v {
		rev[k] = val
	}
	var buf strings.Builder
	keys := make([]string, 0, len(rev))
	for k := range rev {
		keys = append(keys, k)
	}
	for _, k := range keys {
		vs := v[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}

// Set sets the key to value. It replaces any existing
// values.
func (v Parameters) Set(key, value string) {
	v[key] = []string{value}
}
