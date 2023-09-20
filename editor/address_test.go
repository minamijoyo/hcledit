package editor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateAddressFromString(t *testing.T) {
	cases := []struct {
		name    string
		address string
		want    []string
	}{
		{
			name:    "simple address",
			address: "b1",
			want:    []string{"b1"},
		},
		{
			name:    "attribute address",
			address: "b1.l1.a1",
			want:    []string{"b1", "l1", "a1"},
		},
		{
			name:    "escaped address",
			address: `b1.l\.1.a1`,
			want:    []string{"b1", "l.1", "a1"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := createAddressFromString(tc.address)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff(-want +got):\n%v", got, tc.want, diff)
			}
		})
	}
}

func TestCreateStringFromAddress(t *testing.T) {
	cases := []struct {
		name    string
		address []string
		want    string
	}{
		{
			name:    "simple address",
			address: []string{"b1"},
			want:    "b1",
		},
		{
			name:    "simple address",
			address: []string{"b1", "l1", "a1"},
			want:    "b1.l1.a1",
		},
		{
			name:    "simple address",
			address: []string{"b1", `l.1`, "a1"},
			want:    `b1.l\.1.a1`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := createStringFromAddress(tc.address)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff(-want +got):\n%v", got, tc.want, diff)
			}
		})
	}
}
