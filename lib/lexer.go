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
