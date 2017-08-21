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
	"errors"
	"fmt"
	"strconv"
	"strings"

	br "github.com/FabianWe/boolrecognition"
)

// LPBCoeff is the type used for integers in an LPB, i.e. for coefficients and
// the threshold.
//
// There are to reasons (I see) not to use ints directly:
// 1. We can change the type to a float or something else if needed be
// 2. We need ∞ and -∞ later, this is not directly possible with an int
// so we kind of wrap this type around ints.
// Important note: LPBCoeff is intended to be positive always. We will use
// negative values to indicate ∞ and -∞ (see constants later).
//
// If you want to print them nicely (for example with Printf) don't use %d
// as a specifier but %s, the String() method will format it correctly.
//
// Also you can still use some int methods like +, but if one of the values
// is ∞ or -∞ this will not give you what you want. Use Add(), Sub() and
// Mult() instead.
type LPBCoeff int

// If we ever need positive / negative values we should set PositiveInfinity
// to the max. integer and NegativeInfinity to the min. integer.
// See https://stackoverflow.com/questions/6878590/the-maximum-value-for-an-int-type-in-go

const (
	PositiveInfinity LPBCoeff = LPBCoeff(-1) // Value for ∞
	NegativeInfinity LPBCoeff = LPBCoeff(-2) // Value for -∞
)

func (val LPBCoeff) String() string {
	switch val {
	case NegativeInfinity:
		return "-∞"
	case PositiveInfinity:
		return "∞"
	default:
		return strconv.Itoa(int(val))
	}
}

// Add adds two LPBCoeff elements and returns the sum.
// If val1 one is ∞ or -∞ this value is returned, otherwise if it is a
// number the sum will be returned (if val2 is ∞ or -∞ this value gets
// returned as well).
func (val1 LPBCoeff) Add(val2 LPBCoeff) LPBCoeff {
	switch val1 {
	case NegativeInfinity:
		return NegativeInfinity
	case PositiveInfinity:
		return PositiveInfinity
	default:
		switch val2 {
		case NegativeInfinity:
			return NegativeInfinity
		case PositiveInfinity:
			return PositiveInfinity
		default:
			return val1 + val2
		}
	}
}

// Sub computes val1 - val2.
// If neither is ∞ or -∞ it returns val1 - val2. If val1 is ∞ or -∞ this
// value is returned.
// If val1 is a number and val2 is ∞ it returns -∞, if val2 is -∞ it returns ∞.
// It will however not test if the result is still positive! So check this first
// if you are not sure.
func (val1 LPBCoeff) Sub(val2 LPBCoeff) LPBCoeff {
	switch val1 {
	case NegativeInfinity:
		return NegativeInfinity
	case PositiveInfinity:
		return PositiveInfinity
	default:
		switch val2 {
		case NegativeInfinity:
			return PositiveInfinity
		case PositiveInfinity:
			return NegativeInfinity
		default:
			return val1 - val2
		}
	}
}

// Mult computes val1 * val2.
// If either of the value is ∞ or -∞ this value is returned.
func (val1 LPBCoeff) Mult(val2 LPBCoeff) LPBCoeff {
	switch val1 {
	case NegativeInfinity:
		return NegativeInfinity
	case PositiveInfinity:
		return PositiveInfinity
	default:
		switch val2 {
		case PositiveInfinity:
			return PositiveInfinity
		case NegativeInfinity:
			return NegativeInfinity
		default:
			return val1 * val2
		}
	}
}

// Compare compares two values and returns 0 iff val1 = val2, -1 iff
// val1 < val2 and 1 iff val1 > val2.
//
// ∞ and -∞ are handled in the following way:
// ∞ = ∞ and -∞ = -∞.
// No value is greater than ∞ and no value is smaller than -∞.
func (val1 LPBCoeff) Compare(val2 LPBCoeff) int {
	switch val1 {
	case NegativeInfinity:
		if val2 == NegativeInfinity {
			return 0
		} else {
			return -1
		}
	case PositiveInfinity:
		if val2 == PositiveInfinity {
			return 0
		} else {
			return 1
		}
	default:
		switch val2 {
		case NegativeInfinity:
			return 1
		case PositiveInfinity:
			return -1
		default:
			switch {
			case val1 == val2:
				return 0
			case val1 < val2:
				return -1
			default:
				// val1 > val2
				return 1
			}
		}
	}
}

// Equals returns true iff val1 == val2.
// This is the case iff val1.Compare(val2) == 0.
func (val1 LPBCoeff) Equals(val2 LPBCoeff) bool {
	// well that's easy...
	return val1 == val2
}

