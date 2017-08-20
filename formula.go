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

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"sync"

	"github.com/FabianWe/dimacscnf"
)

// Clause is a clause either in a CNF or DNF. It just stores all literals
// in a slice.
//
// How the elements are stored in a clause is up to the implementation and
// problem specific. For example you might say that variables are identified
// by integers ≠ 0 and a positive occurrence of a variable v is expressed by
// storing v whereas ¬v is expressed by storing -a.
//
// Or if you're dealing with only positive occurrences you might start your
// enumeration with 0.
//
// It is however important to document this properly.
type Clause []int

// NewClause returns an empty clause but with the capacity of the underlying
// slice big enough to hold initialCapacity variables.
func NewClause(initialCapacity int) Clause {
	return make([]int, 0, initialCapacity)
}

// Sort sorts the variables in the clause in increasin order.
func (c Clause) Sort() {
	sort.Ints(c)
}

func (c Clause) String() string {
	buffer := new(bytes.Buffer)
	switch len(c) {
	case 0:
		buffer.WriteRune('∅')
	case 1:
		fmt.Fprintf(buffer, "{%d}", c[0])
	default:
		buffer.WriteRune('{')
		buffer.WriteString(strconv.Itoa(c[0]))
		for _, val := range c[1:] {
			fmt.Fprintf(buffer, ", %d", val)
		}
		buffer.WriteRune('}')
	}
	return buffer.String()
}

// ClauseSet is a set of clauses, so a DNF or CNF (or whatever you have in
// mind...).
type ClauseSet []Clause

// NewClauseSet returns an empty clause set but with the capacity of the
// underlying slice big enough to hold initialCapacity clauses.
func NewClauseSet(initialCapacity int) ClauseSet {
	return make([]Clause, 0, initialCapacity)
}

func (phi ClauseSet) String() string {
	buffer := new(bytes.Buffer)
	switch len(phi) {
	case 0:
		buffer.WriteRune('∅')
	case 1:
		fmt.Fprintf(buffer, "{ %s }", phi[0])
	default:
		fmt.Fprintf(buffer, "{ %s", phi[0])
		for _, clause := range phi[1:] {
			fmt.Fprintf(buffer, ", %s", clause)
		}
		buffer.WriteString(" }")
	}
	return buffer.String()
}

// SortAll will sort all clauses in increasing order.
func (phi ClauseSet) SortAll() {
	var wg sync.WaitGroup
	wg.Add(len(phi))
	for _, clause := range phi {
		go func(c Clause) {
			c.Sort()
			wg.Done()
		}(clause)
	}
}

// positiveDimacsParser is a type that implements dimacscnf.DimacsParserHandler
// and is used in ParsePositiveDIMACS to parse the input.
type positiveDimacsParser struct {
	clauses   ClauseSet
	problem   string
	nbvar     int
	nbclauses int
}

func (h *positiveDimacsParser) ProblemLine(problem string, nbvar, nbclauses int) error {
	h.clauses = NewClauseSet(nbclauses)
	h.problem = problem
	h.nbvar = nbvar
	return nil
}

func (h *positiveDimacsParser) NewClause() error {
	h.clauses = append(h.clauses, NewClause(10))
	return nil
}

func (h *positiveDimacsParser) NewVariable(value int) error {
	if len(h.clauses) == 0 {
		return errors.New("Trying to append a variable but no clause was created, probably a parser bug.")
	}
	// check if value is positive, only positive values are allowed
	if value <= 0 {
		return fmt.Errorf("Illegal variable %d: Must be positive", value)
	}
	// check if variable is allowed, i.e. <= then nbvar
	if value > h.nbvar {
		return fmt.Errorf("nbvar was set to %d, but found variable %d", h.nbvar, value)
	}
	i := len(h.clauses) - 1
	// add value - 1, we identify all variables with ints starting with 0
	h.clauses[i] = append(h.clauses[i], value-1)
	return nil
}

func (h *positiveDimacsParser) Done() error {
	return nil
}

// ParsePositiveDIMACS parses a clause set from the reader r.
// It returns the 'name' of the problem (the stuff right after the p in the
// problem line), the number of variables (nbvar) and the ClauseSet itself.
//
// It only parses positive variables, i.e. negative occurrences of variables
// are not allowed and result in an error.
//
// Variables are represented starting with 0, i.e. if you have the clause
// "c 1 4 7" in your DIMACS file this clause will be represented as {0, 3, 6}.
//
// For more information on the DIMACS format see http://www.satcompetition.org/2009/format-benchmarks2009.html
func ParsePositiveDIMACS(r io.Reader) (string, int, ClauseSet, error) {
	h := &positiveDimacsParser{}
	if err := dimacscnf.ParseGenericDimacs(h, r); err != nil {
		return "", -1, nil, err
	}
	return h.problem, h.nbvar, h.clauses, nil
}
