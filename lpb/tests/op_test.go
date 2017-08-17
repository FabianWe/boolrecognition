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

package test

import "testing"
import "github.com/FabianWe/boolrecognition/lpb"

// TestOPSmaus tests Example 6.4 of Smaus' paper.
func TestOPSmaus(t *testing.T) {
	op1 := &lpb.OccurrencePattern{Occurrences: []int{2, 2, 2, 2}, Variable: 1}
	op2 := &lpb.OccurrencePattern{Occurrences: []int{2, 2, 2}, Variable: 2}
	op3 := &lpb.OccurrencePattern{Occurrences: []int{2, 2, 3}, Variable: 3}
	op4 := &lpb.OccurrencePattern{Occurrences: []int{2, 3}, Variable: 4}

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
			t.Errorf("Error comparing OP%d and OP%d: expected %d, got %d",
				tt.first.Variable, tt.second.Variable, tt.expected, actual)
		}
	}
}

// TestOPWenzelmann tests the examples from Example 2.2 from
// Wenzelmanns bachelor thesis.
func TestOPWenzelmann(t *testing.T) {
	op1 := &lpb.OccurrencePattern{Occurrences: []int{2, 2, 3}, Variable: 1}
	op2 := &lpb.OccurrencePattern{Occurrences: []int{2, 3}, Variable: 2}
	op3 := &lpb.OccurrencePattern{Occurrences: []int{2, 3}, Variable: 3}
	op4 := &lpb.OccurrencePattern{Occurrences: []int{3, 3}, Variable: 4}
	op5 := &lpb.OccurrencePattern{Occurrences: []int{3}, Variable: 5}

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
			t.Errorf("Error comparing OP%d and OP%d: expected %d, got %d",
				tt.first.Variable, tt.second.Variable, tt.expected, actual)
		}
	}
}