// Lesser returns true iff val1 < val2.
// This is the case iff val.Compare(val2) < 0.
func (val1 LPBCoeff) Lesser(val2 LPBCoeff) bool {
	return val1.Compare(val2) < 0
}

// Greater returns true iff val1 > val2.
// This is the case iff val1.Compare(val2) > 0.
func (val1 LPBCoeff) Greater(val2 LPBCoeff) bool {
	return val1.Compare(val2) > 0
}

// CoeffMax returns the maximum of both arguments.
func CoeffMax(val1, val2 LPBCoeff) LPBCoeff {
	if val1.Greater(val2) {
		return val1
	}
	return val2
}

// CoeffMin returns the minimum of both arguments.
func CoeffMin(val1, val2 LPBCoeff) LPBCoeff {
	if val1.Lesser(val2) {
		return val1
	}
	return val2
}

// LPB represents an LPB of the form a_1 ⋅ x_1 + ... + a_n ⋅ x_n ≥ d
// Each of the coefficients a_i must be a natural number (positive integer).
// d can be any integer.
type LPB struct {
	Threshold    LPBCoeff
	Coefficients []LPBCoeff
}

// EmptyLPB creates a new LPB with the threshold set to -1 and an empty
// coefficients list.
func EmptyLPB() *LPB {
	return &LPB{Threshold: -1, Coefficients: nil}
}

// NewLPB creates a new LPB with the given threshold and coefficients.
func NewLPB(threshold LPBCoeff, coefficients []LPBCoeff) *LPB {
	return &LPB{Threshold: threshold,
		Coefficients: coefficients}
}

// ParseLPB parses an LPB from the given string, if there is a syntax error
// it returns an error != nil.
//
// The syntax for parsing LPBs is as follows:
// First all the coefficients are separated by a space, then the threshold
// folows, also separated by a space.
//
// So the LPB 2 ⋅ x1 + 1 ⋅ x2 + 1 ⋅ x3 ≥ 2 is represented by "2 1 1 2".
func ParseLPB(str string) (*LPB, error) {
	split := strings.Split(str, " ")
	if len(split) == 0 {
		return nil, errors.New("LPB description is empty.")
	}
	threshold, err := strconv.Atoi(split[len(split)-1])
	if err != nil {
		return nil, err
	}
	coefficients := make([]LPBCoeff, len(split)-1)
	for i, strVal := range split[:len(split)-1] {
		if val, parseErr := strconv.Atoi(strVal); parseErr == nil {
			// no parse error, check if value is positive
			if val < 0 {
				return nil, fmt.Errorf("LPB coefficients must be positive, got %d", val)
			}
			// it's okay, append it to the coefficients
			coefficients[i] = LPBCoeff(val)
		} else {
			// parsing error
			return nil, parseErr
		}
	}
	// everything done, return result
	return NewLPB(LPBCoeff(threshold), coefficients), nil
}

func (lpb *LPB) String() string {
	buffer := new(bytes.Buffer)
	switch len(lpb.Coefficients) {
	case 0:
		buffer.WriteRune('0')
	default:
		fmt.Fprintf(buffer, "%s⋅x1", lpb.Coefficients[0])
		for i, c := range lpb.Coefficients[1:] {
			fmt.Fprintf(buffer, " + %s⋅x%d", c, i+2)
		}
	}
	fmt.Fprintf(buffer, " ≥ %s", lpb.Threshold)
	return buffer.String()
}

// Equals checks if to LPBs are syntactically equal.
func (lpb *LPB) Equals(other *LPB) bool {
	if !lpb.Threshold.Equals(other.Threshold) {
		return false
	}
	// compare the coefficients
	if lpb.Coefficients == nil && other.Coefficients == nil {
		return true
	}
	if lpb.Coefficients == nil || other.Coefficients == nil {
		return false
	}
	if len(lpb.Coefficients) != len(other.Coefficients) {
		return false
	}
	for i, val := range lpb.Coefficients {
		if !val.Equals(other.Coefficients[i]) {
			return false
		}
	}
	return true
}

func (lpb *LPB) ToDNF() br.ClauseSet {
	res := br.NewClauseSet(10)
	var sum LPBCoeff = 0
	for _, coeff := range lpb.Coefficients {
		sum = sum.Add(coeff)
	}
	// check if it represents false
	if sum.Lesser(lpb.Threshold) {
		// res is empty, so that's fine
		return res
	}
	// check if it represents true
	if lpb.Threshold.Compare(0) <= 0 {
		// add an empty clause
		c := br.NewClause(0)
		res = append(res, c)
		return res
	}
	return res
}
