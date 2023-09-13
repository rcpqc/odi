package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/rcpqc/odi/odi"
	"github.com/rcpqc/odi/resolve"
	"github.com/rcpqc/odi/test/cases/case1"
	"github.com/rcpqc/odi/test/cases/case10"
	"github.com/rcpqc/odi/test/cases/case11"
	"github.com/rcpqc/odi/test/cases/case12"
	"github.com/rcpqc/odi/test/cases/case14"
	"github.com/rcpqc/odi/test/cases/case2"
	"github.com/rcpqc/odi/test/cases/case5"
	"github.com/rcpqc/odi/test/cases/case6"
	"github.com/rcpqc/odi/test/cases/case8"
	"github.com/rcpqc/odi/test/cases/case9"
	"github.com/rcpqc/odi/test/config"
	"gopkg.in/yaml.v3"

	_ "github.com/rcpqc/odi/test/cases/case13"
	_ "github.com/rcpqc/odi/test/cases/case3"
	_ "github.com/rcpqc/odi/test/cases/case4"
	_ "github.com/rcpqc/odi/test/cases/case7"
)

func ErrorEqual(err1 error, err2 error) bool {
	if err1 == nil && err2 == nil {
		return true
	}
	if err1 == nil || err2 == nil {
		return false
	}
	return err1.Error() == err2.Error()
}

