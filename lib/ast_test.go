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
	"regexp"
	"testing"
)

func TestOpCmpContain_Eval(t *testing.T) {
	tests := []struct {
		inputs, vals []string
		good         bool
	}{
		// Valid inputs:
		{[]string{"Hello world"}, []string{"world"}, true},
		{[]string{"Hello", "world"}, []string{"wor"}, true},

		// Invalid inputs:
		{[]string{"Hello world"}, []string{"Goodbye"}, false},
		{[]string{"Hello", "world"}, []string{"Hello world"}, false},
	}
	for i, test := range tests {
		ok := OpCmpContain{test.vals}.Eval(test.inputs)
		if !ok {
			if test.good {
				t.Errorf("#%d: Eval() failed: %q", i, test.inputs)
			}
			continue
		}
		if !test.good {
			t.Errorf("#%d: Eval() wrongly succeeded: %q", i, test.inputs)
		}
	}
}

func TestOpCmpEqual_Eval(t *testing.T) {
	tests := []struct {
		inputs, vals []string
		good         bool
	}{
		// Valid inputs:
		{[]string{"Hello world"}, []string{"Hello world"}, true},
		{[]string{"Hello", "world"}, []string{"world"}, true},

		// Invalid inputs:
		{[]string{"Hello world"}, []string{"Goodbye world"}, false},
		{[]string{"Hello", "world"}, []string{"universe"}, false},
	}
	for i, test := range tests {
		ok := OpCmpEqual{test.vals}.Eval(test.inputs)
		if !ok {
			if test.good {
				t.Errorf("#%d: Eval() failed: %q", i, test.inputs)
			}
			continue
		}
		if !test.good {
			t.Errorf("#%d: Eval() wrongly succeeded: %q", i, test.inputs)
		}
	}
}

func TestOpCmpMatch_Eval(t *testing.T) {
	tests := []struct {
		inputs, vals []string
		good         bool
	}{
		// Valid inputs:
		{[]string{"Hello world"}, []string{"w.*d"}, true},
		{[]string{"Hello", "world"}, []string{"wor"}, true},

		// Invalid inputs:
		{[]string{"Hello world"}, []string{"Go*dbye"}, false},
		{[]string{"Hello", "world"}, []string{"Hello  +world"}, false},
	}
	for i, test := range tests {
		rs := make([]*regexp.Regexp, len(test.vals))
		for j := range test.vals {
			rs[j] = regexp.MustCompile(test.vals[j])
		}
		ok := OpCmpMatch{rs}.Eval(test.inputs)
		if !ok {
			if test.good {
				t.Errorf("#%d: Eval() failed: %q", i, test.inputs)
			}
			continue
		}
		if !test.good {
			t.Errorf("#%d: Eval() wrongly succeeded: %q", i, test.inputs)
		}
	}
}
