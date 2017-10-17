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

// Just for playing around a bit and testing stuff.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	br "github.com/FabianWe/boolrecognition"
	"github.com/FabianWe/boolrecognition/lpb"
)

func main() {
	wenzelmannDNF := readDNFFile("wenzelmann.dnf")
	lp := lpb.NewLinearProgram(wenzelmannDNF, 5, true, true)
	computed, err := lp.Solve(lpb.TightenNone, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(computed)
}

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
