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

import "testing"
import br "github.com/FabianWe/boolrecognition"

func TestBinSearch(t *testing.T) {
	s1 := make([]int, 0)
	s2 := []int{42}
	s3 := []int{0, 1, 2, 3, 21, 42, 84, 666}
	tests := []struct {
		s        []int
		val      int
		expected int
	}{
		{s1, 0, -1},
		{s1, 42, -1},
		{s2, 84, -1},
		{s2, 42, 0},
		{s3, 0, 0},
		{s3, 1, 1},
		{s3, 2, 2},
		{s3, 3, 3},
		{s3, 21, 4},
		{s3, 42, 5},
		{s3, 84, 6},
		{s3, 666, 7},
	}
	for _, tt := range tests {
		actual := br.BinSearch(tt.s, tt.val)
		if actual != tt.expected {
			t.Errorf("Searching for %d in %v, expected %d and got %d", tt.val, tt.s, tt.expected, actual)
		}
	}
}
