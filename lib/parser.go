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
	"fmt"
	"io"
	"regexp"
)

type parser struct {
	i  int
	ts []token
}

// expect checks that the current token is of the expected kind and has
// the expected value when it is non-empty. If it does, the index is
// incremented and the value returned.
func (p *parser) expect(k kind, v string) (string, error) {
	if p.i >= len(p.ts) {
		return "", fmt.Errorf("%s: unexpected EOF", p.ts[len(p.ts)-1].pos)
	}
	t := &p.ts[p.i]
	if t.k != k || (v != "" && t.val != v) {
		return "", fmt.Errorf("%s: expected %s, got %s", t.pos, token{k: k, val: v}, t)
	}
	p.i++
	return t.val, nil
}

func (p *parser) expectIdent(v string) error {
	_, err := p.expect(kindIdent, v)
	return err
}

func (p *parser) expectOther(v string) error {
	_, err := p.expect(kindOther, v)
	return err
}

func (p *parser) nextIdent() (string, error) {
	return p.expect(kindIdent, "")
}

func (p *parser) nextOther() (string, error) {
	return p.expect(kindOther, "")
}

func (p *parser) nextStr() (string, error) {
	return p.expect(kindString, "")
}

func (p *parser) parseExprHeader() (expr ExprHeader, err error) {
	oi := p.i
	if err = p.expectIdent("header"); err != nil {
		if err = p.expectIdent("headers"); err != nil {
			return
		}
	}
	if expr.Headers, err = p.parseStrOrListStr(); err != nil {
		return
	}
	var op string
	if op, err = p.nextIdent(); err != nil {
		p.i = oi
		return
	}
	var vals []string
	if vals, err = p.parseStrOrListStr(); err != nil {
		p.i = oi
		return
	}
	expr.Op, err = newOpCmp(op, vals)
	return
}

func (p *parser) parseExprMessage() (expr ExprMessage, err error) {
	oi := p.i
	if err = p.expectIdent("message"); err != nil {
		return
	}
	var op string
	if op, err = p.nextIdent(); err != nil {
		p.i = oi
		return
	}
	var vals []string
	if vals, err = p.parseStrOrListStr(); err != nil {
		p.i = oi
		return
	}
	expr.Op, err = newOpCmp(op, vals)
	return
}

func (p *parser) parseListStr() (ListStr, error) {
	oi := p.i
	if err := p.expectOther("["); err != nil {
		return nil, err
	}
	l := []string{}
	for {
		s, err := p.nextStr()
		if err != nil {
			goto error
		}
		l = append(l, s)
		s, err = p.nextOther()
		if err != nil {
			goto error
		}
		switch s {
		case "]":
			return l, nil
		case ",":
			continue
		default:
			goto error
		}
	error:
		p.i = oi
		return nil, err
	}
}

func (p *parser) parsePredBool() (pred PredBool, err error) {
	if err = p.expectIdent("true"); err == nil {
		return PredBool(true), nil
	}
	if err = p.expectIdent("false"); err == nil {
		return PredBool(false), nil
	}
	return PredBool(false), err
}

func (p *parser) parseStmtLabel() (stmt StmtLabel, err error) {
	if err = p.expectIdent("label"); err != nil {
		return
	}
	stmt.Labels, err = p.parseStrOrListStr()
	return
}

func (p *parser) parseStmtMark() (stmt StmtMark, err error) {
	oi := p.i
	if err = p.expectIdent("mark"); err != nil {
		goto error
	}
	if err = p.expectIdent("as"); err != nil {
		goto error
	}
	if err = p.expectIdent("read"); err != nil {
		goto error
	}
	return
error:
	p.i = oi
	return
}

func (p *parser) parseStmtSkip() (stmt StmtSkip, err error) {
	oi := p.i
	if err = p.expectIdent("skip"); err != nil {
		goto error
	}
	if err = p.expectIdent("inbox"); err != nil {
		goto error
	}
	return
error:
	p.i = oi
	return
}

func (p *parser) parseStmtStop() (stmt StmtStop, err error) {
	oi := p.i
	if err = p.expectIdent("stop"); err != nil {
		p.i = oi
		return
	}
	return
}

func (p *parser) parseStrOrListStr() ([]string, error) {
	if s, err := p.nextStr(); err == nil {
		return []string{s}, nil
	}
	return p.parseListStr()
}

func newOpCmp(op string, vals []string) (OpCmp, error) {
	switch op {
	case "are", "is":
		return OpCmpEqual{vals}, nil
	case "contain", "contains":
		return OpCmpContain{vals}, nil
	case "match", "matches":
		rs := make([]*regexp.Regexp, len(vals))
		for i := range vals {
			r, err := regexp.Compile(vals[i])
			if err != nil {
				return nil, err
			}
			rs[i] = r
		}
		return OpCmpMatch{rs}, nil
	default:
		return nil, fmt.Errorf("unknown binary operator ``%s''", op)
	}
}

func newParser(f string, r io.Reader) (*parser, error) {
	ts, err := lex(f, r)
	if err != nil {
		return nil, err
	}
	return &parser{0, ts}, nil
}
