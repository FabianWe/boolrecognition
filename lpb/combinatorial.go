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
	row := len(c.Tree[col]) - 1
	node.SetRow(row)
	return row
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
	Split(symmetryTest, cut bool)

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

// dnfFinal is a type used to indacte if a dnf is false,
// true or neither.
type dnfFinal int

const (
	IsFalse dnfFinal = iota
	IsTrue
	NotFinal
)

// isFinal checks if a dnf is false, true or neither.
func isFinal(phi br.ClauseSet) dnfFinal {
	if len(phi) == 0 {
		return IsFalse
	}
	if len(phi) == 1 && len(phi[0]) == 1 {
		return IsTrue
	}
	return NotFinal
}

type MainNode struct {
	*GenericSplitNode
	MaxL  int
	Final bool
}

func NewMainNode(lp, up SplitNode, phi br.ClauseSet,
	patterns []*OccurrencePattern, context *TreeContext) *MainNode {
	return &MainNode{NewGenericSplitNode(lp, up, phi, patterns, context), -1, false}
}

func (n *MainNode) calcFinal() {
	switch isFinal(n.Phi) {
	case IsFalse, IsTrue:
		n.Final = true
	default:
		n.Final = false
	}
}

func (n *MainNode) IsFinal() bool {
	return n.Final
}

