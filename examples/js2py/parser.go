package main

import (
	"github.com/agustin-del-pino/aldana/pkg/aldana"
	"github.com/agustin-del-pino/aldana/pkg/aldana/parser"
)

type NodeType int

const (
	Root NodeType = iota
	VarDeclaration
	FuncDeclaration
	ArgsDeclaration
	ClassDeclaration
	ParamDeclaration
	NumberLiteral
	StringLiteral
	BoolLiteral
	BlockStatement
	PrivateClassVariableDeclaration
	MethodClassDeclaration
)

func IsNodeType(t NodeType) aldana.NodePredicator[*Node] {
	return func(n *Node) bool {
		return n.Type == t
	}
}

type Node struct {
	Type     NodeType
	Children []*Node
	Token    *Token
}

func parseRoot(r parser.Reader[*Token], f aldana.ParseRuleFinder[*Token, *Node]) (*Node, error) {
	nd := &Node{
		Type: Root,
	}
	for r.HasTokens() {
		if !IsTokenType(r.GetToken(), Word) {
			return nil, parser.ErrInvalidSyntax
		}
		dl, err := f("declaration", r)
		if err != nil {
			return nil, err
		}
		nd.Children = append(nd.Children, dl)
	}
	return nd, nil
}

func parseBlockStatement(r parser.Reader[*Token], f aldana.ParseRuleFinder[*Token, *Node]) (*Node, error) {
	nd := new(Node)

	if !IsTokenType(r.GetToken(), LeftBrace) {
		return nil, parser.ErrInvalidSyntax
	}

	nd.Type = BlockStatement
	r.Next()

	for r.HasTokens() && !IsTokenType(r.GetToken(), RightBrace) {
		dl, dErr := f("declaration", r)
		if dErr != nil {
			return nil, dErr
		}

		nd.Children = append(nd.Children, dl)
	}

	if !r.HasTokens() && !IsTokenType(r.GetToken(), RightBrace) {
		return nil, parser.ErrInvalidSyntax
	}

	return nd, nil
}

func parseArgsDeclaration(r parser.Reader[*Token], f aldana.ParseRuleFinder[*Token, *Node]) (*Node, error) {
	nd := new(Node)

	if !IsTokenType(r.GetToken(), LeftPrt) {
		return nil, parser.ErrInvalidSyntax
	}
	nd.Type = ArgsDeclaration
	r.Next()
	for r.HasTokens() && !IsTokenType(r.GetToken(), RightPrt) {
		if !IsTokenType(r.GetToken(), Word) {
			return nil, parser.ErrInvalidSyntax
		}
		arg := new(Node)
		arg.Type = ParamDeclaration
		arg.Token = r.GetToken()

		r.Next()

		t := r.GetToken().Type

		if t == RightPrt {
			break
		}

		switch t {
		case Comma:
			r.Next()
			continue
		case Eql:
			vl, vErr := f("value-assignation", r)

			if vErr != nil {
				return nil, vErr
			}
			arg.Children = append(arg.Children, vl)
		default:
			return nil, parser.ErrInvalidSyntax
		}

		nd.Children = append(nd.Children, arg)
	}

	if !r.HasTokens() && !IsTokenType(r.GetToken(), RightPrt) {
		return nil, parser.ErrInvalidSyntax
	}

	return nd, nil
}

func parseValueAssignation(r parser.Reader[*Token], f aldana.ParseRuleFinder[*Token, *Node]) (*Node, error) {
	nd := new(Node)

	switch {
	case IsTokenType(r.GetToken(), Num):
		nd.Type = NumberLiteral
	case IsTokenType(r.GetToken(), Str):
		nd.Type = StringLiteral
	case IsTokenType(r.GetToken(), Word) && (IsTokenValue(r.GetToken(), "false") || IsTokenValue(r.GetToken(), "true")):
		nd.Type = BoolLiteral
	default:
		return nil, parser.ErrUnhandledToken
	}
	nd.Token = r.GetToken()
	r.Next()
	return nd, nil
}

