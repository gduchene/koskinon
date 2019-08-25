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
	"strings"
	"testing"
)

type mockMessage struct {
	body    string
	headers map[string]string
}

var testMessage Message = mockMessage{
	body: `Hello world!

I'm a test message. Something normal people would write. I think.


Cheers,
-- koskinon
`,
	headers: map[string]string{
		"From": "koskinon@example.com",
		"To":   "somebody@example.net",
	},
}

func (m mockMessage) Body() string {
	return m.body
}

func (m mockMessage) Headers() map[string]string {
	return m.headers
}

func TestExprHeader_Eval(t *testing.T) {
	tests := []struct {
		input string
		good  bool
	}{
		// Matching inputs:
		{`header "To" contains "example.net"`, true},
		{`header ["From", "To"] is "koskinon@example.com"`, true},
		{`header ["From", "To"] match "example.com$"`, true},
		{`header ["From", "X-Blah"] match "example.com$"`, true},

		// Non-matching inputs:
		{`header "From" is "koskinon@example.net"`, false},
		{`header "From" matches "^example.net"`, false},
		{`header ["To", "From"] contains "example.org"`, false},
	}
	for i, test := range tests {
		p, err := newParser("", strings.NewReader(test.input))
		if err != nil {
			t.Errorf("#%d: newParser() failed: %s", i, err)
			continue
		}
		expr, err := p.parseExprHeader()
		if err != nil {
			t.Errorf("#%d: parseExprHeader() failed: %s", i, err)
			continue
		}
		ok := expr.Eval(testMessage)
		if ok != test.good {
			t.Errorf("#%d: expected %t, got %t", i, test.good, ok)
		}
	}
}

func TestExprMessage_Eval(t *testing.T) {
	tests := []struct {
		input string
		good  bool
	}{
		// Matching inputs:
		{`message matches "(?m)^-- koskinon"`, true},
		{`message contains "Hello world"`, true},

		// Non-matching inputs:
		{`message is "koskinon"`, false},
		{`message matches ".*beep boop.*"`, false},
		{`message contains "example.org"`, false},
	}
	for i, test := range tests {
		p, err := newParser("", strings.NewReader(test.input))
		if err != nil {
			t.Errorf("#%d: newParser() failed: %s", i, err)
			continue
		}
		expr, err := p.parseExprMessage()
		if err != nil {
			t.Errorf("#%d: parseExprMessage() failed: %s", i, err)
			continue
		}
		ok := expr.Eval(testMessage)
		if ok != test.good {
			t.Errorf("#%d: expected %t, got %t", i, test.good, ok)
		}
	}
}

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
