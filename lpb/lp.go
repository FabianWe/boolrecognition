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

import br "github.com/FabianWe/boolrecognition"

const debug = true

type DNFTreeNodeContent struct {
	phi                   br.ClauseSet
	leftChild, rightChild int
	final                 bool
	depth                 int
}

type DNFTree struct {
	content []*DNFTreeNodeContent
}

func NewDNFTree() *DNFTree {
	return &DNFTree{content: nil}
}

func (tree *DNFTree) CreateNodeEntry(phi br.ClauseSet, depth int, isFinal bool) int {
	n := &DNFTreeNodeContent{phi, -1, -1, isFinal, depth}
	tree.content = append(tree.content, n)
	return len(tree.content) - 1
}

func (tree *DNFTree) CreateRoot(phi br.ClauseSet, depth int, isFinal bool) int {
	return tree.CreateNodeEntry(phi, 0, isFinal)
}

func (tree *DNFTree) CreateLeftChild(nodeID int, phi br.ClauseSet, isFinal bool) int {
	if debug {
		if nodeID < 0 {
			panic("Expected nodeID >= 0 in CreateLeftChild")
		}
		if nodeID >= len(tree.content) {
			panic("Expected nodeID < len(content) in CreateLeftChild")
		}
	}
	n := tree.content[nodeID]
	id := tree.CreateNodeEntry(phi, n.depth+1, isFinal)
	n.leftChild = id
	return id
}

func (tree *DNFTree) CreateRightChild(nodeID int, phi br.ClauseSet, isFinal bool) int {
	if debug {
		if nodeID < 0 {
			panic("Expected nodeID >= 0 in CreateRightChild")
		}
		if nodeID >= len(tree.content) {
			panic("Expected nodeID < len(content) in CreateRightChild")
		}
	}
	n := tree.content[nodeID]
	id := tree.CreateNodeEntry(phi, n.depth+1, isFinal)
	n.rightChild = id
	return id
}

func (tree *DNFTree) IsLeaf(nodeID int) bool {
	n := tree.content[nodeID]
	return n.leftChild < 0 && n.rightChild < 0
}
