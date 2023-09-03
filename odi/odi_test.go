package odi

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/rcpqc/odi/resolve"
	"github.com/rcpqc/odi/test/config"
	"github.com/rcpqc/odi/test/objects"
)

func init() {
	Provide("object_a", func() any { return &objects.A{} })
	Provide("object_b", func() any { return &objects.B{} })
	Provide("object_c", func() any { return &objects.C{} })
	Provide("object_d", func() any { return &objects.D{} })
	Provide("object_e", func() any { return &objects.E{} })
	Provide("object_g", func() any { return &objects.G{} })
	Provide("object_h", func() any { return &objects.H{} })
	Provide("object_j", func() any { return &objects.J{} })
}

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
			name:   "case_1",
			source: config.ReadYaml("../test/cases/1.yaml"),
			opts:   []resolve.Option{resolve.WithTagKey("yaml")},
			want: &objects.A{
				Other:   map[string]interface{}{"object": "object_a", "arg9": 6, "arg10": "t3"},
				Arg0:    123,
				Arg1:    "fafdsa",
				Arg2:    []uint{1, 2, 3},
				ObjectD: objects.D{KK: "kk123"},
				Ifaces: []objects.Interface1{
					&objects.B{XX: 123, YY: "aaf", ZZ: []uint{4, 5, 6}, WW: [2]float32{1.1, 4.3}},
					&objects.C{C: "abcde", E: map[string]int{"a": 3}, F: map[bool]string{true: "T"}},
				},
			},
			errD: fmt.Errorf("ObjectC Dispose"),
		},
		{
			name:   "case_2",
			source: config.ReadJson("../test/cases/2.json"),
			opts:   []resolve.Option{resolve.WithObjectKey("obj"), resolve.WithTagKey("json")},
			want: &objects.A{
				Ifaces: []objects.Interface1{
					&objects.E{DFG: "[xyz]", FF: &struct {
						VC []int "json:\"vc\""
					}{VC: []int{1, 23}}},
				},
			},
		},
		{
			name:   "case_3",
			source: config.ReadYaml("../test/cases/3.yaml"),
			opts:   []resolve.Option{resolve.WithTagKey("yaml")},
			want:   nil,
			errR:   fmt.Errorf("_.arg0: can't convert kind(slice) to int"),
		},
		{
			name:   "case_4",
			source: config.ReadYaml("../test/cases/4.yaml"),
			opts:   []resolve.Option{resolve.WithTagKey("yaml")},
			want:   nil,
			errR:   fmt.Errorf("_.ifaces[0]: container create err: kind(object_f) not registered"),
		},
		{
			name:   "case_5",
			source: config.ReadYaml("../test/cases/5.yaml"),
			opts:   []resolve.Option{resolve.WithTagKey("yaml")},
			want: &objects.D{
				KK: "432",
				B: objects.B{
					XX: 123,
					ZZ: []uint{2, 4},
				},
			},
		},
		{
			name:   "case_6",
			source: map[string]any{"step": []uint{1, 2}},
			opts:   []resolve.Option{resolve.WithObjectKey("step")},
			want:   nil,
			errR:   fmt.Errorf("_: classify err: kind must be a string"),
		},
		{
			name:   "case_7",
			source: map[string]any{"step": "object_g", "mc": map[string][]uint{"x1": {}}},
			opts:   []resolve.Option{resolve.WithObjectKey("step")},
			want:   nil,
			errR:   fmt.Errorf("_.mc[x1]: expect map but slice"),
		},
		{
			name:   "case_8",
			source: map[any]any{"object": "object_g", true: "123", false: true},
			opts:   []resolve.Option{resolve.WithStructFieldNameCompatibility(true)},
			want:   &objects.G{True: 123, False: 1.0},
		},
		{
			name:   "case_9",
			source: map[any]any{"object": "object_h", "a": "fds", "xx": 123},
			want:   nil,
			errR:   fmt.Errorf("_.h1.b: illegal inline type(*objects.B) expect struct or map[string]any"),
		},
		{
			name:   "case_10",
			source: map[any]any{"object": "object_j", "arr_bool": []any{nil, int(-3), uint(0), 3.53, "TRUE", new(float32)}},
			want:   &objects.J{ArrBool: []bool{false, true, false, true, true, false}},
		},
		{
			name:   "case_11",
			source: map[any]any{"object": "object_j", "arr_bool": []any{3, "sdf"}},
			errR:   fmt.Errorf("_.arr_bool[1]: string(sdf) can't convert to boolean"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Resolve
			obj, err := Resolve(tt.source, tt.opts...)
			if !reflect.DeepEqual(obj, tt.want) || !ErrorEqual(err, tt.errR) {
				t.Errorf("Resolve().obj \nresult: %v\nexpect: %v", obj, tt.want)
				t.Errorf("Resolve().err \nresult: %v\nexpect: %v", err, tt.errR)
				return
			}
			// Dispose
			if err := Dispose(obj); !ErrorEqual(err, tt.errD) {
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
		_, _ = Resolve(source, opts...)
	}
}
