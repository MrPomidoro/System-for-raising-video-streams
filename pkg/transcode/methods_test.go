package transcode

import (
	"fmt"
	"reflect"
	"testing"
)

type testS struct {
	name string
	in   interface{}
}

func TestTranscode(t *testing.T) {
	tests := []testS{
		{
			name: "test maps",
			in:   map[int]int{1: 1, 2: 2},
		},
		{
			name: "test slice",
			in:   []int{1, 2},
		},
		{
			name: "test array",
			in:   [2]int{1, 2},
		},
		{
			name: "test string",
			in:   "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			switch in := tt.in.(type) {
			case map[int]int:
				out := make(map[int]int)
				Transcode(in, &out)
				if !reflect.DeepEqual(in, out) {
					t.Errorf("%s --- FAIL: wanted %v, got %v", tt.name, tt.in, out)
				} else {
					fmt.Printf("%s --- OK\n", tt.name)
				}
			case []int:
				out := []int{}
				Transcode(in, &out)
				if !reflect.DeepEqual(in, out) {
					t.Errorf("%s --- FAIL: wanted %v, got %v", tt.name, tt.in, out)
				} else {
					fmt.Printf("%s --- OK\n", tt.name)
				}
			case [2]int:
				out := [2]int{}
				Transcode(in, &out)
				if !reflect.DeepEqual(in, out) {
					t.Errorf("%s --- FAIL: wanted %v, got %v", tt.name, tt.in, out)
				} else {
					fmt.Printf("%s --- OK\n", tt.name)
				}
			case string:
				out := ""
				Transcode(in, &out)
				if !reflect.DeepEqual(in, out) {
					t.Errorf("%s --- FAIL: wanted %v, got %v", tt.name, tt.in, out)
				} else {
					fmt.Printf("%s --- OK\n", tt.name)
				}
			}
		})
	}
}
