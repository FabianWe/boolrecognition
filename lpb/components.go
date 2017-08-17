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
)

// LPB represents an LPB of the form a_1 ⋅ x_1 + ... + a_n ⋅ x_n ≥ d
// Each of the coefficients a_i must be a natural number (positive integer).
// d can be any integer.
type LPB struct {
	Threshold    int
	Coefficients []int
}

// EmptyLPB creates a new LPB with the threshold set to -1 and an empty
// coefficients list.
func EmptyLPB() *LPB {
	return &LPB{Threshold: -1, Coefficients: nil}
}

// NewLPB creates a new LPB with the given threshold and coefficients.
func NewLPB(threshold int, coefficients []int) *LPB {
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
	coefficients := make([]int, len(split)-1)
	for i, strVal := range split[:len(split)-1] {
		if val, parseErr := strconv.Atoi(strVal); parseErr == nil {
			// no parse error, check if value is positive
			if val < 0 {
				return nil, fmt.Errorf("LPB coefficients must be positive, got %d", val)
			}
			// it's okay, append it to the coefficients
			coefficients[i] = val
		} else {
			// parsing error
			return nil, parseErr
		}
	}
	// everything done, return result
	return NewLPB(threshold, coefficients), nil
}

func (lpb *LPB) String() string {
	buffer := new(bytes.Buffer)
	switch len(lpb.Coefficients) {
	case 0:
		buffer.WriteRune('0')
	default:
		fmt.Fprintf(buffer, "%d⋅x1", lpb.Coefficients[0])
		for i, c := range lpb.Coefficients[1:] {
			fmt.Fprintf(buffer, " + %d⋅x%d", c, i+2)
		}
	}
	fmt.Fprintf(buffer, " ≥ %d", lpb.Threshold)
	return buffer.String()
}
