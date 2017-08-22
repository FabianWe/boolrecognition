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

package tests

import (
	"os"
	"path/filepath"
	"testing"

	br "github.com/FabianWe/boolrecognition"
	"github.com/FabianWe/boolrecognition/lpb"
)

// Testcases for both tests we run, the file itself and the expected results
var smausDNF, smausZero, smausOne br.ClauseSet
var wenzelmannDNF, wenzelmannZero, wenzelmannOne br.ClauseSet

// Reads a file from the "dnfs" subdirectory, panics on error
func readDNFFile(filename string) br.ClauseSet {
	f, fErr := os.Open(filepath.Join("dnfs", filename))
	if fErr != nil {
		panic(fErr)
	}
	defer f.Close()
	_, _, phi, err := br.ParsePositiveDIMACS(f)
	if err != nil {
		panic(err)
	}
	return phi
}

func init() {
	smausDNF = readDNFFile("smaus.dnf")
	smausZero = readDNFFile("smaus_0.dnf")
	smausOne = readDNFFile("smaus_1.dnf")
	wenzelmannDNF = readDNFFile("wenzelmann.dnf")
	wenzelmannZero = readDNFFile("wenzelmann_0.dnf")
	wenzelmannOne = readDNFFile("wenzelmann_1.dnf")
}

// cmpDNFS is a rather simple function to compare to DNFs.
// It simply checks if all clauses are equal (must have the same order).
func cmpDNFS(phi1, phi2 br.ClauseSet) bool {
	if len(phi1) != len(phi2) {
		return false
	}
	for i, c1 := range phi1 {
		c2 := phi2[i]
		if len(c1) != len(c2) {
			return false
		}
		for j, var1 := range c1 {
			if var1 != c2[j] {
				return false
			}
		}
	}
	return true
}

// TestSplitSmaus tests the splitting for example 6.6 in the paper
// from Smaus
func TestSplitSmaus(t *testing.T) {
	// create a dummy main node to call split on
	// create some dummy context as well
	ctx := lpb.NewTreeContext(5)
	n := lpb.NewMainNode(nil, nil, smausDNF, nil, ctx)
	n.SetColumn(0)
	zeroSplit := lpb.Split(n, 0, true, false)
	oneSplit := lpb.Split(n, 1, true, true)
	if zeroSplit.Final {
		t.Error("split with k = 0 must produce a non-final DNF, got final dnf")
	}
	if !cmpDNFS(zeroSplit.Phi, smausZero) {
		t.Errorf("For split with k = 0: expected %s, got %s", smausZero, zeroSplit.Phi)
	}

	if oneSplit.Final {
		t.Error("split with k = 1 must produce a non-final DNF, got final dnf")
	}
	if !cmpDNFS(oneSplit.Phi, smausOne) {
		t.Errorf("For split with k = 1: expected %s, got %s", smausOne, oneSplit.Phi)
	}
}

// TestSplitWenzelmann tests the splitting for example 2.2 in Wenzelmanns'
// bachelor thesis.
func TestSplitWenzelmann(t *testing.T) {
	// create a dummy main node to call split on
	// create some dummy context as well
	ctx := lpb.NewTreeContext(5)
	n := lpb.NewMainNode(nil, nil, wenzelmannDNF, nil, ctx)
	n.SetColumn(0)
	zeroSplit := lpb.Split(n, 0, true, false)
	oneSplit := lpb.Split(n, 1, true, true)
	if zeroSplit.Final {
		t.Error("split with k = 0 must produce a non-final DNF, got final dnf")
	}
	if !cmpDNFS(zeroSplit.Phi, wenzelmannZero) {
		t.Errorf("For split with k = 0: expected %s, got %s", wenzelmannZero, zeroSplit.Phi)
	}

	if oneSplit.Final {
		t.Error("split with k = 1 must produce a non-final DNF, got final dnf")
	}
	if !cmpDNFS(oneSplit.Phi, wenzelmannOne) {
		t.Errorf("For split with k = 1: expected %s, got %s", wenzelmannOne, oneSplit.Phi)
	}
}

// Test the SplitBoth method for the example from Smaus.
func TestSplitBothSmaus(t *testing.T) {
	// create a dummy main node to call split on
	// create some dummy context as well
	ctx := lpb.NewTreeContext(5)
	n := lpb.NewMainNode(nil, nil, smausDNF, nil, ctx)
	n.SetColumn(0)
	zeroSplit, oneSplit := lpb.SplitBoth(n, true, false)
	if zeroSplit.Final {
		t.Error("split with k = 0 must produce a non-final DNF, got final dnf")
	}
	if !cmpDNFS(zeroSplit.Phi, smausZero) {
		t.Errorf("For split with k = 0: expected %s, got %s", smausZero, zeroSplit.Phi)
	}

	if oneSplit.Final {
		t.Error("split with k = 1 must produce a non-final DNF, got final dnf")
	}
	if !cmpDNFS(oneSplit.Phi, smausOne) {
		t.Errorf("For split with k = 1: expected %s, got %s", smausOne, oneSplit.Phi)
	}
}

// Test the SplitBoth method for the example from Wenzelmann.
func TestSplitBothWenzelmann(t *testing.T) {
	// create a dummy main node to call split on
	// create some dummy context as well
	ctx := lpb.NewTreeContext(5)
	n := lpb.NewMainNode(nil, nil, wenzelmannDNF, nil, ctx)
	n.SetColumn(0)
	zeroSplit, oneSplit := lpb.SplitBoth(n, true, false)
	if zeroSplit.Final {
		t.Error("split with k = 0 must produce a non-final DNF, got final dnf")
	}
	if !cmpDNFS(zeroSplit.Phi, wenzelmannZero) {
		t.Errorf("For split with k = 0: expected %s, got %s", smausZero, zeroSplit.Phi)
	}

	if oneSplit.Final {
		t.Error("split with k = 1 must produce a non-final DNF, got final dnf")
	}
	if !cmpDNFS(oneSplit.Phi, wenzelmannOne) {
		t.Errorf("For split with k = 1: expected %s, got %s", smausOne, oneSplit.Phi)
	}
}

func TestSmausMin(t *testing.T) {
	expected := lpb.NewLPB(5, []lpb.LPBCoeff{4, 3, 2, 2, 1})
	solver := lpb.NewMinSolver()
	tree := lpb.NewSplittingTree(smausDNF, 5, true, true)
	lpb, err := solver.Solve(tree)
	if err != nil {
		t.Error("Expected LPB, got an error:", err)
	}
	if !lpb.Equals(expected) {
		t.Errorf("Expected LPB %s and got %s", expected, lpb)
	}
}

func TestWenzelmannMin(t *testing.T) {
	expected := lpb.NewLPB(8, []lpb.LPBCoeff{5, 3, 3, 2, 1})
	solver := lpb.NewMinSolver()
	tree := lpb.NewSplittingTree(wenzelmannDNF, 5, true, true)
	lpb, err := solver.Solve(tree)
	if err != nil {
		t.Error("Expected LPB, got an error:", err)
	}
	if !lpb.Equals(expected) {
		t.Errorf("Expected LPB %s and got %s", expected, lpb)
	}
}

// TODO test true and false!
