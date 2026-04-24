package updater

import "testing"

func TestIsNewer(t *testing.T) {
	cases := []struct {
		latest, current string
		want            bool
	}{
		{"1.2.3", "1.2.2", true},
		{"1.2.3", "1.2.3", false},
		{"1.2.3", "1.3.0", false},
		{"2.0.0", "1.99.99", true},
		{"1.2.3", "dev", true},
		{"1.2.3", "", true},
		{"1.2.3-rc1", "1.2.2", true},
		{"1.2.3", "v1.2.3", false},
	}
	for _, c := range cases {
		if got := isNewer(c.latest, c.current); got != c.want {
			t.Errorf("isNewer(%q,%q)=%v want %v", c.latest, c.current, got, c.want)
		}
	}
}
