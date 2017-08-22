// Copyright 2017 Fabian Wenzelmann <fabianwen@posteo.eu>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Just for playing around a bit and testing stuff.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	br "github.com/FabianWe/boolrecognition"
	"github.com/FabianWe/boolrecognition/lpb"
)

// iterativeAverage computes iteratively the average of a series of values.
// Implemented as described here: http://people.revoledu.com/kardi/tutorial/RecursiveStatistic/Time-Average.htm
// In contrast to the method described above t starts with 0, not 1.
// So to compute the average of a series do:
// 1. Initialize your current average to whatever you want
// 2. Initialize t = 0
// 3. For each sample update current = iterativeAverage(t, nextSample, current)
//    and increase t by 1.
func iterativeAverage(t int, value, current float64) float64 {
	return (float64(t)/float64(t+1))*current + (1.0/float64(t+1))*value
}

func main() {
	lpbFileFlag := flag.String("lpb", "", "Path to the lpb file")
	verify := flag.Bool("verify", false, "If true also verify that the produced LPB is correct")
	solverType := flag.String("solver", "minComb", "The solver to use, currently minComb is the only option (incomplete one introduced in the paper by Smaus)")
	numberLoops := flag.Int("N", 5, "The number of times you want to repeat each conversion")
	repeat := flag.Int("R", 3, "How many times to repeat the conversions N times? Best value will be used")
	flag.Parse()
	if *lpbFileFlag == "" {
		fmt.Fprintln(os.Stderr, "lpb must be provided and must point to the file containg all the LPBs")
		os.Exit(1)
	}
	if *solverType != "minComb" {
		fmt.Fprintln(os.Stderr, "Only \"minComb\" is available right now as solver")
		os.Exit(1)
	}
	if *numberLoops <= 0 {
		fmt.Fprintln(os.Stderr, "N must be > 0")
		os.Exit(1)
	}
	if *repeat <= 0 {
		fmt.Fprintln(os.Stderr, "R must be > 0")
		os.Exit(1)
	}
	lpbs, dnfs, parseErr := parseLPBs(*lpbFileFlag)
	var numFailedConv, numNotEqual int
	bestSoFar := -1.0
	if parseErr != nil {
		fmt.Fprintln(os.Stderr, "Error parsing LPBs:", parseErr)
		os.Exit(1)
	}
	for i := 0; i < *repeat; i++ {
		// repeat the test, get average
		var avg float64
		avg, numFailedConv, numNotEqual = runTest(lpbs, dnfs, *numberLoops, *verify)
		if bestSoFar < 0 || avg < bestSoFar {
			bestSoFar = avg
		}
	}
	// print evaluation
	fmt.Printf("Ran tests %d times, showing best average of %d repeats\n", *numberLoops, *repeat)
	failRate := (float64(numFailedConv) / float64(len(lpbs))) * 100.0
	fmt.Printf("Conversion failed on %d of %d tests (%.2f%%)\n", numFailedConv, len(lpbs), failRate)
	if *verify {
		errorRate := (float64(numNotEqual) / float64(len(lpbs)-numFailedConv)) * 100.0
		fmt.Printf("From the times the conversion was successful the output was wrong in %d cases (%.2f%%)\n", numNotEqual, errorRate)
	}
	fmt.Printf("One (single) conversion took %s on average\n", time.Duration(bestSoFar))
}

func parseLPBs(path string) ([]*lpb.LPB, []br.ClauseSet, error) {
	lpbs := make([]*lpb.LPB, 0)
	dnfs := make([]br.ClauseSet, 0)
	f, openErr := os.Open(path)
	if openErr != nil {
		return nil, nil, openErr
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		nextLPB, parseErr := lpb.ParseLPB(scanner.Text())

		if parseErr != nil {
			return nil, nil, parseErr
		}
		lpbs = append(lpbs, nextLPB)
		dnfs = append(dnfs, nextLPB.ToDNF())
	}
	return lpbs, dnfs, nil
}

func runTest(lpbs []*lpb.LPB, dnfs []br.ClauseSet, n int, verify bool) (avg float64, numFailedConv, numNotEqual int) {
	solver := lpb.NewMinSolver()
	avg = 0.0
	t := 0
	for num := 0; num < n; num++ {
		numFailedConv = 0
		numNotEqual = 0
		for i, phi := range dnfs {
			start := time.Now()
			// create tree and try to solve
			tree := lpb.NewSplittingTree(phi, len(lpbs[i].Coefficients), true, true)
			tree.Cut = true
			computedLPB, convErr := solver.Solve(tree)
			dur := time.Since(start)
			ok := true
			if convErr != nil {
				ok = false
				numFailedConv++
			}
			if verify && convErr == nil {
				if !computedLPB.ToDNF().DeepSortedEquals(phi) {
					ok = false
					numNotEqual++
				}
			}
			if ok {
				avg = iterativeAverage(t, float64(dur), avg)
				t++
			}
		}
	}
	return
}
