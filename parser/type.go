package parser

import "hek/ast"

type prefixParseFun func() ast.Expression
type infixParseFun func(expression ast.Expression) ast.Expression
