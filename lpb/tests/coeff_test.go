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
	"testing"

	"github.com/FabianWe/boolrecognition/lpb"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		val1, val2, expected lpb.LPBCoeff
	}{
		{0, 1, 1},
		{21, 21, 42},
		{lpb.PositiveInfinity, 42, lpb.PositiveInfinity},
		{lpb.NegativeInfinity, 42, lpb.NegativeInfinity},
		{42, lpb.PositiveInfinity, lpb.PositiveInfinity},
		{42, lpb.NegativeInfinity, lpb.NegativeInfinity},
	}
	for _, tt := range tests {
		actual := tt.val1.Add(tt.val2)
		if actual != tt.expected {
			t.Errorf("Expected that %s + %s = %s, but got %s", tt.val1, tt.val2, tt.expected, actual)
		}
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		val1, val2, expected lpb.LPBCoeff
	}{
		{42, 21, 21},
		{5, 5, 0},
		{42, lpb.PositiveInfinity, lpb.NegativeInfinity},
		{42, lpb.NegativeInfinity, lpb.PositiveInfinity},
		{lpb.PositiveInfinity, 42, lpb.PositiveInfinity},
		{lpb.NegativeInfinity, 42, lpb.NegativeInfinity},
	}
	for _, tt := range tests {
		actual := tt.val1.Sub(tt.val2)
		if actual != tt.expected {
			t.Errorf("Expected that %s - %s = %s, but got %s", tt.val1, tt.val2, tt.expected, actual)
		}
	}
}

func MultTest(t *testing.T) {
	tests := []struct {
		val1, val2, expected lpb.LPBCoeff
	}{
		{21, 2, 42},
		{3, 5, 15},
		{42, lpb.PositiveInfinity, lpb.PositiveInfinity},
		{42, lpb.NegativeInfinity, lpb.NegativeInfinity},
		{lpb.PositiveInfinity, 42, lpb.PositiveInfinity},
		{lpb.NegativeInfinity, 42, lpb.NegativeInfinity},
	}
	for _, tt := range tests {
		actual := tt.val1.Mult(tt.val2)
		if actual != tt.expected {
			t.Errorf("Expected that %s * %s = %s, but got %s", tt.val1, tt.val2, tt.expected, actual)
		}
	}
}

func CompareTest(t *testing.T) {
	tests := []struct {
		val1, val2 lpb.LPBCoeff
		expected   int
	}{
		{42, 21, 1},
		{21, 42, -1},
		{42, 42, 0},
		{lpb.PositiveInfinity, lpb.PositiveInfinity, 0},
		{lpb.NegativeInfinity, lpb.NegativeInfinity, 0},
		{lpb.PositiveInfinity, lpb.NegativeInfinity, 1},
		{lpb.NegativeInfinity, lpb.PositiveInfinity, -1},
		{lpb.PositiveInfinity, 42, 1},
		{lpb.NegativeInfinity, 42, -1},
	}
	for _, tt := range tests {
		actual := tt.val1.Compare(tt.val2)
		if actual != tt.expected {
			t.Errorf("Expected that %s compared to %s = %s, but got %s", tt.val1, tt.val2, tt.expected, actual)
		}
	}
}

func EqualsTest(t *testing.T) {
	tests := []struct {
		val1, val2 lpb.LPBCoeff
		expected   bool
	}{
		{42, 42, true},
		{42, 21, false},
		{21, 42, false},
		{42, lpb.PositiveInfinity, false},
		{42, lpb.NegativeInfinity, false},
		{lpb.PositiveInfinity, 42, false},
		{lpb.NegativeInfinity, 42, false},
		{lpb.PositiveInfinity, lpb.PositiveInfinity, true},
		{lpb.NegativeInfinity, lpb.NegativeInfinity, true},
		{lpb.PositiveInfinity, lpb.NegativeInfinity, false},
		{lpb.NegativeInfinity, lpb.PositiveInfinity, false},
	}
	for _, tt := range tests {
		actual := tt.val1.Equals(tt.val2)
		if actual != tt.expected {
			t.Errorf("Expected that %s == %s is %s, but got %s", tt.val1, tt.val2, tt.expected, actual)
		}
	}
}
