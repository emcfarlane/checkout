package checkout

import (
	"strconv"
	"testing"
)

func TestParsePAN(t *testing.T) {
	tests := []struct {
		pan     string
		want    string
		wantErr bool
	}{{
		pan:  "4000 0000 0000 0119",
		want: "4000000000000119",
	}, {
		pan:  " 4000-0000_0000,0119",
		want: "4000000000000119",
	}}

	for i, tt := range tests {
		t.Run("pan"+strconv.Itoa(i), func(t *testing.T) {
			got, err := parsePAN(tt.pan)
			if err != nil {
				if tt.wantErr {
					return // okay
				}
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("%s != %s", got, tt.want)
			}
		})
	}
}

func TestCheckLuhn(t *testing.T) {
	tests := []struct {
		pan     string
		wantErr bool
	}{{
		pan:     "",
		wantErr: true,
	}, {
		pan: "4000000000000119",
	}, {
		pan:     "4100000000000119",
		wantErr: true,
	}}

	for i, tt := range tests {
		t.Run("pan"+strconv.Itoa(i), func(t *testing.T) {
			ok := checkLuhn(tt.pan)
			if ok != !tt.wantErr {
				t.Fatalf("check %t want %t", ok, !tt.wantErr)
			}
		})
	}
}

func TestParseCurrencies(t *testing.T) {
	tests := []struct {
		currency string
		want     string
		wantErr  bool
	}{{
		wantErr: true,
	}, {
		currency: "aaa",
		wantErr:  true,
	}, {
		currency: "GbP",
		want:     "gbp",
	}}

	for i, tt := range tests {
		t.Run("currency"+strconv.Itoa(i), func(t *testing.T) {
			got, err := parseCurrency(tt.currency)
			if err != nil {
				if tt.wantErr {
					return // okay
				}
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("%s != %s", got, tt.want)
			}
		})
	}
}
