BASEPATH := ./matrix
OUTPATH := $(BASEPATH)/out

bench/simple:
	go test $(BASEPATH) -bench=BenchmarkSimpleProduct -benchmem \
 		-cpuprofile $(OUTPATH)/simple/cpuprofile.out \
 		-memprofile $(OUTPATH)/simple/memprofile.out \
 		-trace $(OUTPATH)/simple/traceprofile.out \
 		-o $(OUTPATH)/simple/matrix.test

bench/simple/performance:
	go test $(BASEPATH) -bench=BenchmarkSimpleProduct -benchmem

bench/simd:
	GOEXPERIMENT=simd go test $(BASEPATH) -bench=BenchmarkSIMDProduct -benchmem \
 		-cpuprofile $(OUTPATH)/simd/cpuprofile.out \
 		-memprofile $(OUTPATH)/simd/memprofile.out \
 		-trace $(OUTPATH)/simd/traceprofile.out \
 		-o $(OUTPATH)/simd/matrix.test

bench/simd/performance:
	GOEXPERIMENT=simd go test $(BASEPATH) -bench=BenchmarkSIMDProduct -benchmem -count=3

bench/simd/amd64:
	GOEXPERIMENT=simd go test $(BASEPATH) -bench=BenchmarkSIMDAMD64Product -benchmem \
 		-cpuprofile $(OUTPATH)/simd-amd64/cpuprofile.out \
 		-memprofile $(OUTPATH)/simd-amd64/memprofile.out \
 		-trace $(OUTPATH)/simd-amd64/traceprofile.out \
 		-o $(OUTPATH)/simd-amd64/matrix.test

bench/simd/amd64/performance:
	GOEXPERIMENT=simd go test $(BASEPATH) -bench=BenchmarkSIMDAMD64Product -benchmem -count=3
