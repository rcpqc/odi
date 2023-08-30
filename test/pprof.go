package main

import (
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/rcpqc/odi/odi"
	"github.com/rcpqc/odi/resolve"
	"github.com/rcpqc/odi/test/config"
	"github.com/rcpqc/odi/test/objects"
)

func init() {
	odi.Provide("object_a", func() any { return &objects.A{} })
	odi.Provide("object_b", func() any { return &objects.B{} })
	odi.Provide("object_c", func() any { return &objects.C{} })
	odi.Provide("object_d", func() any { return &objects.D{} })
	odi.Provide("object_e", func() any { return &objects.E{} })
	odi.Provide("object_g", func() any { return &objects.G{} })
}

func main() {
	f, _ := os.Create("cpu.profile")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	source := config.ReadYaml("test/cases/1.yaml")
	opts := []resolve.Option{resolve.WithTagKey("yaml")}
	obj, err := odi.Resolve(source, opts...)
	log.Print(obj)
	log.Print(err)
	n := 200000
	st := time.Now()
	for i := 0; i < n; i++ {
		_, _ = odi.Resolve(source, opts...)
	}
	log.Printf("%vns/op", float64(time.Since(st).Nanoseconds())/float64(n))
}
