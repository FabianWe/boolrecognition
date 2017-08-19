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
	br "github.com/FabianWe/boolrecognition"
)

// VariableSetting accumlates some information we have about the variables.
// It stores the occurrence pattern for each variable (as computed from a DNF).
// Furthermore it stores the position of each variable in the global ordering,
// that is: The variables must not be stored according to the occurrence
// patterns, but after creating and sorting the patterns we assign each
// variable x the position in the pattern ordering.
// TODO remove this and create all the stuff somewhere else?
type VariableSetting struct {
	Patterns    []*OccurrencePattern
	VariablePos []int
}

// NewVariableSetting sets up a new variable setting for ϕ.
// It creates the occurrence patterns, sorts them and sets the variable
// ordering array accordingly.
func NewVariableSetting(phi br.ClauseSet, nbvar int) *VariableSetting {
	// first initalize the occurrence patterns
	patterns := OPFromDNF(phi, nbvar)
	// now sort them
	SortPatterns(patterns)
	// save the variable position
	pos := make([]int, nbvar)
	for i, op := range patterns {
		pos[op.Variable] = i
	}
	return &VariableSetting{Patterns: patterns,
		VariablePos: pos}
}

func (vs *VariableSetting) GetVariable(column int) int {
	return vs.Patterns[column].Variable
}

type TreeContext struct {
	Tree    [][]SplitNode
	Nbvar   int
	Setting *VariableSetting
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
	GetPhi() br.ClauseSet
	SetPhi(phi br.ClauseSet)
	GetColumn() int //
	SetColumn(column int)
	GetRow() int
	SetRow(row int)
	GetContext() *TreeContext //
	SetContext(context *TreeContext)
	GetPatterns() []*OccurrencePattern
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

// TODO implement symmetry test.
func Split(n SplitNode, k int, symmetryTest bool) *SplitResult {
	// ctx := n.GetContext()
	// column := n.GetColumn()
	// newOccurrences := make([]*OccurrencePattern, ctx.Nbvar-column-1)
	// // maybe too big...
	// newDNF := br.NewClauseSet(len(n.GetPhi()))
	// variable := ctx.Setting.GetVariable(column)
	// fmt.Println(newOccurrences)
	return nil
}
