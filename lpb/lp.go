// Copyright 2017 Fabian Wenzelmann <fabianwen@posteo.eu>, Christian Schilling,
// Jan-Georg Smaus
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
	"sort"
	"sync"

	br "github.com/FabianWe/boolrecognition"
)

// debug is used to panic in some conditions, if tested properly set to false.
const debug = true

// DNFTreeNodeContent is a node in the tree we construct for the regularity
// test. Each node stores a DNF, the information if that DNF is final (i.e.
// true or false), its depth and two children.
//
// The children are stored by ID, in the tree we store a list of all nodes
// and can therefore retrieve the actual node.
type DNFTreeNodeContent struct {
	phi                   br.ClauseSet
	leftChild, rightChild int
	final                 bool
	depth                 int
}

// A DNFTree is a collection of DNFTreeNodeContent objects.
// The root note is stored on position 0.
type DNFTree struct {
	Content []*DNFTreeNodeContent
	Nbvar   int
}

// NewDNFTree returns an empty tree containing no nodes.
func NewDNFTree(nbvar int) *DNFTree {
	return &DNFTree{Content: nil, Nbvar: nbvar}
}

// CreateNodeEntry creates a new node given its DNF, depth and the information
// if that DNF is final.
//
// It will append the new node to the tree and return the index of the new node.
func (tree *DNFTree) CreateNodeEntry(phi br.ClauseSet, depth int, isFinal bool) int {
	n := &DNFTreeNodeContent{phi, -1, -1, isFinal, depth}
	tree.Content = append(tree.Content, n)
	return len(tree.Content) - 1
}

// CreateRoot creates a root node and returns its ID (should be 0).
func (tree *DNFTree) CreateRoot(phi br.ClauseSet, isFinal bool) int {
	return tree.CreateNodeEntry(phi, 0, isFinal)
}

// CreateLeftChild creates a new node and sets the left child of nodeID
// to this node. Returns the ID of the new node.
func (tree *DNFTree) CreateLeftChild(nodeID int, phi br.ClauseSet, isFinal bool) int {
	if debug {
		if nodeID < 0 {
			panic("Expected nodeID >= 0 in CreateLeftChild")
		}
		if nodeID >= len(tree.Content) {
			panic("Expected nodeID < len(content) in CreateLeftChild")
		}
	}
	n := tree.Content[nodeID]
	id := tree.CreateNodeEntry(phi, n.depth+1, isFinal)
	n.leftChild = id
	return id
}

// CreateLeftChild creates a new node and sets the right child of nodeID
// to this node. Returns the ID of the new node.
func (tree *DNFTree) CreateRightChild(nodeID int, phi br.ClauseSet, isFinal bool) int {
	if debug {
		if nodeID < 0 {
			panic("Expected nodeID >= 0 in CreateRightChild")
		}
		if nodeID >= len(tree.Content) {
			panic("Expected nodeID < len(content) in CreateRightChild")
		}
	}
	n := tree.Content[nodeID]
	id := tree.CreateNodeEntry(phi, n.depth+1, isFinal)
	n.rightChild = id
	return id
}

// IsLeaf checks if the node is a leaf (has no child nodes).
func (tree *DNFTree) IsLeaf(nodeID int) bool {
	n := tree.Content[nodeID]
	return n.leftChild < 0 && n.rightChild < 0
}

type LPSplitResult struct {
	Final bool
	Phi   br.ClauseSet
}

func (tree *DNFTree) Split(nodeID int) (*LPSplitResult, *LPSplitResult) {
	n := tree.Content[nodeID]
	variable := n.depth
	first, second := br.NewClauseSet(len(n.phi)), br.NewClauseSet(len(n.phi))
	for _, clause := range n.phi {
		// check if the variable is contained
		if len(clause) > 0 && clause[0] == variable {
			// remove variable and add to the first result
			first = append(first, clause[1:])
		} else {
			// variable not contained, so just add the complete clause
			second = append(second, clause)
		}
	}
	isFirstFinal, isSecondFinal := isFinal(first) != NotFinal, isFinal(second) != NotFinal
	return &LPSplitResult{isFirstFinal, first}, &LPSplitResult{isSecondFinal, second}
}

// BuildTree will build the whole tree. The root note must be set already.
func (tree *DNFTree) BuildTree() {
	if debug {
		if len(tree.Content) != 1 {
			panic("Expected a tree containing exactly one node (the root) in BuildTree")
		}
	}
	if tree.Content[0].final {
		// for true and false there is nothing to do
		return
	}
	// create a queue that stores the node ids that must be explored
	// add first node (root) to it
	waiting := []int{0}
	for len(waiting) != 0 {
		nextID := waiting[0]
		waiting = waiting[1:]
		next := tree.Content[nextID]
		if next.final {
			// no splitting required for final node
			continue
		}
		// split the node
		first, second := tree.Split(nextID)
		if first.Final {
			if len(first.Phi) != 0 {
				leftID := tree.CreateLeftChild(nextID, first.Phi, true)
				waiting = append(waiting, leftID)
			}
			// TODO why only in this case?
		} else {
			leftID := tree.CreateLeftChild(nextID, first.Phi, false)
			waiting = append(waiting, leftID)
		}

		if second.Final {
			if len(second.Phi) != 0 {
				rightID := tree.CreateRightChild(nextID, second.Phi, true)
				waiting = append(waiting, rightID)
			}
		} else {
			rightID := tree.CreateRightChild(nextID, second.Phi, true)
			waiting = append(waiting, rightID)
		}
	}
}

