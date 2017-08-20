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

package lpb

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"sync"

	br "github.com/FabianWe/boolrecognition"
)

// OccurrencePattern is a multiset of sorted integer values.
// For a DNF ϕ and a variable x in ϕ the occurrence pattern is defined
// as the multiset having one occurrence of n for each clause of length n
// in ϕ that contains x.
// (Definition 6.3 in On Boolean Functions Encodable as a Single Linear
// Pseudo-Boolean Constraint by Jan-Georg Smaus)
// They're implemented as a slice of integer values, we also store the
// id of the variable x.
//
// Occurrence patterns are sorted multisets, we first construct the whole
// pattern and then sort it. So the Insert method does not (re)sort the OP.
//
// Note that occurrence patterns are not safe for concurrent use (don't insert).
// from multiple goroutines.
type OccurrencePattern struct {
	// The actual pattern
	Occurrences []int
	// VariableId is the id of the variable, this is only needed in the first run!
	// When we create OPs for DNFs in the second or later column this value will
	// not make sense! It only makes sense in the first column of the tree.
	VariableId int
}

// EmptyOccurrencePattern greates a new occurrence pattern for the variable
// and the slice of the numbers is initialized with the specified capacity.
func EmptyOccurrencePattern(initialCapacity int) *OccurrencePattern {
	return &OccurrencePattern{Occurrences: make([]int, 0, initialCapacity), VariableId: -1}
}

func (op *OccurrencePattern) String() string {
	buffer := new(bytes.Buffer)
	switch len(op.Occurrences) {
	case 0:
		buffer.WriteRune('∅')
	case 1:
		fmt.Fprintf(buffer, "⦃%d⦄", op.Occurrences[0])
	default:
		buffer.WriteRune('⦃')
		buffer.WriteString(strconv.Itoa(op.Occurrences[0]))
		for _, val := range op.Occurrences[1:] {
			fmt.Fprintf(buffer, ", %d", val)
		}
		buffer.WriteRune('⦄')
	}
	return buffer.String()
}

// Insert adds the value to the pattern. It will not (re)sort the array.
func (op *OccurrencePattern) Insert(val int) {
	op.Occurrences = append(op.Occurrences, val)
}

// Sort sorts the occurrences in increasing order.
func (op *OccurrencePattern) Sort() {
	sort.Ints(op.Occurrences)
}

// CompareTo compares the occurrence pattern to another one.
// It returns 0 iff op == other, -1 iff op ≺ other and 1 iff op ≻ other.
func (op *OccurrencePattern) CompareTo(other *OccurrencePattern) int {
	n, m := len(op.Occurrences), len(other.Occurrences)
	i, j := 0, 0
	for i < n && j < m {
		v1, v2 := op.Occurrences[i], other.Occurrences[j]
		i++
		j++
		// we've already found the longest common prefix if v1 != v2
		// in this case we return the result
		// if v1 == v2 we continue the loop
		switch {
		case v1 < v2:
			return 1
		case v2 < v1:
			return -1
		}
	}
	// the for loop is done so either both are empty or one of the has some
	// elements left
	switch {
	case i < n:
		return 1
	case j < m:
		return -1
	default:
		// in this case both conditions above failed, therefore the occurrence
		// patterns are equal
		return 0
	}
}

// EmptyPatterns returns a slice of occurrence pattern of the given size, each
// one is initialized to an empty pattern.
func EmptyPatterns(size int) []*OccurrencePattern {
	res := make([]*OccurrencePattern, size)
	for i := 0; i < size; i++ {
		newPattern := EmptyOccurrencePattern(10)
		newPattern.VariableId = i
		res[i] = newPattern
	}
	return res
}

// updateOP will update the occurrence patterns given a new clause.
func updateOP(patterns []*OccurrencePattern, clause br.Clause, nbvar, column int) {
	n := len(clause)
	for _, x := range clause {
		patterns[x-column].Insert(n)
	}
}

// OPFromDNF will build the occurrence patterns for the DNF ϕ.
// The number of variables (nbvar) must be known in advance.
// No variable in the DNF must be >= nbvar, this will not be checked though!
//
// This function returns a slice of length nbvar where the OP for variable x
// is stored on position x.
//
// This method does not sort the patterns (neither the patterns themselves
// or the whole pattern slice)! You must first make sure that all the patterns
// are sorted (the elements in each pattern), see SortAll for this.
// Then you must make sure that the elements are sorted according the ≽
// order, see SortPatterns for this.
func OPFromDNF(phi br.ClauseSet, nbvar int) []*OccurrencePattern {
	return OPFromDNFShift(phi, nbvar, 0)
}

// OPFromDNFShift will build the occurrence patterns for the DNF ϕ.
//
// This is a rather internal method, but I exported it for testing and playing
// around. If you want to build the patterns for a whole DNF (no splitting) or
// something use OPFromDNF.
//
// nbvar is the number of variables in the whole dnf, this means if you already
// used split some variables you always use the nbvar for the original DNF!
//
// column refers to the column this DNF is created for in the tree.
// So if we call Split on a node in column zero (the first one) column will
// be one because we create the successors that are stored in column one
// (the second one)
func OPFromDNFShift(phi br.ClauseSet, nbvar, column int) []*OccurrencePattern {
	res := EmptyPatterns(nbvar - column)
	// we could think if we want to do some concurrent stuff here, but we would
	// have to lock the occurrence patterns... and we don't really want that
	for _, clause := range phi {
		updateOP(res, clause, nbvar, column)
	}
	return res
}

// SortPatterns will sort the occurrence patterns according to importance of
// the variables.
// That is the greatest element according to ≽ comes first.
func SortPatterns(patterns []*OccurrencePattern) {
	comp := func(i, j int) bool {
		return patterns[i].CompareTo(patterns[j]) > 0
	}
	sort.Slice(patterns, comp)
}

// SortAll will sort each occurrence pattern in patterns.
// So don't confuse this method with SortPatterns, this will sort the patterns
// according to ≽, but each pattern itself must be sorted first with this
// method.
//
// This method will sort all patterns concurrently.
//
// TODO potential improvement? Sort the clauses according to length first.
// So when constructing new patterns we don't have to sort again and again.
// But is this always correct when we split away variables? I don't think so.
func SortAll(patterns []*OccurrencePattern) {
	var wg sync.WaitGroup
	wg.Add(len(patterns))
	for _, pattern := range patterns {
		go func(op *OccurrencePattern) {
			op.Sort()
			wg.Done()
		}(pattern)
	}
	wg.Wait()
}

// ComputeMaxL computes the max l s.t. the first l patterns are equal.
//
// There are some TODO here.
func ComputeMaxL(patterns []*OccurrencePattern) int {
	l := 1
	// TODO: Strange behavior if the DNF cannot be represented as a LPB
	// first can be null in this case?
	// why? JGS: let's see if implementing a symmetry test will remedy this.
	first := patterns[0]
	for l < len(patterns) {
		next := patterns[l]
		if next == nil {
			// TODO create new empty OP here, correct?
			next = EmptyOccurrencePattern(10)
			patterns[l] = next
		}
		if first == nil { // TODO JGS: to fix the segmentation fault
			l++
		} else {
			if first.CompareTo(next) == 0 {
				l++
			} else {
				break
			}
		}
	}
	return l
}