// TODO JGS: cut is being ignored because I don't really understand it :)
func (n *MainNode) Split(symmetryTest, cut bool) {
	n.MaxL = ComputeMaxL(n.GetPatterns())
	nodeVal := isFinal(n.GetPhi())
	if n.MaxL == 1 {
		if cut {
			switch nodeVal {
			case IsFalse:
				splitRes := Split(n, 1, true, symmetryTest)
				lowerChild := NewMainNode(nil, n, splitRes.Phi, splitRes.Occurrences, n.GetContext())
				lowerChild.Final = splitRes.Final
				n.SetLowerChild(lowerChild)
				return
			case IsTrue:
				splitRes := Split(n, 0, true, symmetryTest)
				upperChild := NewMainNode(n, nil, splitRes.Phi, splitRes.Occurrences, n.GetContext())
				upperChild.Final = splitRes.Final
				n.SetUpperChild(upperChild)
				return
			}
		}
		first, second := SplitBoth(n, true, symmetryTest)

		upperChild := NewMainNode(n, nil, first.Phi, first.Occurrences, n.GetContext())
		upperChild.Final = first.Final
		n.SetUpperChild(upperChild)

		lowerChild := NewMainNode(nil, n, second.Phi, second.Occurrences, n.GetContext())
		lowerChild.Final = second.Final
		n.SetLowerChild(lowerChild)
	} else {
		if cut {
			switch nodeVal {
			case IsFalse:
				splitRes := Split(n, 1, false, symmetryTest)
				lowerChild := NewAuxNode(nil, n, splitRes.Phi, splitRes.Occurrences,
					n.GetContext(), n.MaxL, 1)
				n.SetLowerChild(lowerChild)
				return
			case IsTrue:
				splitRes := Split(n, 0, false, symmetryTest)
				upperChild := NewAuxNode(n, nil, splitRes.Phi, splitRes.Occurrences,
					n.GetContext(), n.MaxL, 1)
				n.SetUpperChild(upperChild)
				return
			}
		}
		first, second := SplitBoth(n, false, symmetryTest)

		upperChild := NewAuxNode(n, nil, first.Phi, first.Occurrences,
			n.GetContext(), n.MaxL, 1)
		n.SetUpperChild(upperChild)

		lowerChild := NewAuxNode(nil, n, second.Phi, second.Occurrences,
			n.GetContext(), n.MaxL, 1)
		n.SetLowerChild(lowerChild)
	}
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

func (n *AuxNode) createMainNode() bool {
	return n.LPrime == (n.LValue - 1)
}

func (n *AuxNode) Split(symmetryTest, cut bool) {
	createBoth := false
	nodeVal := isFinal(n.GetPhi())
	if cut {
		switch nodeVal {
		case IsFalse:
			if n.createMainNode() {
				splitRes := Split(n, 1, true, symmetryTest)
				lowerChild := NewMainNode(nil, n, splitRes.Phi, splitRes.Occurrences, n.GetContext())
				lowerChild.Final = splitRes.Final
				n.SetLowerChild(lowerChild)
			} else {
				splitRes := Split(n, 1, false, symmetryTest)
				lowerChild := NewAuxNode(nil, n, splitRes.Phi, splitRes.Occurrences,
					n.GetContext(), n.LValue, n.LPrime+1)
				n.SetLowerChild(lowerChild)
			}
			return
		case IsTrue:
			if n.GetUpperParent().GetUpperChild() == nil {
				// TODO remove debug once tested thoroughly
				panic("Debug error: split aux node, upperParent.upperChild is nil!")
			}
			n.SetUpperChild(n.GetUpperParent().GetUpperChild().GetLowerChild())
			return
		}
	}
	if n.GetUpperParent() != nil {
		if n.GetUpperParent().GetUpperChild() == nil {
			// TODO remove debug once tested thoroughly
			panic("Debug error: split aux node, upperParent.upperChild is nil!")
		}
		n.SetUpperChild(n.GetUpperParent().GetUpperChild().GetLowerChild())
		n.GetUpperChild().SetLowerParent(n)
	} else {
		createBoth = true
	}
	if createBoth {
		if n.createMainNode() {
			first, second := SplitBoth(n, true, symmetryTest)

			upperChild := NewMainNode(n, nil, first.Phi, first.Occurrences, n.GetContext())
			upperChild.Final = first.Final
			n.SetUpperChild(upperChild)

			lowerChild := NewMainNode(nil, n, second.Phi, second.Occurrences, n.GetContext())
			lowerChild.Final = second.Final
			n.SetLowerChild(lowerChild)
		} else {
			first, second := SplitBoth(n, false, symmetryTest)

			upperChild := NewAuxNode(n, nil, first.Phi, first.Occurrences,
				n.GetContext(), n.LValue, n.LPrime+1)
			n.SetUpperChild(upperChild)

			lowerChild := NewAuxNode(nil, n, second.Phi, second.Occurrences,
				n.GetContext(), n.LValue, n.LPrime+1)
			n.SetLowerChild(lowerChild)
		}
	} else {
		// here is the code for the idea in the TODO below
		// start a go routine that computes Split with k = 0
		// and write it to a channel

		// ch := make(chan *SplitResult)
		// if symmetryTest {
		// 	go func() {
		// 		ch <- Split(n, 0, false, symmetryTest)
		// 	}()
		// }

		// end of the symmetry test code
		if n.createMainNode() {
			splitRes := Split(n, 1, true, symmetryTest)
			lowerChild := NewMainNode(nil, n, splitRes.Phi, splitRes.Occurrences, n.GetContext())
			// TODO in the C++ version I calculated isFinal on the node, why not do it
			// this way, is this wrong?
			lowerChild.Final = splitRes.Final
			n.SetLowerChild(lowerChild)
		} else {
			splitRes := Split(n, 1, false, symmetryTest)
			lowerChild := NewAuxNode(nil, n, splitRes.Phi, splitRes.Occurrences,
				n.GetContext(), n.LValue, n.LValue+1)
			n.SetLowerChild(lowerChild)
		}
		if symmetryTest {
			// TODO here was the symmetry test, removed it
			// have to think about how to do it best,
			// I think it's not a good idea to call split again?
			// of course we could use SplitBoth but this would needlessly create
			// occurrence patterns, even if we don't need it
			// I think it's best if we simply split before we even decide if we're
			// in a main node (in a different go routine)
			// however checking if two DNFs are equal is not that easy with this
			// style, or is it?
			// all clauses are sorted, so we simply must compare the length and
			// then iterate over each element in the slice,
			// such an implementation is already present in split_test
			// now that I think of it that should do the trick...

			// here is the code if we would want to receive the result
			// we computed concurrently

			// testDnf := <-ch

			// now test it and so on
			// end of this code snippet
		}
	}
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
// If createPatterns is true the occurrence patterns will be created
//
// TODO implement symmetry test.
func Split(n SplitNode, k int, createPatterns, symmetryTest bool) *SplitResult {
	nbvar := n.GetContext().Nbvar
	column := n.GetColumn()
	// we can update the patterns while we iterate over the dnf
	var newOccurrences []*OccurrencePattern
	if createPatterns {
		newOccurrences = EmptyPatterns(nbvar - column - 1)
	}
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
				if createPatterns {
					updateOP(newOccurrences, newClause, nbvar, column+1)
				}
			}
		}
	} else {
		// if k is one we copy all clauses that contain the variable, but we
		// remove the variable form the clauses
		for _, clause := range n.GetPhi() {
			if len(clause) == 0 {
				// empty clause! So return Split with k = 0
				return Split(n, 0, createPatterns, symmetryTest)
			}
			// if the variable is contained copy the clause and remove the variable
			// this means to simply remove the first element
			if clause[0] == variable {
				newClause := clause[1:]
				if len(newClause) == 0 {
					isResFinal = true
				}
				newDNF = append(newDNF, newClause)
				if createPatterns {
					updateOP(newOccurrences, newClause, nbvar, column+1)
				}
			}
		}
	}
	if len(newDNF) == 0 {
		isResFinal = true
	}
	// sort new occurrence patterns
	if createPatterns {
		SortAll(newOccurrences)
	}
	return NewSplitResult(isResFinal, newDNF, newOccurrences)
}

