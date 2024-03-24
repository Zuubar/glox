package interpreter

import (
	"fmt"
	"glox/parser"
	"glox/scanner"
)

type loxTrait struct {
	stmt parser.TraitStmt
}

func newTrait(stmt parser.TraitStmt) *loxTrait {
	return &loxTrait{stmt: stmt}
}

func (t *loxTrait) Name() scanner.Token {
	return t.stmt.Name
}

func (t *loxTrait) Methods() []parser.FunctionStmt {
	return t.stmt.Methods
}

func (t *loxTrait) StaticMethods() []parser.FunctionStmt {
	return t.stmt.StaticMethods
}

func (t *loxTrait) String() string {
	return fmt.Sprintf("<trait %s>", t.Name().Lexeme)
}