func TestResolveAndDispose(t *testing.T) {
	tests := []struct {
		name   string
		source any
		opts   []resolve.Option
		want   any
		errR   error
		errD   error
	}{
		{
			name:   "case1",
			source: config.ReadYaml("cases/case1/cfg.yaml"),
			opts:   []resolve.Option{resolve.WithTagKey("yaml")},
			want: &case1.A{
				Other:   map[string]interface{}{"object": "case1_a", "arg9": 6, "arg10": "t3"},
				Arg0:    123,
				Arg1:    "fafdsa",
				Arg2:    []uint{1, 2, 3},
				ObjectD: case1.D{KK: "kk123"},
				Ifaces: []case1.Interface{
					&case1.B{XX: 123, YY: "aaf", ZZ: []uint{4, 5, 6}, WW: [2]float32{1.1, 4.3}},
					&case1.C{C: "abcde", E: map[string]int{"a": 3}, F: map[bool]string{true: "T"}},
				},
			},
			errD: fmt.Errorf("ObjectC Dispose"),
		},
		{
			name:   "case2",
			source: config.ReadJson("cases/case2/cfg.json"),
			opts:   []resolve.Option{resolve.WithObjectKey("obj"), resolve.WithTagKey("json")},
			want: &case2.A{
				Ifaces: []case2.Interface{
					&case2.E{
						DFG: "[xyz]",
						CX:  321,
						FF: &struct {
							VC []int "json:\"vc\""
						}{VC: []int{1, 23}},
					},
				},
			},
		},
		{
			name:   "case3",
			source: config.ReadYaml("cases/case3/cfg.yaml"),
			opts:   []resolve.Option{resolve.WithTagKey("yaml")},
			want:   nil,
			errR:   fmt.Errorf("_.arg0: can't convert kind(slice) to int"),
		},
		{
			name:   "case4",
			source: config.ReadYaml("cases/case4/cfg.yaml"),
			opts:   []resolve.Option{resolve.WithTagKey("yaml")},
			want:   nil,
			errR:   fmt.Errorf("_.ifaces[0]: container create err: kind(case4_f) not registered"),
		},
		{
			name:   "case5",
			source: config.ReadYaml("cases/case5/cfg.yaml"),
			opts:   []resolve.Option{resolve.WithTagKey("yaml")},
			want: &case5.D{
				KK: "432",
				B: case5.B{
					XX: 123,
					ZZ: []uint{2, 4},
				},
			},
		},
		{
			name:   "case6.1",
			source: map[string]any{"step": []uint{1, 2}},
			opts:   []resolve.Option{resolve.WithObjectKey("step")},
			want:   nil,
			errR:   fmt.Errorf("_: classify err: kind must be a string"),
		},
		{
			name:   "case6.2",
			source: map[string]any{"step": "case6_g", "mc": map[string][]uint{"x1": {}}},
			opts:   []resolve.Option{resolve.WithObjectKey("step")},
			want:   nil,
			errR:   fmt.Errorf("_.mc[x1]: expect map but slice"),
		},
		{
			name:   "case6.3",
			source: map[any]any{"object": "case6_g", true: "123", false: true},
			opts:   []resolve.Option{resolve.WithStructFieldNameCompatibility(true)},
			want:   &case6.G{True: 123, False: 1.0},
		},
		{
			name:   "case7",
			source: map[any]any{"object": "case7_h", "a": "fds", "xx": 123},
			want:   nil,
			errR:   fmt.Errorf("_.h1.b: illegal inline type(*case7.B) expect struct or map[string]any"),
		},
		{
			name:   "case8.1",
			source: map[any]any{"object": "case8_bools", "arr_bool": []any{nil, int(-3), uint(0), 3.53, "TRUE", new(float32)}},
			want:   &case8.Bools{ArrBool: []bool{false, true, false, true, true, false}},
		},
		{
			name:   "case8.2",
			source: map[any]any{"object": "case8_bools", "arr_bool": []any{3, "sdf"}},
			errR:   fmt.Errorf("_.arr_bool[1]: can't convert string(sdf) to bool"),
		},
		{
			name:   "case8.3",
			source: map[any]any{"object": "case8_bools", "arr_bool": []any{(*int)(nil)}},
			errR:   fmt.Errorf("_.arr_bool[0]: can't convert nil pointer to bool"),
		},
		{
			name:   "case8.4",
			source: map[any]any{"object": "case8_bools", "arr_bool": []any{complex(1, 2)}},
			errR:   fmt.Errorf("_.arr_bool[0]: can't convert kind(complex128) to bool"),
		},
		{
			name:   "case9.1",
			source: map[any]any{"object": "case9_ints", "arr_int": []any{nil, "-3", uint(0), 3.53, new(float32), true, false}},
			want:   &case9.Ints{ArrInt: []int{0, -3, 0, 3, 0, 1, 0}},
		},
		{
			name:   "case9.2",
			source: map[any]any{"object": "case9_ints", "arr_int": []any{3, "sdf"}},
			errR:   fmt.Errorf("_.arr_int[1]: can't convert string(sdf) to int"),
		},
		{
			name:   "case9.3",
			source: map[any]any{"object": "case9_ints", "arr_int": []any{(*int)(nil)}},
			errR:   fmt.Errorf("_.arr_int[0]: can't convert nil pointer to int"),
		},
		{
			name:   "case10.1",
			source: map[any]any{"object": "case10_uints", "arr_uint": []any{nil, "54", uint(0), 3.53, new(float32), true, false}},
			want:   &case10.Uints{ArrUint: []uint{0, 54, 0, 3, 0, 1, 0}},
		},
		{
			name:   "case10.2",
			source: map[any]any{"object": "case10_uints", "arr_uint": []any{3, "sdf"}},
			errR:   fmt.Errorf("_.arr_uint[1]: can't convert string(sdf) to uint"),
		},
		{
			name:   "case10.3",
			source: map[any]any{"object": "case10_uints", "arr_uint": []any{(*int)(nil)}},
			errR:   fmt.Errorf("_.arr_uint[0]: can't convert nil pointer to uint"),
		},
		{
			name:   "case10.4",
			source: map[any]any{"object": "case10_uints", "arr_uint": []any{complex(1, 2)}},
			errR:   fmt.Errorf("_.arr_uint[0]: can't convert kind(complex128) to uint"),
		},
		{
			name:   "case11.1",
			source: map[any]any{"object": "case11_floats", "arr_float": []any{nil, "55.2", 2, uint(0), new(float32), true, false}},
			want:   &case11.Floats{ArrFloat: []float64{0, 55.2, 2, 0, 0, 1, 0}},
		},
		{
			name:   "case11.2",
			source: map[any]any{"object": "case11_floats", "arr_float": []any{3, "sdf"}},
			errR:   fmt.Errorf("_.arr_float[1]: can't convert string(sdf) to float"),
		},
		{
			name:   "case11.3",
			source: map[any]any{"object": "case11_floats", "arr_float": []any{(*int)(nil)}},
			errR:   fmt.Errorf("_.arr_float[0]: can't convert nil pointer to float"),
		},
		{
			name:   "case11.4",
			source: map[any]any{"object": "case11_floats", "arr_float": []any{complex(1, 2)}},
			errR:   fmt.Errorf("_.arr_float[0]: can't convert kind(complex128) to float"),
		},
		{
			name:   "case12.1",
			source: map[any]any{"object": "case12_strings", "arr_string": []any{nil, -2, uint(0), 23.443, new(float32), true, false}},
			want:   &case12.Strings{ArrString: []string{"", "-2", "0", "23.443", "0", "true", "false"}},
		},
		{
			name:   "case12.2",
			source: map[any]any{"object": "case12_strings", "arr_string": []any{(*int)(nil)}},
			errR:   fmt.Errorf("_.arr_string[0]: can't convert nil pointer to string"),
		},
		{
			name:   "case12.3",
			source: map[any]any{"object": "case12_strings", "arr_string": []any{complex(1, 2)}},
			errR:   fmt.Errorf("_.arr_string[0]: can't convert kind(complex128) to string"),
		},
		{
			name:   "case13.1",
			source: nil,
			errR:   fmt.Errorf("_: classify err: expect map but invalid"),
		},
		{
			name:   "case13.2",
			source: []any{map[any]any{"object": "case13_a"}},
			errR:   fmt.Errorf("_: classify err: expect map but slice"),
		},
		{
			name:   "case13.3",
			source: map[any]any{"xxx": "case13_a"},
			errR:   fmt.Errorf("_: classify err: not exist kind field(object)"),
		},
		{
			name:   "case13.4",
			source: map[any]any{"object": "case13_b", "channel": make(chan int)},
			errR:   fmt.Errorf("_.channel: not support kind: chan"),
		},
		{
			name:   "case13.5",
			source: map[any]any{"object": "case13_c", "x": "sdfsdf"},
			errR:   fmt.Errorf("_.x: expect slice or array but string"),
		},
		{
			name:   "case13.6",
			source: map[any]any{"object": "case13_c", "y": "sdfsdf"},
			errR:   fmt.Errorf("_.y: expect slice or array but string"),
		},
		{
			name:   "case13.7",
			source: map[any]any{"object": "case13_c", "z": "sdfsdf"},
			errR:   fmt.Errorf("_.z: expect map but string"),
		},
		{
			name:   "case13.8",
			source: map[any]any{"object": "case13_c", "i": map[any]any{"object": "case13_a"}},
			errR:   fmt.Errorf("_.i: the injected object does not implement the interface(case13.Iter)"),
		},
		{
			name:   "case13.9",
			source: map[any]any{"object": "case13_c", "x": []any{"1231", 423, func() {}}},
			errR:   fmt.Errorf("_.x[2]: can't convert kind(func) to string"),
		},
		{
			name:   "case13.10",
			source: map[any]any{"object": "case13_c", "z": map[any]any{complex(1, 2): 4}},
			errR:   fmt.Errorf("_.z[]: can't convert kind(complex128) to string"),
		},
		{
			name:   "case13.11",
			source: map[any]any{"object": "case13_c", "x": []any{"1231", 423}},
			errR:   fmt.Errorf("_.x: expect array's length to be 3 but 2"),
		},
		{
			name:   "case14",
			source: map[any]any{"object": "case14_a", "m1": nil, "s2": nil, "b2": nil},
			want:   &case14.A{M1: map[string]string{}, S2: []int{}, B2: &case14.B{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Resolve
			obj, err := odi.Resolve(tt.source, tt.opts...)
			if !reflect.DeepEqual(obj, tt.want) || !ErrorEqual(err, tt.errR) {
				t.Errorf("Resolve().obj \nresult: %v\nexpect: %v", obj, tt.want)
				bytesobj, _ := yaml.Marshal(obj)
				t.Errorf(string(bytesobj))
				t.Errorf("Resolve().err \nresult: %v\nexpect: %v", err, tt.errR)
				return
			}
			// Dispose
			if err := odi.Dispose(obj); !ErrorEqual(err, tt.errD) {
				t.Errorf("Dispose().err \nresult: %v\nexpect: %v", err, tt.errD)
			}
		})
	}
}

func BenchmarkResolve(b *testing.B) {
	source := config.ReadYaml("../test/cases/1.yaml")
	opts := []resolve.Option{resolve.WithTagKey("yaml")}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = odi.Resolve(source, opts...)
	}
}
