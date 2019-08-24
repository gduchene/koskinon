/*
Copyright 2019 The koskinon Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lib

import (
	"reflect"
	"strings"
	"testing"
)

func TestLex(t *testing.T) {
	tests := []struct {
		input  string
		result []string
	}{
		// Valid inputs:
		{"hello, world", []string{"hello", ",", "world"}},
		{"`hello world`", []string{"hello world"}},
		{`"hello world"`, []string{"hello world"}},
		{"[`Hello`, world]", []string{"[", "Hello", ",", "world", "]"}},

		// Unsupported tokens in input:
		{"42 nope", nil},
		{"still nope 42.2", nil},
		{"forever nope 'c'", nil},
	}
	for i, test := range tests {
		ts, err := lex("", strings.NewReader(test.input))
		if err != nil {
			if test.result != nil {
				t.Errorf("#%d: lexing failed unexpectedly: %s", i, err)
			}
			continue
		}
		output := make([]string, len(ts))
		for j := range ts {
			output[j] = ts[j].val
		}
		if !reflect.DeepEqual(test.result, output) {
			t.Errorf("#%d: expected %v, got %v", i, test.result, output)
		}
	}
}
