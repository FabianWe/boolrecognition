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
	column := n.GetColumn()
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
				newDNF = append(newDNF, newClause)
				if len(newClause) == 0 {
					isResFinal = true
				}
			}
		}
	}
	if len(newDNF) == 0 {
		isResFinal = true
	}
	// create occurrence pattern
	occurrences := OPFromDNFShift(newDNF, n.GetContext().Nbvar, column+1)
	return NewSplitResult(isResFinal, newDNF, occurrences)
}
