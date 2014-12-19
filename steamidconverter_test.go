package steamidconverter

import (
	"testing"
)

var textTests = []struct {
	in uint64
	out string
}{
	{76561197960430077, "STEAM_0:1:82174"},
}

var sixtyFourTests = []struct {
	in string
	out uint64
}{
	{"STEAM_0:1:82174", 76561197960430077},
}

var steam3Tests = []struct {
	in string
	out string
}{
	{"STEAM_0:1:82174", "[U:1:164349]"},
}

func TestConvertToText(t *testing.T) {
	for _, tt := range textTests {
		sidText := ConvertToText(tt.in)
		if sidText != tt.out {
			t.Errorf("ConvertToText(%d) => %q, want %q", tt.in, sidText, tt.out)
		}
	}
}

func TestConvertTo64(t *testing.T) {
	for _, tt := range sixtyFourTests {
		sid64, err := ConvertTo64(tt.in)
		if err != nil {
			t.Error("ConvertTo64(%q) => _, %q", tt.in, err)
		}
		if sid64 != tt.out {
			t.Errorf("ConvertTo64(%q) => %d, want %d", tt.in, sid64, tt.out)
		}
	}
}

func TestConvertToSteam3(t *testing.T) {
	for _, tt := range steam3Tests {
		sid3, err := ConvertToSteam3(tt.in)
		if err != nil {
			t.Error("ConvertToSteam3(%q) => _, %q", tt.in, err)
		}
		if sid3 != tt.out {
			t.Errorf("ConvertToSteam3(%q) => %q, want %q", tt.in, sid3, tt.out)
		}
	}
}