func (tree *DNFTree) IsImplicant(mtp br.BooleanVector) bool {
	uID := 0
	for k := 0; k < len(mtp); k++ {
		u := tree.Content[uID]

		if tree.IsLeaf(uID) {
			return true
		}

		leftChild, rightChild := u.leftChild, u.rightChild
		if mtp[k] {
			if leftChild >= 0 {
				uID = leftChild
				continue
			} else {
				if debug {
					if rightChild < 0 {
						panic("rightChild must not be nil in IsImplicant")
					}
				}
				uID = rightChild
			}
		} else {
			if rightChild >= 0 {
				uID = rightChild
				continue
			} else {
				return false
			}
		}
	}
	if debug {
		if !(tree.Content[uID].leftChild < 0 && tree.Content[uID].rightChild < 0) {
			panic("rightChild and leftChild must be nil in IsImplicant")
		}
	}
	return true
}

func (tree *DNFTree) IsRegular(mtps []br.BooleanVector) bool {
	numRuns := tree.Nbvar - 1
	res := true
	// we will do this concurrently:
	// for each mtp iterate over all variable combinations and perform the test
	// and write the result to a channel
	// this also has some drawback: we need to wait for all mtps to finish
	// otherwise we would need some context wish would be too much here
	// so they all must write a result, even if one already returns false...
	report := make(chan bool, 10)
	// channel to report once we read all results
	done := make(chan bool)
	go func() {
		for i := 0; i < len(mtps); i++ {
			nxt := <-report
			if !nxt {
				res = false
			}
		}
		done <- true
	}()
	for k := 0; k < len(mtps); k++ {
		go func(index int) {
			mtp := mtps[index]
			check := true
			for i := 0; i < numRuns; i++ {
				if (!mtp[i]) && (mtp[i+1]) {
					// change the positions in the point, after the implicant test
					// we will change them again
					mtp[i] = true
					mtp[i+1] = false
					isImplicant := tree.IsImplicant(mtp)
					mtp[i] = false
					mtp[i+1] = true
					if !isImplicant {
						check = false
						break
					}
				}
			}
			report <- check
		}(k)
	}
	// wait until all results are there
	<-done
	return res
}

// TightenMode describes different modes to tighten the linear program
// before solving it.
//
// There are three different modes described below.
type TightenMode int

const (
	TightenNone       TightenMode = iota // Add only constraings necessary for solving the problem
	TightenNeighbours                    // Add also constraings between variables x(i) and x(i + 1)
	TightenAll                           // Add additional constraints between all variable pairs
)

type LinearProgram struct {
	Renaming, ReverseRenaming []int
	SymTest                   bool
	Tree                      *DNFTree
	Winder                    br.WinderMatrix
}

// NewLinearProgram creates a new lp given the DNF ϕ.
//
// It will however not create the actual program or the tree, this must be done
// somewhere else, it only creates the root node.
//
// Important note: For our algorithm to work the variables must be sorted
// according to their importance. Since this is not always the case (only
// during testing and some very special cases) this method will do this for
// you, i.e. it will create the winder matrix and then rename all
// variables accordingly. So the DNF we store in the root node is the
// renamed DNF. But we also store the mapping that caused this renaming
// in the field Renaming. This slice stores for each "old" variable
// the id in the new tree, i.e. a lookup tree.Renaming[id] gives you the
// id of the variable in the new tree.
// The reverse mapping, i.e. new variable → old variable is stored in
// ReverseRenaming.
//
// If you don't need the renaming set sortMatrix to false, in this case
// the matrix will work properly but the variables don't get sorted.
// That is only set it to false if you know that the ordering of the variables
// is already correct.
// Renaming and ReverseRenaming will be set to nil in this case.
//
// Also the clauses in the DNF must be sorted in increasing order.
// If you don't want the clauses to get sorted set sortClauses to false.
// Of course this only makes sense if also sortMatrix is set to false,
// otherwise the new dnf might not be sorted.
// This functions will sort them in this case nonetheless.
// TODO is this correct? I guess we need it later...
//
// The variables in the DNF have to be 0 <= v < nbar (so nbvar must be correct
// and variables start with 0).
// Also each variable should appear at least once in the DNF, what happens
// otherwise is not tested yet.
func NewLinearProgram(phi br.ClauseSet, nbvar int, sortMatrix, sortClauses bool) *LinearProgram {
	tree := NewDNFTree(nbvar)
	newDNF, winder, renaming, reverseRenaming := initLP(phi, nbvar, sortMatrix)
	if sortMatrix || sortClauses {
		newDNF.SortAll()
	}
	dnfType := isFinal(newDNF)
	rootID := tree.CreateRoot(newDNF, dnfType != NotFinal)
	if debug {
		if rootID != 0 {
			panic("Expected root id to be 0, in NewLinearProgram")
		}
	}
	return &LinearProgram{Renaming: renaming,
		ReverseRenaming: reverseRenaming,
		Tree:            tree,
		Winder:          winder}
}

