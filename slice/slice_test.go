// Copyright 2013 com authors

//

// Licensed under the Apache License, Version 2.0 (the "License"): you may

// not use this file except in compliance with the License. You may obtain

// a copy of the License at

//

//     http://www.apache.org/licenses/LICENSE-2.0

//

// Unless required by applicable law or agreed to in writing, software

// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT

// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the

// License for the specific language governing permissions and limitations

// under the License.

package slice

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRemoveDuplicate(t *testing.T) {
	Convey("Remove duplicates from a slice of strings", t, func() {
		Convey("Remove duplicates from a slice with duplicates", func() {
			input := []string{"a", "b", "a", "c", "b"}
			expected := []string{"a", "b", "c"}
			result := RemoveDuplicate(input)
			So(result, ShouldResemble, expected)
		})

		Convey("Remove duplicates from a slice with no duplicates", func() {
			input := []string{"a", "b", "c"}
			expected := []string{"a", "b", "c"}
			result := RemoveDuplicate(input)
			So(result, ShouldResemble, expected)
		})

		Convey("Remove duplicates from an empty slice", func() {
			input := []string{}
			expected := []string{}
			result := RemoveDuplicate(input)
			So(result, ShouldResemble, expected)
		})
	})

	Convey("Remove duplicates from a slice of integers", t, func() {
		Convey("Remove duplicates from a slice with duplicates", func() {
			input := []int{1, 2, 1, 3, 2}
			expected := []int{1, 2, 3}
			result := RemoveDuplicate(input)
			So(result, ShouldResemble, expected)
		})

		Convey("Remove duplicates from a slice with no duplicates", func() {
			input := []int{1, 2, 3}
			expected := []int{1, 2, 3}
			result := RemoveDuplicate(input)
			So(result, ShouldResemble, expected)
		})

		Convey("Remove duplicates from an empty slice", func() {
			input := []int{}
			expected := []int{}
			result := RemoveDuplicate(input)
			So(result, ShouldResemble, expected)
		})
	})

	Convey("Remove duplicates from a slice of float64", t, func() {
		Convey("Remove duplicates from a slice with duplicates", func() {
			input := []float64{1.1, 2.2, 1.1, 3.3, 2.2}
			expected := []float64{1.1, 2.2, 3.3}
			result := RemoveDuplicate(input)
			So(result, ShouldResemble, expected)
		})

		Convey("Remove duplicates from a slice with no duplicates", func() {
			input := []float64{1.1, 2.2, 3.3}
			expected := []float64{1.1, 2.2, 3.3}
			result := RemoveDuplicate(input)
			So(result, ShouldResemble, expected)
		})

		Convey("Remove duplicates from an empty slice", func() {
			input := []float64{}
			expected := []float64{}
			result := RemoveDuplicate(input)
			So(result, ShouldResemble, expected)
		})
	})
}
