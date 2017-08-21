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
	"errors"
	"fmt"
	"sync"

	br "github.com/FabianWe/boolrecognition"
)

type TreeContext struct {
	Tree  [][]SplitNode
	Nbvar int
}

func NewTreeContext(nbvar int) *TreeContext {
	return &TreeContext{Tree: make([][]SplitNode, nbvar+1),
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
	Split(symmetryTest, cut bool) error

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

// RegisterNode will register the node in the context.
// This method must be called each time a node gets created, so probably
// in a New... method.
// It will determine the column and row of the node and set the values
// on the node.
func RegisterNode(n SplitNode) {
	column := 0
	switch {
	case n.GetLowerParent() != nil:
		column = n.GetLowerParent().GetColumn() + 1
	case n.GetUpperParent() != nil:
		column = n.GetUpperParent().GetColumn() + 1
	}
	n.SetColumn(column)
	n.GetContext().AddNode(n)
	// fmt.Printf("Added %s: %s in col %d, row %d\n", reflect.TypeOf(n), n.GetPhi(), n.GetColumn(), n.GetRow())
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
	if len(phi) == 1 && len(phi[0]) == 0 {
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
	n := &MainNode{NewGenericSplitNode(lp, up, phi, patterns, context), -1, false}
	RegisterNode(n)
	return n
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
func (n *MainNode) Split(symmetryTest, cut bool) error {
	n.SetAlreadySplit(true)
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
				return nil
			case IsTrue:
				splitRes := Split(n, 0, true, symmetryTest)
				upperChild := NewMainNode(n, nil, splitRes.Phi, splitRes.Occurrences, n.GetContext())
				upperChild.Final = splitRes.Final
				n.SetUpperChild(upperChild)
				return nil
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
				return nil
			case IsTrue:
				splitRes := Split(n, 0, false, symmetryTest)
				upperChild := NewAuxNode(n, nil, splitRes.Phi, splitRes.Occurrences,
					n.GetContext(), n.MaxL, 1)
				n.SetUpperChild(upperChild)
				return nil
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
	return nil
}

type AuxNode struct {
	*GenericSplitNode
	LValue, LPrime int
}

func NewAuxNode(lp, up SplitNode, phi br.ClauseSet,
	patterns []*OccurrencePattern, context *TreeContext,
	lValue, lPrime int) *AuxNode {
	n := &AuxNode{NewGenericSplitNode(lp, up, phi, patterns, context), lValue, lPrime}
	RegisterNode(n)
	return n
}

func (n *AuxNode) IsFinal() bool {
	return false
}

func (n *AuxNode) createMainNode() bool {
	return n.LPrime == (n.LValue - 1)
}

func (n *AuxNode) Split(symmetryTest, cut bool) error {
	n.SetAlreadySplit(true)
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
			return nil
		case IsTrue:
			if n.GetUpperParent().GetUpperChild() == nil {
				// TODO remove debug once tested thoroughly
				panic("Debug error: split aux node, upperParent.upperChild is nil!")
			}
			n.SetUpperChild(n.GetUpperParent().GetUpperChild().GetLowerChild())
			return nil
		}
	}
	if n.GetUpperParent() != nil {
		if n.GetUpperParent().GetUpperChild() == nil {
			// TODO remove debug once tested thoroughly
			panic("Debug error: split aux node, upperParent.upperChild is nil!")
		}
		n.SetUpperChild(n.GetUpperParent().GetUpperChild().GetLowerChild())
		if n.GetUpperChild() == nil {
			fmt.Printf("Problem in col %d, row %d, dnf %s\n", n.GetColumn(), n.GetRow(), n.GetPhi())
			fmt.Println("Upper parent is", n.GetUpperParent().GetPhi())
		}
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
				n.GetContext(), n.LValue, n.LPrime+1)
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
			// I also implemented the compare method: SortedEquals for Clause

			// here is the code if we would want to receive the result
			// we computed concurrently

			// testDnf := <-ch

			// now test it and so on
			// if test fails return ErrNotSymmetric
			// end of this code snippet

			// TODO JGS You implemented this smmyetry test with comparing the DNFs,
			// does that mean Split and SplitBoth don't really need the symmetry
			// test variable?
		}
	}
	return nil
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

// ErrNotSymmetric may be returned by Split if the symmetric property
// is violated.
var ErrNotSymmetric error = errors.New("Found variables that are not symmetric.")

// SplittingTree represents the tree for a DNF.
type SplittingTree struct {
	Root                      *MainNode    // The root node
	Context                   *TreeContext // The context of the tree
	Renaming, ReverseRenaming []int        // See NewSplittingTree
	SymTest                   bool         // If true the test for symmetric variables is performed
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
//
// By default Cut and SymTest are set to true, so if you want
// to debug better set it by hand before calling CreateTree.
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
	return &SplittingTree{Root: root,
		Context:         context,
		Renaming:        renaming,
		ReverseRenaming: reverseRenaming,
		Cut:             true,
		SymTest:         true}
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

// CreateTree creates the whole splitting tree and returns ErrNotSymmetric
// if the symmetric property was violated.
//
// Think about a concurrent approach?
func (t *SplittingTree) CreateTree() error {
	// initialize the queue, we initialize it with some size
	// not a very good sice probably but it's something
	waiting := make([]SplitNode, 0, t.Root.GetContext().Nbvar)
	// add root node to queue
	waiting = append(waiting, t.Root)
	// loop while queue not empty
	for len(waiting) > 0 {
		// get next element
		next := waiting[0]
		waiting[0] = nil
		waiting = waiting[1:]
		if next.IsFinal() || next.IsAlreadySplit() {
			continue
		}
		if err := next.Split(t.SymTest, t.Cut); err != nil {
			return err
		}
		child1, child2 := next.GetUpperChild(), next.GetLowerChild()
		if child1 != nil && !child1.IsAlreadySplit() {
			waiting = append(waiting, child1)
		}
		if child2 != nil && !child2.IsAlreadySplit() {
			waiting = append(waiting, child2)
		}
	}
	return nil
}

// Interval defines an of the form interval (a, b].
// a and b are natural numbers, but we allow ∞ and -∞ as well.
//
// We will also used it sometimes if we just need a pair of numbers.
type Interval struct {
	LHS, RHS LPBCoeff
}

func NewInterval(lhs, rhs LPBCoeff) Interval {
	return Interval{lhs, rhs}
}

func (i Interval) String() string {
	return fmt.Sprintf("(%s, %s]", i.LHS, i.RHS)
}

// TreeSolver is an interface for everything that can transform a
// tree into an LPB or returns an error if that is not possible.
//
// The tree is not fully created when passed to Solve, i.e. it has to be created
// with CreateTree.
type TreeSolver interface {
	Solve(t *SplittingTree) (*LPB, error)
}

// SolverState provides the solver with certain information about the current
// search space, like current coefficients and so on.
type SolverState struct {
	Coefficients    []LPBCoeff   // The coefficients for each column, must be multiplied with the coeff factor.
	CoeffSums       []LPBCoeff   // Stores the sum of all coefficients including column k, gets updated in SetCoeff
	Intervals       [][]Interval // For each column contains all intervals in the tree
	IntervalFactors []int        // Factor intervals in a certain column must be multiplied with.
}

// NewSolverState a new solver space that sets all factors to one and
// intervals big enough to contain all intervals for the tree.
func NewSolverState(t *SplittingTree) *SolverState {
	size := t.Context.Nbvar + 1
	coefficients := make([]LPBCoeff, size)
	coeffSums := make([]LPBCoeff, size)
	intervals := make([][]Interval, size)
	intervalFactors := make([]int, size)
	// initialize all factors with one, for intervals create a slice for each
	// column that is big enough to hold all values
	for i := 0; i < size; i++ {
		coefficients[i] = NegativeInfinity // just something that is not a number
		coeffSums[i] = 0
		intervalFactors[i] = 1
		numRows := len(t.Context.Tree[i])
		intervals[i] = make([]Interval, numRows)
	}
	return &SolverState{Coefficients: coefficients,
		CoeffSums:       coeffSums,
		Intervals:       intervals,
		IntervalFactors: intervalFactors}
}

// SetCoeff sets the coefficient in the specified column, also updating the
// coefficient sum in that column.
func (s *SolverState) SetCoeff(column int, val LPBCoeff) {
	s.Coefficients[column] = val
	// in the last column there is no value to add
	if column == len(s.CoeffSums)-1 {
		s.CoeffSums[column] = val
	} else {
		s.CoeffSums[column] = s.CoeffSums[column+1].Add(val)
	}
}

// GetSumAfter returns the sum of all coefficients *after* the specified column,
// so not including this column.
func (s *SolverState) GetSumAfter(column int) LPBCoeff {
	if column == len(s.CoeffSums)-1 {
		return 0
	}
	return s.CoeffSums[column+1]
}

// GetInterval returns the current interval value taking into consideration
// the intervals for that column.
// Note that a new interval gets returned, so chaning the value you receive
// won't change anything here.
func (s *SolverState) GetInterval(column, row int) Interval {
	current := s.Intervals[column][row]
	factor := LPBCoeff(s.IntervalFactors[column])
	return NewInterval(current.LHS.Mult(factor), current.RHS.Mult(factor))
}

// SolveConflict solves a conflict of the form α < α_{k+1} < α + 1.
// It simply doubles the whole system (intervals and coefficients).
//
// If we are forced to solve a conflict in column k this means that we have
// to double all the intervals in column k and all following and all
// coefficients starting in column k + 1 (we haven't set a coefficient for k
// yet).
func (s *SolverState) SolveConflict(column int) {
	s.IntervalFactors[column] *= 2
	for k := column + 1; k < len(s.Coefficients); k++ {
		s.Coefficients[k] = s.Coefficients[k].Mult(2)
		s.CoeffSums[k] = s.Coefficients[k].Mult(2)
		s.IntervalFactors[k] *= 2
	}
}

// ColumnHandler is used to choose coefficients (coefficients α s.t. a < α < b)
// and the degree of the LPB // given an interval (a, b] and additional
// information such as the tree.
//
// It also must handle the columns with the HandleColumn function.
// This function takes the column that should be handled next. By handling we
// mean that it must compute all the intervals in that column and return the
// values a, b in which we must choose the coefficient (i.e. choose α s.t.
// a < α < b). If column = 0 we don't have to choose a coefficient, so in this
// case it may return whatever it wants.
//
// If you need more information than what is stored in the solver state this
// would be the right place, simply add all the information you need in your
// own type and make sure they're handled correctly.
// Of course you can also just implement your own solver at your will if none
// of the provided ones fits your purposes.
// When implementing your own chooser you should check the ComputeInterval
// function that computes the interval for a given node in the tree.
// MinColumnHandler is an example of such an implementation.
//
// All solving algorithms should make sure that these functions get only called
// with valid intervals, i.e. there must be a valid coefficient to choose.
// So implementations here must not check that a < b or something like that.
// Also it is not possible that we have a conflict of the form b = a + 1, this
// must be checked in another place.
//
// The Init function gets called each the handler should be initialized for a
// new tree.
type ColumnHandler interface {
	Init(t *SplittingTree)
	ChooseCoeff(i Interval, s *SolverState, t *SplittingTree, column int) (LPBCoeff, error)
	ChooseDegree(i Interval, s *SolverState, t *SplittingTree) (LPBCoeff, error)
	HandleColumn(s *SolverState, t *SplittingTree, column int) Interval
}

// ComputeInterval will compute the interval for the node in the specified
// column and row.
// For true and false the intervals are easy, for everything else the intervals
// and coefficient in the next column are looked at and s and b get computed
// accordingly.
//
// It will also set the interval in s.
func ComputeInterval(s *SolverState, t *SplittingTree, column, row int) Interval {
	var res Interval
	// get the node
	n := t.Root.GetContext().Tree[column][row]
	// check if dnf is true, false or something else
	sumSoFar := s.GetSumAfter(column)
	// TODO here were some asserts in the C++ code that don't make sense to me
	switch isFinal(n.GetPhi()) {
	case IsTrue:
		res = NewInterval(NegativeInfinity, 0)
	case IsFalse:
		res = NewInterval(sumSoFar, PositiveInfinity)
	default:
		uc := n.GetUpperChild()
		lc := n.GetLowerChild()
		switch {
		case uc == nil:
			res = NewInterval(sumSoFar, PositiveInfinity)
		case lc == nil:
			res = NewInterval(NegativeInfinity, 0)
		default:
			// neither is nil, so get the saved intervals
			// uc and lc column must be column+1
			upper := s.GetInterval(column+1, uc.GetRow())
			lower := s.GetInterval(column+1, lc.GetRow())
			s0, b0, s1, b1 := upper.LHS, upper.RHS, lower.LHS, lower.RHS
			lastCoeff := s.Coefficients[column+1]
			s := CoeffMax(s0, s1.Add(lastCoeff))
			b := CoeffMin(b0, b1.Add(lastCoeff))
			res = NewInterval(s, b)
		}
	}
	s.Intervals[column][row] = res
	return res
}

// degreeError is a small helper function that provides an error with the
// message that we can't choose a value that satisfies the conditions.
// I.e. we can't find a value α with a < α < b.
func degreeError(i Interval) error {
	return fmt.Errorf("Can't choose a value α s.t. %s < α < %s", i.LHS, i.RHS)
}

type MinColumnHandler struct{}

func NewMinColumnHandler() MinColumnHandler {
	return MinColumnHandler{}
}

func (handler MinColumnHandler) Init(t *SplittingTree) {
	// do nothing, no setup required
}

func (handler MinColumnHandler) ChooseCoeff(i Interval, s *SolverState, t *SplittingTree, column int) (LPBCoeff, error) {
	// we don't have to check if a < b or anything, this will already be done,
	// we just check for some weird cases... not entirely sure which of them
	// can ever happen
	switch i.LHS {
	case PositiveInfinity:
		return -1, degreeError(i)
	case NegativeInfinity:
		return 0, nil
	default:
		return i.LHS + 1, nil
	}
}

func (handler MinColumnHandler) ChooseDegree(i Interval, s *SolverState, t *SplittingTree) (LPBCoeff, error) {
	return handler.ChooseCoeff(i, s, t, -1)
}

func (handler MinColumnHandler) HandleColumn(s *SolverState, t *SplittingTree, column int) Interval {
	treeColumn := t.Root.GetContext().Tree[column]
	minSoFar := PositiveInfinity
	maxSoFar := NegativeInfinity
	// compute first interval for that column
	// we will need the last variable later to update the max and min
	// so we will later refer to last as "the interval before", that's why it
	// is called so
	last := ComputeInterval(s, t, column, 0)
	// iterate over all other rows
	numRows := len(treeColumn)
	for row := 1; row < numRows; row++ {
		current := ComputeInterval(s, t, column, row)
		n := treeColumn[row]
		if n.GetUpperParent() != nil { // TODO this is simpler than in C++, but should be ok? If it has an upper parent we must compare?
			diff1 := last.LHS.Sub(current.RHS)
			diff2 := last.RHS.Sub(current.LHS)
			maxSoFar = CoeffMax(maxSoFar, diff1)
			minSoFar = CoeffMin(minSoFar, diff2)
			// fmt.Printf("Min: %s - %s = %s\tMax: %s - %s = %s\n", last.RHS, current.LHS, diff2, last.LHS, current.RHS, diff1)
		}
		last = current
	}
	return NewInterval(maxSoFar, minSoFar)
}

type SimpleTreeSolver struct {
	handler ColumnHandler
	s       *SolverState
}

func NewSimpleTreeSolver(handler ColumnHandler) SimpleTreeSolver {
	return SimpleTreeSolver{handler: handler, s: nil}
}

func NewMinSolver() TreeSolver {
	return NewSimpleTreeSolver(NewMinColumnHandler())
}

func (solver SimpleTreeSolver) Solve(t *SplittingTree) (*LPB, error) {
	if err := t.CreateTree(); err != nil {
		return nil, err
	}
	solver.handler.Init(t)
	solver.s = NewSolverState(t)
	k := len(solver.s.Coefficients) - 1
	for k >= 0 {
		interval := solver.handler.HandleColumn(solver.s, t, k)
		if k == 0 {
			k--
			continue
		}
		// check if the interval makes sense, i.e. we don't have a conflict
		// and there is a possible solution
		// of course this must not be done in the first column
		max, min := interval.LHS, interval.RHS
		switch {
		case max.Add(1).Equals(min):
			// conflict, solve it!
			solver.s.SolveConflict(k)
			// multiply interval with 2
			interval.LHS *= 2
			interval.RHS *= 2
		case max.Compare(min) >= 0:
			// we can't choose a coefficient here!
			return nil, degreeError(interval)
		}
		// if we have reached this point we can choose a coefficient!
		coeff, chooseErr := solver.handler.ChooseCoeff(interval, solver.s, t, k)
		if chooseErr != nil {
			return nil, chooseErr
		}
		// now we can set the new coefficient
		solver.s.SetCoeff(k, coeff)
		k--
	}
	// once we reach this point we can choose the degree
	// first however we have to check if the interval is invalid
	rootInterval := solver.s.GetInterval(0, 0)
	if rootInterval.LHS.Compare(rootInterval.RHS) >= 0 {
		return nil, fmt.Errorf("Can't choose a degree in the interval %s", rootInterval)
	}
	degree, chooseErr := solver.handler.ChooseDegree(rootInterval, solver.s, t)
	if chooseErr != nil {
		return nil, chooseErr
	}
	// success, return the LPB!
	return NewLPB(degree, solver.s.Coefficients[1:]), nil
}
