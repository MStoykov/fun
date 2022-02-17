package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/dop251/goja_nodejs/console"
	v8 "rogchap.com/v8go"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/robertkrimen/otto"
)

func main() {
	// benchmarkOttoParallel()
	// benchmarkGojaParallel()
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
	benchmarkV8Worker(benchmarks)
	benchmarkOtto(benchmarks)
	benchmarkGoja(benchmarks)
}

func mustReadFile(source string) string {
	buf, err := ioutil.ReadFile(source)
	if err != nil {
		log.Fatalf("Failed to read file: %s", source)
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
	log.Println("benchmark goja")
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
			vm := goja.New()
			source := `
		var sum = 0;
		for(var i = 0 ; i < 10000000 ; i++){
			sum++;
		}
	`
			vm.RunString(source)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("otto parallel thread :", time.Now().Sub(now))
}

func benchmarkV8Worker(benchmarks []string) {
	log.Println("benchmark v8Worker\n")
	iso := v8.NewIsolate()
	printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		fmt.Printf("%v\n", info.Args()) // when the JS function is called this Go callback will execute
		return nil                      // you can return a value back to the JS caller if required
	})
	global := v8.NewObjectTemplate(iso) // a template that represents a JS Object
	global.Set("$print", printfn)       // sets the "print" property of the Object to our function
	ctx := v8.NewContext(iso, global)   // new Context with the global Object set to our object template
	for _, source := range benchmarks {
		_, err := ctx.RunScript(mustReadFile(source), source)
		if err != nil {
			fmt.Println(err)
		}
	}
	source := "benchmark_v8.js"
	_, err := ctx.RunScript(mustReadFile(source), source)
	if err != nil {
		fmt.Println(err)
	}
}
