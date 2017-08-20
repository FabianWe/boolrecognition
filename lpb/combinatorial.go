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
	"fmt"
	"sync"

	br "github.com/FabianWe/boolrecognition"
)

type TreeContext struct {
	Tree  [][]SplitNode
	Nbvar int
}

func NewTreeContext(nbvar int) *TreeContext {
	return &TreeContext{Tree: make([][]SplitNode, nbvar),
		Nbvar: nbvar}
}

func (c *TreeContext) AddNode(node SplitNode) int {
	col := node.GetColumn()
	c.Tree[col] = append(c.Tree[col], node)
	return len(c.Tree[col]) - 1
}

type GenericSplitNode struct {
	LowerParent, UpperParent, LowerChild, UpperChild SplitNode
	Phi                                              br.ClauseSet
	Context                                          *TreeContext
	Column, Row                                      int
	Patterns                                         []*OccurrencePattern
	AlreadySplit                                     bool
}

func NewGenericSplitNode(lp, up SplitNode, phi br.ClauseSet,
	patterns []*OccurrencePattern, context *TreeContext) *GenericSplitNode {
	// TODO move all this somewhere else, whoever creates the node!
	// var column int
	// switch {
	// case lp != nil:
	// 	column = lp.GetColumn() + 1
	// case up != nil:
	// 	column = up.GetColumn() + 1
	// default:
	// 	column = 0
	// }
	node := &GenericSplitNode{LowerParent: lp,
		UpperParent:  up,
		LowerChild:   nil,
		UpperChild:   nil,
		Phi:          phi,
		Column:       -1,
		Row:          -1,
		Context:      context,
		Patterns:     patterns,
		AlreadySplit: false}
	// TODO check if this is even possible...
	// if context == nil {
	// 	node.Row = -1
	// } else {
	// 	node.Row = context.AddNode(node)
	// }
	// TODO we need to do all this somewhere else!
	return node
}

func (n *GenericSplitNode) GetLowerParent() SplitNode {
	return n.LowerParent
}

func (n *GenericSplitNode) SetLowerParent(node SplitNode) {
	n.LowerParent = node
}

func (n *GenericSplitNode) GetUpperParent() SplitNode {
	return n.UpperParent
}

func (n *GenericSplitNode) SetUpperParent(node SplitNode) {
	n.UpperParent = node
}

func (n *GenericSplitNode) GetLowerChild() SplitNode {
	return n.LowerChild
}

func (n *GenericSplitNode) SetLowerChild(node SplitNode) {
	n.LowerChild = node
}

func (n *GenericSplitNode) GetUpperChild() SplitNode {
	return n.UpperChild
}

func (n *GenericSplitNode) SetUpperChild(node SplitNode) {
	n.UpperChild = node
}

func (n *GenericSplitNode) GetPhi() br.ClauseSet {
	return n.Phi
}

func (n *GenericSplitNode) SetPhi(phi br.ClauseSet) {
	n.Phi = phi
}

func (n *GenericSplitNode) GetColumn() int {
	return n.Column
}

func (n *GenericSplitNode) SetColumn(column int) {
	n.Column = column
}

func (n *GenericSplitNode) GetRow() int {
	return n.Row
}

func (n *GenericSplitNode) SetRow(row int) {
	n.Row = row
}

func (n *GenericSplitNode) GetContext() *TreeContext {
	return n.Context
}

func (n *GenericSplitNode) SetContext(context *TreeContext) {
	n.Context = context
}

func (n *GenericSplitNode) GetPatterns() []*OccurrencePattern {
	return n.Patterns
}

func (n *GenericSplitNode) SetPatterns(patterns []*OccurrencePattern) {
	n.Patterns = patterns
}

func (n *GenericSplitNode) IsAlreadySplit() bool {
	return n.AlreadySplit
}

func (n *GenericSplitNode) SetAlreadySplit(val bool) {
	n.AlreadySplit = val
}

