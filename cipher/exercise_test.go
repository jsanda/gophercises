package cipher

import "testing"

func TestEncode(t *testing.T) {
	tables := []struct {
		s string
		k int
		z string
	} {
		{"ABC", 1, "BCD"},
		{"EZ", 3, "HC"},
		{"abc", 1, "bcd"},
		{"ez", 3, "hc"},
		{"z", 40, "n"},
		{"Cx-m3R", 57, "Hc-r3W"},
	}
	for _, table := range tables {
		if encoded := Encode(table.s, table.k); encoded != table.z {
			t.Errorf("Encode(%s, %d) was wrong, expected %s, got %s", table.s, table.k, table.z, encoded)
		}
	}
}