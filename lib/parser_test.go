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
	r "regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_parseExprHeader(t *testing.T) {
	for i, test := range []struct {
		input string
		want  ExprHeader
	}{
		{
			`header "From" contains "foo"`,
			ExprHeader{[]string{"From"}, OpCmpContain{[]string{"foo"}}},
		},
		{
			`header "From" contains ["foo", "bar"]`,
			ExprHeader{[]string{"From"}, OpCmpContain{[]string{"foo", "bar"}}},
		},
		{
			`headers ["From", "To"] contain "foo"`,
			ExprHeader{[]string{"From", "To"}, OpCmpContain{[]string{"foo"}}},
		},
		{
			`headers ["From", "To"] contain ["foo", "bar"]`,
			ExprHeader{[]string{"From", "To"}, OpCmpContain{[]string{"foo", "bar"}}},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p, err := newParser("", strings.NewReader(test.input))
			require.NoError(t, err)
			expr, err := p.parseExprHeader()
			require.NoError(t, err)
			assert.Equal(t, test.want, expr)
		})
	}
}

func TestParser_parseExprMessage(t *testing.T) {
	for i, test := range []struct {
		input string
		want  ExprMessage
	}{
		{
			`message contains "foo"`,
			ExprMessage{OpCmpContain{[]string{"foo"}}},
		},
		{
			`message contains ["foo", "bar"]`,
			ExprMessage{OpCmpContain{[]string{"foo", "bar"}}},
		},
		{
			`message matches "fo+"`,
			ExprMessage{OpCmpMatch{[]*r.Regexp{r.MustCompile("fo+")}}},
		},
		{
			`message is "foo"`,
			ExprMessage{OpCmpEqual{[]string{"foo"}}},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p, err := newParser("", strings.NewReader(test.input))
			require.NoError(t, err)
			expr, err := p.parseExprMessage()
			require.NoError(t, err)
			assert.Equal(t, test.want, expr)
		})
	}
}

func TestParser_parseListStr(t *testing.T) {
	for i, test := range []struct {
		input   string
		want    ListStr
		wantErr bool
	}{
		// Valid inputs:
		{`["hello world"]`, ListStr{"hello world"}, false},
		{`["hello", "world"]`, ListStr{"hello", "world"}, false},
		{"[\"hello\", `world`]", ListStr{"hello", "world"}, false},

		// Invalid inputs:
		{`"hello world"`, ListStr{}, true},
		{`["hello",]`, ListStr{}, true},
		{`[,"hello",]`, ListStr{}, true},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p, err := newParser("", strings.NewReader(test.input))
			require.NoError(t, err)
			l, err := p.parseListStr()
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, l)
		})
	}
}

func TestParser_parseStmtLabel(t *testing.T) {
	for _, test := range []struct {
		input   string
		want    StmtLabel
		wantErr bool
	}{
		// Valid inputs:
		{`label "foo"`, StmtLabel{ListStr{"foo"}}, false},
		{`label ["foo"]`, StmtLabel{ListStr{"foo"}}, false},
		{`label ["foo", "bar"]`, StmtLabel{ListStr{"foo", "bar"}}, false},

		// Invalid inputs:
		{`label`, StmtLabel{}, true},
		{`label []`, StmtLabel{}, true},
		{`"label"`, StmtLabel{}, true},
	} {
		p, err := newParser("", strings.NewReader(test.input))
		require.NoError(t, err)
		stmt, err := p.parseStmtLabel()
		if test.wantErr {
			assert.Error(t, err)
			continue
		}
		require.NoError(t, err)
		assert.Equal(t, test.want, stmt)
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
	for i, test := range []struct {
		input   string
		want    StmtSkip
		wantErr bool
	}{
		// Valid input:
		{`skip inbox`, StmtSkip{}, false},

		// Invalid inputs:
		{`"skip inbox"`, StmtSkip{}, true},
		{`skip`, StmtSkip{}, true},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p, err := newParser("", strings.NewReader(test.input))
			require.NoError(t, err)
			stmt, err := p.parseStmtSkip()
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, stmt)
		})
	}
}

func TestParser_parseStmtStop(t *testing.T) {
	for i, test := range []struct {
		input   string
		want    StmtStop
		wantErr bool
	}{
		{"stop", StmtStop{}, false},
		{`"stop"`, StmtStop{}, true},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p, err := newParser("", strings.NewReader(test.input))
			require.NoError(t, err)
			stmt, err := p.parseStmtStop()
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, stmt)
		})
	}
}

func TestNewOpCmp(t *testing.T) {
	tests := []struct {
		op     string
		vals   []string
		result OpCmp
	}{
		// Valid inputs:
		{"is", []string{"foo", "bar"}, OpCmpEqual{[]string{"foo", "bar"}}},
		{"are", []string{"foo", "bar"}, OpCmpEqual{[]string{"foo", "bar"}}},
		{"contains", []string{"foo", "bar"}, OpCmpContain{[]string{"foo", "bar"}}},
		{"contain", []string{"foo", "bar"}, OpCmpContain{[]string{"foo", "bar"}}},
		{"matches", []string{"f."}, OpCmpMatch{[]*r.Regexp{r.MustCompile("f.")}}},
		{"match", []string{"f."}, OpCmpMatch{[]*r.Regexp{r.MustCompile("f.")}}},

		// Invalid inputs:
		{"matches", []string{"?"}, nil},
		{"matches", []string{"f(?:"}, nil},
		{"match", []string{"f**"}, nil},
	}
	for i, test := range tests {
		op, err := newOpCmp(test.op, test.vals)
		if err != nil {
			if test.result != nil {
				t.Errorf("#%d: newOpCmp() failed: %s", i, err)
			}
			continue
		}
		if !reflect.DeepEqual(test.result, op) {
			t.Errorf("#%d: expected %#v, got %#v", i, test.result, op)
		}
	}
}
