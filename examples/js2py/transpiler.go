package main

import (
	"github.com/agustin-del-pino/aldana/pkg/aldana"
	"github.com/agustin-del-pino/aldana/pkg/aldana/transpiler"
)

func transpileVar(n *Node, t transpiler.Transpiler[*Node]) ([]byte, error) {
	var b []byte

	b = append(b, n.Token.Value...)
	b = append(b, ' ', '=', ' ')

	v, vErr := t.Transpile(n.Children[0])

	if vErr != nil {
		return nil, vErr
	}

	b = append(b, v...)

	return b, nil
}

func transpileNumLit(n *Node, t transpiler.Transpiler[*Node]) ([]byte, error) {
	var b []byte

	b = append(b, n.Token.Value...)

	return b, nil
}

func transpileStrLit(n *Node, t transpiler.Transpiler[*Node]) ([]byte, error) {
	var b []byte
	b = append(b, '"')
	b = append(b, n.Token.Value...)
	b = append(b, '"')
	return b, nil
}

func transpileBoolLit(n *Node, t transpiler.Transpiler[*Node]) ([]byte, error) {
	var b []byte

	switch {
	case IsTokenValue(n.Token, "true"):
		b = append(b, []byte("True")...)
	case IsTokenValue(n.Token, "false"):
		b = append(b, []byte("False")...)
	default:
		return nil, transpiler.ErrUnexpectedToken
	}

	return b, nil
}

func transpileRoot(n *Node, t transpiler.Transpiler[*Node]) ([]byte, error) {
	var b []byte
	for _, c := range n.Children {
		cb, cErr := t.Transpile(c)
		if cErr != nil {
			return nil, cErr
		}

		b = append(b, cb...)
		b = append(b, '\n')
	}
	return b, nil
}

func NewJSTranspiler() transpiler.Transpiler[*Node] {
	return aldana.NewTranspiler(&aldana.TranspilerOptions[*Node]{
		TranspileRules: []aldana.TranspileRule[*Node]{
			aldana.NewTranspileRule(IsNodeType(Root), transpileRoot),
			aldana.NewTranspileRule(IsNodeType(VarDeclaration), transpileVar),
			aldana.NewTranspileRule(IsNodeType(NumberLiteral), transpileNumLit),
			aldana.NewTranspileRule(IsNodeType(StringLiteral), transpileStrLit),
			aldana.NewTranspileRule(IsNodeType(BoolLiteral), transpileBoolLit),
		},
	})

}
