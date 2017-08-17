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
// Note that occurrence patterns are not safe for concurrent use (don't insert)
// from multiple goroutines.
type OccurrencePattern struct {
	Occurrences []int
	Variable    int
}

// EmptyOccurrencePattern greates a new occurrence pattern for the variable
// and the slice of the numbers is initialized with the specified capacity.
func EmptyOccurrencePattern(variable, initialCapacity int) *OccurrencePattern {
	return &OccurrencePattern{Variable: variable,
		Occurrences: make([]int, 0, initialCapacity)}
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
