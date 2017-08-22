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

// Example 1.28 / 1.29 from BOOLEAN FUNCTIONS Theory, Algorithms, and
// Applications by Yves Crama and Peter L. Hammer (page 58)
var phi br.ClauseSet = []br.Clause{
	[]int{0, 1},
	[]int{0, 2, 3},
	[]int{1, 2},
}

// TODO write the actual tests, but this is really annoying with the concurrent
// stuff, we have to ignore the order etc.
func TestMinMaxPoints(t *testing.T) {
	// mtps := lpb.ComputeMTPs(phi, 4)
	// mfps := lpb.ComputeMFPs(mtps, true)
	// tests here
}
