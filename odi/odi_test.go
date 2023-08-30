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
		name       string
		source     any
		objKey     string
		tagKey     string
		want       any
		resolveErr error
		disposeErr error
	}{
		{
			name:   "case_1",
			source: config.ReadYaml("../test/cases/1.yaml"),
			tagKey: "yaml",
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
			disposeErr: fmt.Errorf("ObjectC Dispose"),
		},
		{
			name:   "case_2",
			source: config.ReadJson("../test/cases/2.json"),
			objKey: "obj",
			tagKey: "json",
			want: &objects.A{
				Ifaces: []objects.Interface1{
					&objects.E{DFG: "[xyz]", FF: &struct {
						VC []int "json:\"vc\""
					}{VC: []int{1, 23}}},
				},
			},
		},
		{
			name:       "case_3",
			source:     config.ReadYaml("../test/cases/3.yaml"),
			tagKey:     "yaml",
			want:       nil,
			resolveErr: fmt.Errorf("_.arg0: can't convert kind(slice) to int"),
		},
		{
			name:       "case_4",
			source:     config.ReadYaml("../test/cases/4.yaml"),
			tagKey:     "yaml",
			want:       nil,
			resolveErr: fmt.Errorf("_.ifaces[0]: container create err: kind(object_f) not registered"),
		},
		{
			name:   "case_5",
			source: config.ReadYaml("../test/cases/5.yaml"),
			tagKey: "yaml",
			want: &objects.D{
				KK: "432",
				B: objects.B{
					XX: 123,
					ZZ: []uint{2, 4},
				},
			},
		},
		{
			name:       "case_6",
			source:     map[string]any{"step": []uint{1, 2}},
			objKey:     "step",
			want:       nil,
			resolveErr: fmt.Errorf("_: classify err: kind must be a string"),
		},
		{
			name:       "case_7",
			source:     map[string]any{"step": "object_g", "mc": map[string][]uint{"x1": {}}},
			objKey:     "step",
			want:       nil,
			resolveErr: fmt.Errorf("_.mc[x1]: expect map but slice"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []resolve.Option{}
			if tt.objKey != "" {
				opts = append(opts, resolve.WithObjectKey(tt.objKey))
			}
			if tt.tagKey != "" {
				opts = append(opts, resolve.WithTagKey(tt.tagKey))
			}

			// Resolve
			got, err := Resolve(tt.source, opts...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resolve() res\ngot:  %v\nwant: %v", got, tt.want)
			}
			if !ErrorEqual(err, tt.resolveErr) {
				t.Errorf("Resolve() err\ngot:  %v\nwant: %v", err, tt.resolveErr)
			}

			// Dispose
			if err := Dispose(got); !ErrorEqual(err, tt.disposeErr) {
				t.Errorf("Dispose() err\ngot:  %v\nwant: %v", err, tt.disposeErr)
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