func parseFieldDeclaration(r parser.Reader[*Token], f aldana.ParseRuleFinder[*Token, *Node]) (*Node, error) {
	nd := new(Node)

	switch {
	case IsTokenType(r.GetToken(), Hash):
		r.Next()
		if !IsTokenType(r.GetToken(), Word) {
			return nil, parser.ErrInvalidSyntax
		}
		nd.Type = PrivateClassVariableDeclaration
		nd.Token = r.GetToken()
		r.Next()

		if !IsTokenType(r.GetToken(), Eql) {
			return nil, parser.ErrInvalidSyntax
		}

		vl, vErr := f("value-assignation", r)

		if vErr != nil {
			return nil, vErr
		}

		nd.Children = append(nd.Children, vl)

	case IsTokenType(r.GetToken(), Word):
		nd.Type = MethodClassDeclaration
		nd.Token = r.GetToken()
		args, aErr := f("args-declaration", r)

		if aErr != nil {
			return nil, aErr
		}

		nd.Children = append(nd.Children, args)

		body, bErr := f("block-statement", r)

		if bErr != nil {
			return nil, bErr
		}

		nd.Children = append(nd.Children, body)

	default:
		return nil, parser.ErrUnhandledToken
	}

	return nd, nil
}

func parseDeclaration(r parser.Reader[*Token], f aldana.ParseRuleFinder[*Token, *Node]) (*Node, error) {
	nd := new(Node)
	switch {
	case IsTokenValue(r.GetToken(), "class"):
		nd.Type = FuncDeclaration
		r.Next()
		if !IsTokenType(r.GetToken(), Word) {
			return nil, parser.ErrInvalidSyntax
		}
		nd.Token = r.GetToken()
		for r.HasTokens() && !IsTokenType(r.GetToken(), RightBrace) {
			fields, fErr := f("field-statement", r)

			if fErr != nil {
				return nil, fErr
			}

			nd.Children = append(nd.Children, fields)
		}

		if !r.HasTokens() && !IsTokenType(r.GetToken(), RightBrace) {
			return nil, parser.ErrInvalidSyntax
		}
		r.Next()

	case IsTokenValue(r.GetToken(), "function"):
		nd.Type = FuncDeclaration
		r.Next()
		if !IsTokenType(r.GetToken(), Word) {
			return nil, parser.ErrInvalidSyntax
		}
		nd.Token = r.GetToken()

		args, aErr := f("args-declaration", r)

		if aErr != nil {
			return nil, aErr
		}

		nd.Children = append(nd.Children, args)
		body, bErr := f("block-statement", r)

		if bErr != nil {
			return nil, bErr
		}

		nd.Children = append(nd.Children, body)
	case IsTokenValue(r.GetToken(), "var"), IsTokenValue(r.GetToken(), "const"), IsTokenValue(r.GetToken(), "let"):
		nd.Type = VarDeclaration
		r.Next()

		if !IsTokenType(r.GetToken(), Word) {
			return nil, parser.ErrInvalidSyntax
		}
		nd.Token = r.GetToken()
		r.Next()

		if !IsTokenType(r.GetToken(), Eql) {
			return nil, parser.ErrInvalidSyntax
		}
		r.Next()
		vl, err := f("value-assignation", r)

		if err != nil {
			return nil, err
		}

		nd.Children = append(nd.Children, vl)
	default:
		return nil, parser.ErrInvalidSyntax
	}

	return nd, nil
}

func NewPythonParser() parser.Parser[*Token, *Node] {
	return aldana.NewParser(&aldana.ParserOptions[*Token, *Node]{
		ParseRules: map[string]aldana.NodeRule[*Token, *Node]{
			"root":              parseRoot,
			"declaration":       parseDeclaration,
			"args-declaration":  parseArgsDeclaration,
			"field-declaration": parseFieldDeclaration,
			"value-assignation": parseValueAssignation,
			"block-statement":   parseBlockStatement,
		},
		Root: "root",
	})
}
