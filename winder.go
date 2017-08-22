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

package boolrecognition

import "sort"

// Class representing the Winder Matrix of a Boolean function in DNF
// representation.
//
// For a positive Boolean function f := {x1, ..., xn} → {0, 1} and a DNF ϕ
// consisting of the prime implicants of f (so ϕ represents f) the Winder matrix
// of f is defined as the n × n matrix R s.t. r[i][d] is the number of prime
// implicants of f that contain the variable xi and contain d variables.
// (Definition from "Boolean Functions: Theory, Algorithms, and Applications"
// by Yves Crama and Peter L. Hammer)
//
// For example consider the DNF ϕ ≡ { {x1, x2}, {x1, x3}, {x1, x4, x5}, {x2, x3, x4} }.
// The Winder matrix of ϕ (or better to say the function f it represents) is
//
// 0 2 1 0 0,
// 0 1 1 0 0,
// 0 1 1 0 0,
// 0 0 2 0 0,
// 0 0 1 0 0
//
// For example the first row counts all occurrences of x1. So the two 2 comes
// from the fact that x1 occurs in two clauses of size 2.
//
// The last two columns contain 0 everywhere because there is no clause of size
// 4 or 5.
//
// The matrix as we compute it also contains an addition column that saves the
// id of the variable the corresponding row was created for, so we actually
// create the following matrix:
//
// 0 2 1 0 0 | 1,
// 0 1 1 0 0 | 2,
// 0 1 1 0 0 | 3,
// 0 0 2 0 0 | 4,
// 0 0 1 0 0 | 5
//
// We will also assume that ϕ is zero based, meaning that
type WinderMatrix [][]int

// NewWinderMatrix returns a new winder matrix with the given size.
// That is it returns a matrix of size nbvar × nbvar + 1 (because in each colum)
// we also save the variable id.
//
// The matrix gets initialized with the correct values if create = true.
// Otherwise the matrix will be initalized with 0 everywhere (excet the last
// row that becomes the variable id).
func NewWinderMatrix(phi ClauseSet, nbvar int, create bool) WinderMatrix {
	var res WinderMatrix = make([][]int, nbvar, nbvar)
	for i := 0; i < nbvar; i++ {
		res[i] = make([]int, nbvar+1)
		res[i][nbvar] = i
	}
	if create {
		res.create(phi)
	}
	return res
}

// create initializes the matrix with the DNF ϕ, that is it sets up the
// correct occurrences in the matrix.
func (matrix WinderMatrix) create(phi ClauseSet) {
	for _, clause := range phi {
		length := len(clause)
		for _, v := range clause {
			matrix[v][length-1]++
		}
	}
}

func (matrix WinderMatrix) Equals(other WinderMatrix) bool {
	if len(matrix) != len(other) {
		return false
	}
	for i, row := range matrix {
		otherRow := other[i]
		if len(row) != len(otherRow) {
			return false
		}
		for j, entry := range row {
			if entry != otherRow[j] {
				return false
			}
		}
	}
	return true
}

// CompareMatrixEntry returns two rows in a Winder matrix.
//
// That is needed in order to sort the matrix later on.
// It returns -1 if row1 ≺ row2, 1 if row1 ≻ row2 and 0 if they're equal.
//
// row1 ≺ row2 iff for the first entry where row1[i] and row2[i] are diffferent
// it holds that row1[i] < row2[i].
// row1 ≻ row2 iff for the first entry where row1[i] and row2[i] are diffferent
// it holds that row[i] > row2[i].
func CompareMatrixEntry(row1, row2 []int) int {
	size := len(row1) - 1
	for i := 0; i < size; i++ {
		val1, val2 := row1[i], row2[i]
		switch {
		case val1 < val2:
			return -1
		case val1 > val2:
			return 1
		}
	}
	return 0
}

// Sort sorts the whole matrix according to the ≻ order as defined in
// CompareMatrixEntry.
// That is the most important variable comes first etc.
//
// This function is implemented with the sort package, however I read in
// Boolean Functions: Theory, Algorithms, and Applications by Yves Crama and
// Peter L. Hammer that the runtime can be improved here.
func (matrix WinderMatrix) Sort() {
	cmp := func(i, j int) bool {
		return CompareMatrixEntry(matrix[i], matrix[j]) > 0
	}
	sort.Slice(matrix, cmp)
}
