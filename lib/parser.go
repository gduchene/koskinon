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

func newParser(f string, r io.Reader) (*parser, error) {
	ts, err := lex(f, r)
	if err != nil {
		return nil, err
	}
	return &parser{0, ts}, nil
}