// TODO remove methods we don't need!
type SplitNode interface {
	IsFinal() bool
	GetLowerParent() SplitNode
	SetLowerParent(node SplitNode)
	GetUpperParent() SplitNode
	SetUpperParent(node SplitNode)
	GetLowerChild() SplitNode
	SetLowerChild(node SplitNode)
	GetUpperChild() SplitNode
	SetUpperChild(node SplitNode)
	GetPhi() br.ClauseSet //
	SetPhi(phi br.ClauseSet)
	GetColumn() int       //
	SetColumn(column int) //
	GetRow() int
	SetRow(row int)
	GetContext() *TreeContext //
	SetContext(context *TreeContext)
	GetPatterns() []*OccurrencePattern //
	SetPatterns(patterns []*OccurrencePattern)
	IsAlreadySplit() bool
	SetAlreadySplit(val bool)
}

type MainNode struct {
	*GenericSplitNode
	MaxL  int
	Final bool
}

func NewMainNode(lp, up SplitNode, phi br.ClauseSet,
	patterns []*OccurrencePattern, context *TreeContext,
	maxL int, final bool) *MainNode {
	return &MainNode{NewGenericSplitNode(lp, up, phi, patterns, context), maxL, final}
}

func (n *MainNode) IsFinal() bool {
	return n.Final
}

type AuxNode struct {
	*GenericSplitNode
	LValue, LPrime int
}

func NewAuxNode(lp, up SplitNode, phi br.ClauseSet,
	patterns []*OccurrencePattern, context *TreeContext,
	lValue, lPrime int) *AuxNode {
	return &AuxNode{NewGenericSplitNode(lp, up, phi, patterns, context), lValue, lPrime}
}

func (n *AuxNode) IsFinal() bool {
	return false
}

type SplitResult struct {
	Final       bool
	Phi         br.ClauseSet
	Occurrences []*OccurrencePattern
}

func NewSplitResult(final bool, phi br.ClauseSet, occurrences []*OccurrencePattern) *SplitResult {
	return &SplitResult{Final: final,
		Phi:         phi,
		Occurrences: occurrences}
}

// Split will split away the next variable. The variable that must be split
// away is given by the column of the node (in column k we split away variable
// k).
//
// TODO implement symmetry test.
func Split(n SplitNode, k int, symmetryTest bool) *SplitResult {
	nbvar := n.GetContext().Nbvar
	column := n.GetColumn()
	// we can update the patterns while we iterate over the dnf
	newOccurrences := EmptyPatterns(nbvar - column - 1)
	isResFinal := false
	// maybe too big...
	newDNF := br.NewClauseSet(len(n.GetPhi()))
	if n.GetPatterns() != nil && len(n.GetPatterns()) == 0 {
		// TODO debug, can this happen?
		fmt.Println("Debug the split method! Weird case...")
	}
	// just to make clear where the variable comes from
	variable := column
	if k == 0 {
		// in this case copy all clauses that do not contain the variable
		for _, clause := range n.GetPhi() {
			// check if the clause does not contain the variable
			if len(clause) == 0 || clause[0] != variable {
				// we simply use the old clause, this makes the code a bit more
				// understandable I hope
				newClause := clause
				if len(newClause) == 0 {
					isResFinal = true
				}
				// add clause
				newDNF = append(newDNF, newClause)
				updateOP(newOccurrences, newClause, nbvar, column+1)
			}
		}
	} else {
		// if k is one we copy all clauses that contain the variable, but we
		// remove the variable form the clauses
		for _, clause := range n.GetPhi() {
			if len(clause) == 0 {
				// empty clause! So return Split with k = 0
				return Split(n, 0, symmetryTest)
			}
			// if the variable is contained copy the clause and remove the variable
			// this means to simply remove the first element
			if clause[0] == variable {
				newClause := clause[1:]
				if len(newClause) == 0 {
					isResFinal = true
				}
				newDNF = append(newDNF, newClause)
				updateOP(newOccurrences, newClause, nbvar, column+1)
			}
		}
	}
	if len(newDNF) == 0 {
		isResFinal = true
	}
	// sort new occurrence patterns
	SortAll(newOccurrences)
	return NewSplitResult(isResFinal, newDNF, newOccurrences)
}

