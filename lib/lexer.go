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

type token struct {
	kind int
	pos  scanner.Position
	val  string
}

const (
	kindUnknown = iota
	kindIdent
	kindOther
	kindString
)

func lex(f string, r io.Reader) (ts []token, err error) {
	s := scanner.Scanner{}
	s.Init(r)
	s.Filename = f
	for t := s.Scan(); t != scanner.EOF; t = s.Scan() {
		kind := kindUnknown
		val := s.TokenText()
		switch t {
		case scanner.Char, scanner.Int, scanner.Float:
			err = fmt.Errorf("%s: unexpected ``%s''", s.Position, s.TokenText())
			return
		case scanner.Ident:
			kind = kindIdent
		case scanner.RawString, scanner.String:
			kind = kindString
			val = val[1 : len(val)-1]
		default:
			kind = kindOther
		}
		ts = append(ts, token{kind, s.Position, val})
	}
	return
}