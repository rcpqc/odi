package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/rcpqc/odi/odi"
	"github.com/rcpqc/odi/resolve"
	"github.com/rcpqc/odi/test/cases/case1"
	"github.com/rcpqc/odi/test/cases/case2"
	_ "github.com/rcpqc/odi/test/cases/case3"
	_ "github.com/rcpqc/odi/test/cases/case4"
	"github.com/rcpqc/odi/test/cases/case5"
	"github.com/rcpqc/odi/test/cases/case6"
	_ "github.com/rcpqc/odi/test/cases/case7"
	"github.com/rcpqc/odi/test/cases/case8"
	"github.com/rcpqc/odi/test/config"
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
					&case2.E{DFG: "[xyz]", FF: &struct {
						VC []int "json:\"vc\""
					}{VC: []int{1, 23}}},
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
			source: map[any]any{"object": "case8_j", "arr_bool": []any{nil, int(-3), uint(0), 3.53, "TRUE", new(float32)}},
			want:   &case8.J{ArrBool: []bool{false, true, false, true, true, false}},
		},
		{
			name:   "case8.2",
			source: map[any]any{"object": "case8_j", "arr_bool": []any{3, "sdf"}},
			errR:   fmt.Errorf("_.arr_bool[1]: string(sdf) can't convert to boolean"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Resolve
			obj, err := odi.Resolve(tt.source, tt.opts...)
			if !reflect.DeepEqual(obj, tt.want) || !ErrorEqual(err, tt.errR) {
				t.Errorf("Resolve().obj \nresult: %v\nexpect: %v", obj, tt.want)
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
