// SPDX-FileCopyrightText: © 2019 The koskinon Authors
// SPDX-License-Identifier: Apache-2.0

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

type PredBool bool

func (p PredBool) Eval(Message) bool {
	return bool(p)
}

type StmtLabel struct {
	Labels []string
}

type StmtMark struct{}

type StmtSkip struct{}

type StmtStop struct{}
