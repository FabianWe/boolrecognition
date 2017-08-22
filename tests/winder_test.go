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

	br "github.com/FabianWe/boolrecognition"
)

func TestWinder(t *testing.T) {
	var phi br.ClauseSet = []br.Clause{
		[]int{0, 1},
		[]int{0, 2},
		[]int{0, 3, 4},
		[]int{1, 2, 3},
	}
	winder := br.NewWinderMatrix(phi, 5, true)
	var expected br.WinderMatrix = [][]int{
		[]int{0, 2, 1, 0, 0, 0},
		[]int{0, 1, 1, 0, 0, 1},
		[]int{0, 1, 1, 0, 0, 2},
		[]int{0, 0, 2, 0, 0, 3},
		[]int{0, 0, 1, 0, 0, 4},
	}
	if !winder.Equals(expected) {
		t.Errorf("Expected winder matrix %s, but got %s", expected, winder)
	}
}
