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

func TestParser_parseListStr(t *testing.T) {
	tests := []struct {
		input  string
		result ListStr
	}{
		// Valid inputs:
		{`["hello world"]`, ListStr{"hello world"}},
		{`["hello", "world"]`, ListStr{"hello", "world"}},
		{"[\"hello\", `world`]", ListStr{"hello", "world"}},

		// Invalid inputs:
		{`"hello world"`, ListStr{}},
		{`["hello",]`, ListStr{}},
		{`[,"hello",]`, ListStr{}},
	}
	for i, test := range tests {
		p, err := newParser("", strings.NewReader(test.input))
		if err != nil {
			t.Errorf("#%d: newParser() failed: %s", i, err)
			continue
		}
		l, err := p.parseListStr()
		if err != nil {
			if len(test.result) != 0 {
				t.Errorf("#%d: parseListStr() failed: %s", i, err)
			}
			continue
		}
		if !reflect.DeepEqual(test.result, l) {
			t.Errorf("#%d: expected %q, got %q", i, test.result, l)
		}
	}
}

func TestParser_parseStmtLabel(t *testing.T) {
	tests := []struct {
		input  string
		result StmtLabel
	}{
		// Valid inputs:
		{`label "foo"`, StmtLabel{ListStr{"foo"}}},
		{`label ["foo"]`, StmtLabel{ListStr{"foo"}}},
		{`label ["foo", "bar"]`, StmtLabel{ListStr{"foo", "bar"}}},

		// Invalid inputs:
		{`label`, StmtLabel{}},
		{`label []`, StmtLabel{}},
		{`"label"`, StmtLabel{}},
	}
	for i, test := range tests {
		p, err := newParser("", strings.NewReader(test.input))
		if err != nil {
			t.Errorf("#%d: newParser() failed: %s", i, err)
			continue
		}
		stmt, err := p.parseStmtLabel()
		if err != nil {
			if len(test.result.Labels) != 0 {
				t.Errorf("#%d: parseStmtLabel() failed: %s", i, err)
			}
			continue
		}
		if !reflect.DeepEqual(test.result, stmt) {
			t.Errorf("#%d: expected %q, got %q", i, test.result, stmt)
		}
	}
}

func TestParser_parseStmtMark(t *testing.T) {
	tests := []struct {
		input string
		good  bool
	}{
		// Valid input:
		{`mark as read`, true},

		// Invalid inputs:
		{`"mark as read"`, false},
		{`mark`, false},
	}
	for i, test := range tests {
		p, err := newParser("", strings.NewReader(test.input))
		if err != nil {
			t.Errorf("#%d: newParser() failed: %s", i, err)
			continue
		}
		if _, err := p.parseStmtMark(); err != nil {
			if test.good {
				t.Errorf("#%d: parseStmtMark() failed: %s", i, err)
			}
			continue
		}
		if !test.good {
			t.Errorf("#%d: parseStmtMark() wrongly succeeded: %s", i, test.input)
		}
	}
}

func TestParser_parseStmtSkip(t *testing.T) {
	tests := []struct {
		input string
		good  bool
	}{
		// Valid input:
		{`skip inbox`, true},

		// Invalid inputs:
		{`"skip inbox"`, false},
		{`skip`, false},
	}
	for i, test := range tests {
		p, err := newParser("", strings.NewReader(test.input))
		if err != nil {
			t.Errorf("#%d: newParser() failed: %s", i, err)
			continue
		}
		if _, err := p.parseStmtSkip(); err != nil {
			if test.good {
				t.Errorf("#%d: parseStmtSkip() failed: %s", i, err)
			}
			continue
		}
		if !test.good {
			t.Errorf("#%d: parseStmtSkip() wrongly succeeded: %s", i, test.input)
		}
	}
}

func TestParser_parseStmtStop(t *testing.T) {
	tests := []struct {
		input string
		good  bool
	}{
		// Valid input:
		{`stop`, true},

		// Invalid input:
		{`"stop"`, false},
	}
	for i, test := range tests {
		p, err := newParser("", strings.NewReader(test.input))
		if err != nil {
			t.Errorf("#%d: newParser() failed: %s", i, err)
			continue
		}
		if _, err := p.parseStmtStop(); err != nil {
			if test.good {
				t.Errorf("#%d: parseStmtStop() failed: %s", i, err)
			}
			continue
		}
		if !test.good {
			t.Errorf("#%d: parseStmtStop() wrongly succeeded: %s", i, test.input)
		}
	}
}
