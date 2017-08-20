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

// BinSearch performs a binary search on s to check if val is present.
// There is a similar function in the sort package (sort.Search), but this one
// uses a function as an argument.
// This makes the approach in sort slower, since binary search is a very
// simple procedure I've reimplemented it here.
//
// This method will return the index of val in the slice if it is present
// and -1 if val is not present.
//
// Some runtime comparisons (slice of size 100000):
// With sort.Search: 10130768 ns/op
// With this method: 7126174 ns/op
// which is a factor of ≈ 1.5
func BinSearch(s []int, val int) int {
	l, r := 0, len(s)-1
	for l <= r {
		m := l + (r-l)/2
		nxt := s[m]
		if nxt == val {
			return m
		}
		if nxt < val {
			l = m + 1
		} else {
			r = m - 1
		}
	}
	return -1
}
