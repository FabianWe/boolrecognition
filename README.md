# boolrecognition
Recognitation and transformation procedures for boolean functions written in Go.

## Installation
First you have to install [lpsolve](http://lpsolve.sourceforge.net/) (the linear program solver) and the Go bindings for lpsolve. Due to copyright problems lpsolve is not shipped with boolrecognition. You can find the installation instructions [here](https://github.com/draffensperger/golp/). After this just use `go get github.com/FabianWe/boolrecognition/...` and then build the binary with `go build cmd/benchmarklpb/benchmarklpb.go` (from the directory `FabianWe/boolrecognition`).

## Usage
Currenty benchmarklp accepts text files where each line contains an LPB in the format:
First all the coefficients are separated by a space, then the threshold follows, also separated by a space.

So the LPB 2 ⋅ x1 + 1 ⋅ x2 + 1 ⋅ x3 ≥ 2 is represented by "2 1 1 2".

You can find benchmarks [here](https://github.com/FabianWe/lpb_benchmarks).  Example:

    ./benchmarklpb -lpb lpb_benchmarks/full/lpb/full_6.lpb -verify

(adjust the path to the lpb file). This uses the combinatorial solver (known to be not complete). To use the linear program solver use

    ./benchmarklpb -lpb lpb_benchmarks/full/lpb/full_6.lpb -verify -solver lp

For more options see `./benchmarklpb -help`.