//
// TODO implement symmetry test.
func SplitBoth(n SplitNode, symmetryTest bool) (*SplitResult, *SplitResult) {
	nbvar := n.GetContext().Nbvar
	column := n.GetColumn()
	// again we can update the occurrence patterns while iterating over the
	// dnf, this time we use some concurrency:
	// when we add to the first dnf and then encounter another clause for the
	// second dnf we can run this concurrently.
	// that's to say: each dnf might write values to a channel (below)
	// a goroutine listens on that channels and handels the update
	// we use a waitgroup that we add in both gorutines to and later wait
	// for all updates to finish
	var wg sync.WaitGroup
	updateChanOne := make(chan br.Clause)
	updateChanTwo := make(chan br.Clause)
	// defer closing the channels so we don't forget it
	defer close(updateChanOne)
	defer close(updateChanTwo)
	// initialize both occurrence patterns, we also create them concurrently
	var patternsOne, patternsTwo []*OccurrencePattern
	wg.Add(2)
	go func() {
		defer wg.Done()
		patternsOne = EmptyPatterns(nbvar - column - 1)
	}()
	go func() {
		defer wg.Done()
		patternsTwo = EmptyPatterns(nbvar - column - 1)
	}()
	wg.Wait()
	// start the go routines
	go func() {
		for clause := range updateChanOne {
			updateOP(patternsOne, clause, nbvar, column+1)
			wg.Done()
		}
	}()
	go func() {
		for clause := range updateChanTwo {
			updateOP(patternsTwo, clause, nbvar, column+1)
			wg.Done()
		}
	}()
	isFirstFinal, isSecondFinal := false, false
	containsEmptyClause := false
	firstDNF, secondDNF := br.NewClauseSet(len(n.GetPhi())), br.NewClauseSet(len(n.GetPhi()))
	// just to make clear where the variable comes from
	variable := column

	for _, clause := range n.GetPhi() {
		if len(clause) == 0 {
			containsEmptyClause = true
		}
		// now check if variable is contained or not
		if len(clause) == 0 || clause[0] != variable {
			// we simply use the old clause, this makes the code a bit more
			// understandable I hope
			newClause := clause
			if len(newClause) == 0 {
				isFirstFinal = true
			}
			// add clause
			firstDNF = append(firstDNF, newClause)
			// add occurrence pattern to the channel
			// run a go routine and add 1 to the wait group
			wg.Add(1)
			go func(c br.Clause) {
				updateChanOne <- c
			}(newClause)
		} else if clause[0] == variable {
			// if the variable is contained copy the clause and remove the variable
			// this means to simply remove the first element
			if clause[0] == variable {
				newClause := clause[1:]
				if len(newClause) == 0 {
					isSecondFinal = true
				}
				// add clause
				secondDNF = append(secondDNF, newClause)
				// add occurrence pattern to the channel
				// run a go routine and add 1 to the wait group
				wg.Add(1)
				go func(c br.Clause) {
					updateChanTwo <- c
				}(newClause)
			}
		}
	}
	if len(firstDNF) == 0 {
		isFirstFinal = true
	}
	if len(secondDNF) == 0 {
		isSecondFinal = true
	}
	// wait for the update go routines to finish
	wg.Wait()
	if containsEmptyClause {
		// return first result for both
		// sort the patterns first
		SortAll(patternsOne)
		res := NewSplitResult(isFirstFinal, firstDNF, patternsOne)
		return res, res
	}
	// if empty clause is not contained return both results
	// again, sort the patterns
	// again concurrently, why not?
	wg.Add(2)
	go func() {
		SortAll(patternsOne)
		wg.Done()
	}()
	go func() {
		SortAll(patternsTwo)
		wg.Done()
	}()
	// wait until sorting is finished
	wg.Wait()
	res1 := NewSplitResult(isFirstFinal, firstDNF, patternsOne)
	res2 := NewSplitResult(isSecondFinal, secondDNF, patternsTwo)
	return res1, res2
}
