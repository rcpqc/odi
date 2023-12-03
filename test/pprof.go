package main

import (
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/rcpqc/odi"
	"github.com/rcpqc/odi/resolve"
	"github.com/rcpqc/odi/test/config"

	_ "github.com/rcpqc/odi/test/cases/case1"
)

func main() {
	f, _ := os.Create("cpu.profile")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	source := config.ReadYaml("cases/case1/cfg.yaml") // 13401
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
