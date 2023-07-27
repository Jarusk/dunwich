package cluster

import (
	"net/netip"
	"reflect"
	"testing"
)

func TestNodeListToStringList(t *testing.T) {
	tests := map[string]struct {
		input    []netip.AddrPort
		expected []string
	}{
		"empty": {
			input:    []netip.AddrPort{},
			expected: []string{},
		},
		"one": {
			input: []netip.AddrPort{
				netip.MustParseAddrPort("192.168.0.10:1234"),
			},
			expected: []string{
				"192.168.0.10:1234",
			},
		},
	}
	t.Parallel()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res := nodeListToStringList(tc.input)
			if !reflect.DeepEqual(res, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, res)
			}
		})
	}
}
