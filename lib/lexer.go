// SPDX-FileCopyrightText: © 2019 The koskinon Authors
// SPDX-License-Identifier: Apache-2.0

package lib

import (
	"fmt"
	"io"
	"text/scanner"
)

type kind int

const (
	kindUnknown kind = iota
	kindIdent
	kindOther
	kindString
)

func (k kind) String() string {
	return [...]string{"UNKNOWN", "IDENT", "OTHER", "STRING"}[k]
}

type token struct {
	k   kind
	pos scanner.Position
	val string
}

func (t token) String() string {
	return fmt.Sprintf("%s(%s)", t.k, t.val)
}

func lex(f string, r io.Reader) (ts []token, err error) {
	s := scanner.Scanner{}
	s.Init(r)
	s.Filename = f
	for t := s.Scan(); t != scanner.EOF; t = s.Scan() {
		k := kindUnknown
		v := s.TokenText()
		switch t {
		case scanner.Char, scanner.Int, scanner.Float:
			err = fmt.Errorf("%s: unexpected ``%s''", s.Position, s.TokenText())
			return
		case scanner.Ident:
			k = kindIdent
		case scanner.RawString, scanner.String:
			k = kindString
			v = v[1 : len(v)-1]
		default:
			k = kindOther
		}
		ts = append(ts, token{k, s.Position, v})
	}
	return
}
