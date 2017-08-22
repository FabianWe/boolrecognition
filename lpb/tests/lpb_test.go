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

func TestConvertFalse(t *testing.T) {
	falseLPB := lpb.NewLPB(10, []lpb.LPBCoeff{1, 1, 1, 3})
	res := falseLPB.ToDNF()
	if len(res) != 0 {
		t.Errorf("LPB %s should be false, but got DNF %s", falseLPB, res)
	}
}

func TestConvertTrue(t *testing.T) {
	trueLPB := lpb.NewLPB(0, []lpb.LPBCoeff{})
	res := trueLPB.ToDNF()
	if !(len(res) == 1 && len(res[0]) == 0) {
		t.Errorf("LPB %s should be true, but got DNF %s", trueLPB, res)
	}
}
