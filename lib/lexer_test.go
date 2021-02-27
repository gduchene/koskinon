// SPDX-FileCopyrightText: © 2019 The koskinon Authors
// SPDX-License-Identifier: Apache-2.0

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
				t.Errorf("#%d: lex() failed: %s", i, err)
			}
			continue
		}
		output := make([]string, len(ts))
		for j := range ts {
			output[j] = ts[j].val
		}
		if !reflect.DeepEqual(test.result, output) {
			t.Errorf("#%d: expected %q, got %q", i, test.result, output)
		}
	}
}
