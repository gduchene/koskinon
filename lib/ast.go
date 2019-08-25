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
)

type ListStr []string

type Expr interface {
	Eval(Message) bool
}

type ExprHeader struct {
	Headers []string
	Op      OpCmp
}

func (e ExprHeader) Eval(m Message) bool {
	vals := []string{}
	headers := m.Headers()
	for _, header := range e.Headers {
		if val, ok := headers[header]; ok {
			vals = append(vals, val)
		}
	}
	return e.Op.Eval(vals)
}

type ExprMessage struct {
	Op OpCmp
}

func (e ExprMessage) Eval(m Message) bool {
	return e.Op.Eval([]string{m.Body()})
}

type Message interface {
	Body() string
	Headers() map[string]string
}

type OpCmp interface {
	Eval([]string) bool
}

type OpCmpContain struct {
	Values []string
}

func (o OpCmpContain) Eval(l []string) bool {
	for _, s1 := range l {
		for _, s2 := range o.Values {
			if strings.Contains(s1, s2) {
				return true
			}
		}
	}
	return false
}

type OpCmpEqual struct {
	Values []string
}

func (o OpCmpEqual) Eval(l []string) bool {
	for _, s1 := range l {
		for _, s2 := range o.Values {
			if s1 == s2 {
				return true
			}
		}
	}
	return false
}

type OpCmpMatch struct {
	Regexps []*regexp.Regexp
}

func (o OpCmpMatch) Eval(l []string) bool {
	for _, s1 := range l {
		for _, s2 := range o.Regexps {
			if s2.MatchString(s1) {
				return true
			}
		}
	}
	return false
}

type StmtLabel struct {
	Labels []string
}

type StmtMark struct{}

type StmtSkip struct{}

type StmtStop struct{}
