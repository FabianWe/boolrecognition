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

/*
Package boolrecognition is designed to recognize certain classes of boolean
functions, but also provides functions to work with boolean functions.

State of the Project

Currently we're working on solver for the so called threshold synthesis problem
(transforming a minimal DNF to a linear pseudo-Boolean constraint).
But we would welcome more contributions and other projects.

The code is contained in the subpackage lpb.

Components of this Package

This package defines some base types, such as clauses (a set of literals) and
clause sets (a set of clauses). We chose a simple form of representation: A
clause is just a slice of integers and a clause set (DNF or CNF) is a slice of
clauses. Thus each variable is identified by an integer id. How these ids are
used may depend on the problem domain. For example it is ok to identify each
variable with an integer 1 ≤ i ≤ n. Positive occurrences of i are denoted by
i and negative occurrences of i are denoted by -i. Sometimes it also can be
useful to start variable ids with 0 (for example for positive DNF only use
positive integers 0 ≤ i < n).

Also some algorithm may require that the clauses (or even the DNF / CNF) are
sorted in a particular way, make sure to document the code properly.

Winder matrix

We implemented to so called winder matrix (Threshold Logic by Robert O. Winder).
*/
package boolrecognition