// initLP initializes the lp, that is it creates the Winder matrix for
// (the renamed) ϕ.
// It will also compute Renaming and ReverseRenaming as discussed in
// NewLinearProgram.
//
// It returns first the renamedDNF, the Winder matrix, then Renaming and then
// ReverseRenaming.
// If sortMatrix is false the old dnf will be returned.
func initLP(phi br.ClauseSet, nbvar int, sortMatrix bool) (br.ClauseSet, br.WinderMatrix, []int, []int) {
	newDNF := phi
	var renaming, reverseRenaming []int = nil, nil
	winder := br.NewWinderMatrix(phi, nbvar, true)
	if sortMatrix {
		renaming = make([]int, nbvar)
		reverseRenaming = make([]int, nbvar)
		// sort the matrix
		winder.Sort()
		// create the renaming
		for newVariableId, row := range winder {
			renaming[row[len(row)-1]] = newVariableId
			reverseRenaming[newVariableId] = row[len(row)-1]
		}
		newDNF = make([]br.Clause, len(phi))
		// clone each clause
		// we'll do that concurrently
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
	return newDNF, winder, renaming, reverseRenaming
}

// ComputeMTPs computes the set of minimal true points of a minimal ϕ.
// Since ϕ is minimal this is easy: Each clause defines exactly one minimal
// true point.
func ComputeMTPs(phi br.ClauseSet, nbvar int) []br.BooleanVector {
	res := make([]br.BooleanVector, len(phi))
	for i, clause := range phi {
		point := br.NewBooleanVector(nbvar)
		res[i] = point
		for _, v := range clause {
			point[v] = true
		}
	}
	return res
}

// TODO test this with 0, I don't know what happens to the wait group
// otherwise, or just never call it with a DNF with zero clauses
func ComputeMFPs(mtps []br.BooleanVector, sortPoints bool) []br.BooleanVector {
	// first sort the mtps
	if sortPoints {
		cmp := func(i, j int) bool {
			p1, p2 := mtps[i], mtps[j]
			if debug {
				if len(p1) != len(p2) {
					panic("MTPS must be of same length in ComputeMFPs")
				}
			}
			size := len(p1)
			for k := 0; k < size; k++ {
				val1, val2 := p1[k], p2[k]
				if (!val1) && val2 {
					return true
				} else if val1 && (!val2) {
					return false
				}
			}
			if debug {
				panic("Must not reach this state in ComputeMFPs")
			}
			return false
		}
		sort.Slice(mtps, cmp)
	}
	// compute nu, we do this concurrently
	var wg sync.WaitGroup
	wg.Add(len(mtps) - 1)
	nu := make([]int, len(mtps))
	fmt.Println(len(mtps))
	for i := 1; i < len(mtps); i++ {
		go func(index int) {
			vars := len(mtps[index])
			for j := 0; j < vars; j++ {
				val1 := mtps[index-1][j]
				val2 := mtps[index][j]
				if (!val1) && val2 {
					nu[index] = j + 1
					break
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	// create the actual points, again we do that concurrently and communicate
	// via a channel
	// we range over that channel so we must not forget to close it!
	res := make([]br.BooleanVector, 0, 10)
	// start a function that listens on the channel and adds all points to the
	// result
	// we use a done channel to signal when all points have been added
	resChan := make(chan br.BooleanVector, 10)
	done := make(chan bool)
	go func() {
		for point := range resChan {
			res = append(res, point)
		}
		done <- true
	}()
	// in the wait group we wait until for all i we've added all points
	// after all points were written to the channel we close the channel and then
	// wait until they have been added to result
	wg.Add(len(mtps))
	for i := 0; i < len(mtps); i++ {
		go func(index int) {
			point := mtps[index]
			vars := len(point)
			for j := nu[index]; j < vars; j++ {
				if point[j] {

					if debug {
						if nu[index] > j {
							panic("nu[i] must be <= j in ComputeMFPs")
						}
					}

					newPoint := point.Clone()
					newPoint[j] = false
					for k := j + 1; k < vars; k++ {
						newPoint[k] = true
					}
					resChan <- newPoint
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	close(resChan)
	// now wait until all points were added to res
	<-done
	return res
}