//
// TODO implement symmetry test.
func SplitBoth(n SplitNode, createPatterns, symmetryTest bool) (*SplitResult, *SplitResult) {
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
	if createPatterns {
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
	}
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
			if createPatterns {
				wg.Add(1)
				go func(c br.Clause) {
					updateChanOne <- c
				}(newClause)
			}
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
				if createPatterns {
					wg.Add(1)
					go func(c br.Clause) {
						updateChanTwo <- c
					}(newClause)
				}
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

// SplittingTree represents the tree for a DNF.
type SplittingTree struct {
	Root                      *MainNode    // The root node
	Context                   *TreeContext // The context of the tree
	Renaming, ReverseRenaming []int        // See NewSplittingTree
	DeleteContent             bool         // If set to true the content of a node (dnf and OPs) are set to nil and deleted by garbage collection
	Cut                       bool         // TODO JGS Not entirely sure what this is supposed to mean
}

// NewSplittingTree creates a new tree given the DNF ϕ.
// Important note: For our algorithm to work the variables must be sorted
// according to their importance. Since this is not always the case (only
// during testing and some very special cases) this method will do this for
// you, i.e. it will create the occurrence patterns and then rename all
// variables accordingly. So the DNF we store in the root node is the
// renamed DNF. But we also store the mapping that caused this renaming
// in the tree in the field Renaming. This slice stores for each "old" variable
// the id in the new tree, i.e. a lookup tree.Renaming[id] gives you the
// id of the variable in the new tree.
// The reverse mapping, i.e. new variable → old variable is stored in
// ReverseRenaming.
//
// If you don't need the renaming set sortPatterns to false, in this case
// the patterns will work properly but the patterns don't get sorted.
// That is only set it to false if you know that the ordering of the variables
// is already correct.
//
// Also the clauses in the DNF must be sorted in increasing order.
// If you don't want the clauses to get sorted set sortClauses to false.
// Of course this only makes sense if also sortPatterns is set to false,
// otherwise the new dnf might not be sorted.
// This functions will sort them in this case nonetheless.
//
// The variables in the DNF have to be 0 <= v < nbar (so nbvar must be correct
// and variables start with 0).
// Also each variable should appear at least once in the DNF, what happens
// otherwise is not tested yet.
func NewSplittingTree(phi br.ClauseSet, nbvar int, sortPatterns, sortClauses bool) *SplittingTree {
	context := NewTreeContext(nbvar)
	// setup the patterns and the renamings
	newDNF, patterns, renaming, reverseRenaming := initOPs(phi, nbvar, sortPatterns)
	if sortPatterns || sortClauses {
		newDNF.SortAll()
	}
	// create root node
	// set final by computing it by hand
	root := NewMainNode(nil, nil, newDNF, patterns, context)
	// compute if node is already final
	root.calcFinal()
	// set column for root
	root.SetColumn(0)
	// insert root in the context
	context.AddNode(root)
	return &SplittingTree{Root: root,
		Context:         context,
		Renaming:        renaming,
		ReverseRenaming: reverseRenaming}
}

// initOPs initializes the occurrence patterns for ϕ.
// That is: It creates all patterns and sorts them.
// It will also compute Renaming and ReverseRenaming as discussed in
// NewSplittingTree.
//
// It returns first the renamedDNF, the patterns, then Renaming and then ReverseRenaming.
// If sortPatterns is false the old dnf will be returned.
//
// TODO never tested, but seems reasonable
func initOPs(phi br.ClauseSet, nbvar int, sortPatterns bool) (br.ClauseSet, []*OccurrencePattern, []int, []int) {
	newDNF := phi
	// intialize the renaming stuff
	renaming := make([]int, nbvar)
	reverseRenaming := make([]int, nbvar)
	// initialize the occurrence patterns for the DNF
	patterns := OPFromDNF(phi, nbvar)
	// sort each pattern
	SortAll(patterns)
	// sort the pattern slice only if sortPatterns is set
	if sortPatterns {
		SortPatterns(patterns)
		// now also create the mappings
		for newVariableId, pattern := range patterns {
			renaming[pattern.VariableId] = newVariableId
			reverseRenaming[newVariableId] = pattern.VariableId
		}
		// we also must rename each variable in the dnf and return the new dnf
		newDNF = make([]br.Clause, len(phi))
		// now clone each clause, we'll do that concurrently
		var wg sync.WaitGroup
		wg.Add(len(phi))
		for i := 0; i < len(phi); i++ {
			go func(index int) {
				clause := phi[index]
				var newClause br.Clause = make([]int, len(clause))
				for j, oldID := range clause {
					newClause[j] = renaming[oldID]
				}
				newDNF[index] = newClause
				wg.Done()
			}(i)
		}
		wg.Wait()
	}
	return newDNF, patterns, renaming, reverseRenaming
}
