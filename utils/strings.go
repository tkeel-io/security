/*
Copyright 2021 The tKeel Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

// StringsIndexOf returns index position in slice from given string
// If value is -1, the string does not found.
func StringsIndexOf(slice []string, s string) int {
	for i, v := range slice {
		if v == s {
			return i
		}
	}

	return -1
}

// StringsInclude returns true or false if given string is in slice.
func StringsInclude(slice []string, s string) bool {
	return StringsIndexOf(slice, s) >= 0
}

// StringsUniqueAppend appends a string if not exist in the slice.
func StringsUniqueAppend(slice []string, str ...string) []string {
	for i := range str {
		if StringsIndexOf(slice, str[i]) != -1 {
			continue
		}

		slice = append(slice, str[i])
	}

	return slice
}
