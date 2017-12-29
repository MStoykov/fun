package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/dop251/goja_nodejs/console"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/robertkrimen/otto"
)

func main() {
	benchmarkOttoParallel()
	benchmarkGojaParallel()
	// benchmarkGolangParallel()

	octaneBenchmark()
}

func octaneBenchmark() {
	benchmarks := []string{
		"octane-benchmark/base.js",
		"octane-benchmark/box2d.js",
		"octane-benchmark/code-load.js",
		"octane-benchmark/crypto.js",
		"octane-benchmark/deltablue.js",
		"octane-benchmark/earley-boyer.js",
		"octane-benchmark/gbemu-part1.js",
		"octane-benchmark/gbemu-part2.js",
		"octane-benchmark/mandreel.js",
		"octane-benchmark/navier-stokes.js",
		"octane-benchmark/pdfjs.js",
		"octane-benchmark/raytrace.js",
		//		"octane-benchmark/regexp.js",
		"octane-benchmark/richards.js",
		"octane-benchmark/run.js",
		"octane-benchmark/splay.js",
		//"octane-benchmark/typescript-compiler.js",
		//"octane-benchmark/typescript-input.js",
		//"octane-benchmark/typescript.js",
		// "octane-benchmark/zlib-data.js",
		// "octane-benchmark/zlib.js",
	}
	benchmarkOtto(benchmarks)
	benchmarkGoja(benchmarks)
}

func mustReadFile(source string) string {
	buf, err := ioutil.ReadFile(source)
	if err != nil {
		log.Fatal("Failed to read file: %s", source)
		return ""
	}
	return string(buf)
}

func benchmarkOtto(benchmarks []string) {
	log.Println("benchmark otto\n")
	vm := otto.New()
	for _, source := range benchmarks {
		script, err := vm.Compile(source, mustReadFile(source))
		if err != nil {
			log.Fatal("error code : %s %s", source, err)
		}
		vm.Run(script)
	}
	source := "benchmark_otto.js"
	script, err := vm.Compile(source, mustReadFile(source))
	if err != nil {
		log.Fatal("error code : %s %s", source, err)
	}
	vm.Run(script)
}

func benchmarkGoja(benchmarks []string) {
	log.Println("benchmark goja\n")
	vm := goja.New()
	new(require.Registry).Enable(vm)
	console.Enable(vm)

	for _, source := range benchmarks {
		vm.RunScript(source, mustReadFile(source))
	}
	source := "benchmark_otto.js"
	vm.RunScript(source, mustReadFile(source))
}

func benchmarkOttoParallel() {
	fmt.Println("otto parallel test")
	var wg sync.WaitGroup
	now := time.Now()

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			vm := otto.New()
			source := `
		var sum = 0;
		for(var i = 0 ; i < 10000000 ; i++){
			sum++;
		}
	`
			vm.Run(source)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("otto parallel thread :", time.Now().Sub(now))
}

func benchmarkGojaParallel() {
	fmt.Println("goja parallel test")
	var wg sync.WaitGroup
	now := time.Now()

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			vm := goja.New()
			source := `
		var sum = 0;
		for(var i = 0 ; i < 10000000 ; i++){
			sum++;
		}
	`
			vm.RunScript("summ", source)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("goja parallel thread :", time.Now().Sub(now))
}

func benchmarkGolangParallel() {
	fmt.Println("golang parallel test")
	var wg sync.WaitGroup
	now := time.Now()
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			sum := 0
			for j := 0; j < 10000000; j++ {
				sum++
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("golang parallel thread :", time.Now().Sub(now))
}
