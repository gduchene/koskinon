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
	"golang.org/x/exp/ebnf"
	"os"
	"testing"
)

func TestGrammarIsValid(t *testing.T) {
	fd, err := os.Open("grammar.ebnf")
	if err != nil {
		t.Fatal(err)
	}
	g, err := ebnf.Parse("grammar.ebnf", fd)
	if err != nil {
		t.Fatal(err)
	}
	if err = ebnf.Verify(g, "Stmt"); err != nil {
		t.Error(err)
	}
}
