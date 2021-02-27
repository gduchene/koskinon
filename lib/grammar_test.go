// SPDX-FileCopyrightText: © 2019 The koskinon Authors
// SPDX-License-Identifier: Apache-2.0

package lib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/ebnf"
)

func TestGrammarIsValid(t *testing.T) {
	fd, err := os.Open("grammar.ebnf")
	require.NoError(t, err)
	g, err := ebnf.Parse("grammar.ebnf", fd)
	require.NoError(t, err)
	assert.NoError(t, ebnf.Verify(g, "Stmt"))
}
