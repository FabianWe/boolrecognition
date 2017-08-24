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
	tighten := lpb.TightenNone
	lpbFileFlag := flag.String("lpb", "", "Path to the lpb file")
	verify := flag.Bool("verify", false, "If true also verify that the produced LPB is correct")
	solverType := flag.String("solver", "minComb", "The solver to use, currently \"mincomb\" and \"lp\" are available")
	numberLoops := flag.Int("N", 5, "The number of times you want to repeat each conversion")
	repeat := flag.Int("R", 3, "How many times to repeat the conversions N times? Best value will be used")
	tightenFlag := flag.String("tighten", "none", "If the solver is lp solver this describes how to tighten the lp:"+
		" \"none\" for now additional constraints, \"neighbours\" for constraints v(i) and v(i + 1) and \"all\""+
		" for constraints between all v(i) and v(j). Default is \"none\"")
	flag.Parse()
	var converter lpb.DNFToLPB
	if *lpbFileFlag == "" {
		fmt.Fprintln(os.Stderr, "lpb must be provided and must point to the file containg all the LPBs")
		os.Exit(1)
	}
	switch *solverType {
	case "minComb":
		converter = lpb.NewCombinatorialSolver(lpb.NewMinSolver())
	case "lp":
		switch *tightenFlag {
		case "none":
		case "neighbours":
			tighten = lpb.TightenNeighbours
		case "all":
			tighten = lpb.TightenAll
		default:
			fmt.Fprintln(os.Stderr, "Tighten type must be either \"none\", \"neighbours\" or \"all\", got", *tightenFlag)
			os.Exit(1)
		}
		converter = lpb.NewLPSolver(tighten)
	default:
		fmt.Fprintln(os.Stderr, "Only \"minComb\" and \"lp\" are valid solvers, got", *solverType)
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
	bestSoFarSucc := -1.0
	bestSoFarAll := -1.0
	if parseErr != nil {
		fmt.Fprintln(os.Stderr, "Error parsing LPBs:", parseErr)
		os.Exit(1)
	}
	for i := 0; i < *repeat; i++ {
		// repeat the test, get average
		var avgSucc, avgAll float64
		// run verify only in the last run, no need to always do it
		avgSucc, avgAll, numFailedConv, numNotEqual = runTest(lpbs, dnfs, *numberLoops, *verify && (i == *repeat-1), converter)
		if bestSoFarSucc < 0 || avgSucc < bestSoFarSucc {
			bestSoFarSucc = avgSucc
		}
		if bestSoFarAll < 0 || avgAll < bestSoFarAll {
			bestSoFarAll = avgAll
		}
	}
	// print evaluation
	fmt.Printf("Ran tests %d times, showing best average of %d repeats\n\n", *numberLoops, *repeat)
	failRate := (float64(numFailedConv) / float64(len(lpbs))) * 100.0
	fmt.Printf("Conversion failed on %d of %d tests (%.2f%%)\n", numFailedConv, len(lpbs), failRate)
	if *verify {
		errorRate := (float64(numNotEqual) / float64(len(lpbs)-numFailedConv)) * 100.0
		fmt.Printf("From the times the conversion was successful the output was wrong in %d cases (%.2f%%)\n", numNotEqual, errorRate)
	}
	fmt.Println("\nRuntime results:")
	fmt.Printf("One single conversion took %s on average on all successful runs\n", time.Duration(bestSoFarSucc))
	fmt.Printf("One single conversion took %s on average on all runs (including failed ones)\n", time.Duration(bestSoFarAll))
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

func runTest(lpbs []*lpb.LPB, dnfs []br.ClauseSet, n int, verify bool, converter lpb.DNFToLPB) (avgSucc, avgAll float64, numFailedConv, numNotEqual int) {
	avgSucc = 0.0
	avgAll = 0.0
	tSucc := 0
	tAll := 0
	for num := 0; num < n; num++ {
		numFailedConv = 0
		numNotEqual = 0
		for i, phi := range dnfs {
			start := time.Now()
			// start the solver
			computedLPB, convErr := converter.Convert(phi, len(lpbs[i].Coefficients))
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
				avgSucc = iterativeAverage(tSucc, float64(dur), avgSucc)
				tSucc++
			}
			// always update avgAll
			avgAll = iterativeAverage(tAll, float64(dur), avgAll)
			tAll++
		}
	}
	return
}
