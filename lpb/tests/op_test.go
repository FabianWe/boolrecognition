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
	"sort"
	"testing"

	"github.com/FabianWe/boolrecognition/lpb"
)

// TestOPSmaus tests Example 6.4 of Smaus' paper.
func TestOPSmaus(t *testing.T) {
	op1 := &lpb.OccurrencePattern{Occurrences: []int{2, 2, 2, 2}}
	op2 := &lpb.OccurrencePattern{Occurrences: []int{2, 2, 2}}
	op3 := &lpb.OccurrencePattern{Occurrences: []int{2, 2, 3}}
	op4 := &lpb.OccurrencePattern{Occurrences: []int{2, 3}}

	tests := []struct {
		first, second *lpb.OccurrencePattern
		expected      int
	}{
		// order as in the paper
		{op1, op2, 1},
		{op2, op3, 1},
		{op3, op4, 1},
		// each with itself
		{op1, op1, 0},
		{op2, op2, 0},
		{op3, op3, 0},
		{op4, op4, 0},
		// order as in the paper, the other way around
		{op4, op3, -1},
		{op3, op2, -1},
		{op2, op1, -1},
		// some random other tests
		{op1, op3, 1},
		{op2, op4, 1},
		{op4, op1, -1},
	}
	for _, tt := range tests {
		actual := tt.first.CompareTo(tt.second)
		if actual != tt.expected {
			t.Errorf("Error comparing OPs: expected %d, got %d",
				tt.expected, actual)
		}
	}
}

// TestOPWenzelmann tests the examples from Example 2.2 from
// Wenzelmanns bachelor thesis.
func TestOPWenzelmann(t *testing.T) {
	op1 := &lpb.OccurrencePattern{Occurrences: []int{2, 2, 3}}
	op2 := &lpb.OccurrencePattern{Occurrences: []int{2, 3}}
	op3 := &lpb.OccurrencePattern{Occurrences: []int{2, 3}}
	op4 := &lpb.OccurrencePattern{Occurrences: []int{3, 3}}
	op5 := &lpb.OccurrencePattern{Occurrences: []int{3}}

	tests := []struct {
		first, second *lpb.OccurrencePattern
		expected      int
	}{
		{op1, op2, 1},
		{op2, op3, 0},
		{op4, op3, -1},
		{op5, op4, -1},
		{op5, op1, -1},
		{op2, op4, 1},
		{op2, op2, 0},
	}
	for _, tt := range tests {
		actual := tt.first.CompareTo(tt.second)
		if actual != tt.expected {
			t.Errorf("Error comparing OPs: expected %d, got %d",
				tt.expected, actual)
		}
	}
}

const size int = 100000

var mySlice []int = make([]int, size)
var myMap map[int]struct{} = make(map[int]struct{}, size)

func init() {
	for i := 0; i < size; i++ {
		mySlice = append(mySlice, i)
		myMap[i] = struct{}{}
	}
}

func BenchmarkAddSlice(b *testing.B) {
	for n := 0; n < b.N; n++ {
		s := make([]int, size)
		for i := 0; i < size; i++ {
			s = append(s, i)
		}
	}
}

func BenchmarkAddMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		m := make(map[int]struct{}, size)
		for i := 0; i < size; i++ {
			m[i] = struct{}{}
		}
	}
}

// Just some benchmarks on how to implement DNFs, wrong place but
// they're there just if you're curious
var x bool

func containsSlice() {
	half := size / 2
	end := size + half
	for i := half; i < end; i++ {
		index := sort.SearchInts(mySlice, i)
		x = index < len(mySlice)
	}
}

func containsMap() {
	half := size / 2
	end := size + half
	for i := half; i < end; i++ {
		_, contains := myMap[i]
		x = contains
	}
}

func binSearch(s []int, val int) bool {
	l, r := 0, len(s)-1
	for l <= r {
		m := l + (r-l)/2
		nxt := s[m]
		if nxt == val {
			return true
		}
		if nxt < val {
			l = m + 1
		} else {
			r = m - 1
		}
	}
	return false
}

func containsBinSearch() {
	half := size / 2
	end := size + half
	for i := half; i < end; i++ {
		x = binSearch(mySlice, i)
	}
}

func BenchmarkLookupSlice(b *testing.B) {
	for n := 0; n < b.N; n++ {
		containsSlice()
	}
}

func BenchmarkLookupMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		containsMap()
	}
}

func BenchmarkBinSearch(b *testing.B) {
	for n := 0; n < b.N; n++ {
		containsBinSearch()
	}
}
