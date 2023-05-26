package aldana

import "github.com/agustin-del-pino/aldana/pkg/aldana/parser"

type ParseRuleFinder[Tt any, Tn any] func(n string, r parser.Reader[Tt]) (Tn, error)

// NodeRule is a function that returns the node made by the tokens and an error.
type NodeRule[Tt any, Tn any] func(r parser.Reader[Tt], f ParseRuleFinder[Tt, Tn]) (Tn, error)

// ParserOptions contains the parser's options.
type ParserOptions[Tt any, Tn any] struct {
	// ParseRules are the rules for parse the tokens identified by a name.
	ParseRules map[string]NodeRule[Tt, Tn]
	// Root is the ParserRule's as root Parser-Rule.
	Root string
}

// defaultParser implements parser.Parser.
type defaultParser[Tt any, Tn any] struct {
	// ops are the parser's options.
	ops *ParserOptions[Tt, Tn]
}

func (p *defaultParser[Tt, Tn]) Parse(r parser.Reader[Tt]) (Tn, error) {
	pr, ok := p.ops.ParseRules[p.ops.Root]

	if !ok {
		return *new(Tn), ErrNotFoundRootParserRule
	}

	r.Next()

	if !r.HasTokens() {
		return *new(Tn), ErrNoTokenToParser
	}

	nd, err := pr(r, p.findRule)

	if err != nil {
		return nd, err
	}

	if r.HasTokens() {
		return nd, parser.ErrUnhandledToken
	}

	return nd, nil
}

func (p *defaultParser[Tt, Tn]) findRule(n string, r parser.Reader[Tt]) (Tn, error) {
	pr, ok := p.ops.ParseRules[n]

	if !ok {
		return *new(Tn), ErrNotFundParserRule
	}

	return pr(r, p.findRule)
}

// NewParser returns the default implementation of parser.Parser.
//
// # Example
//
//	func parseExpression(r parser.Reader[*Token], p TokenPredicate[*Token], f ParseRuleFinder[*Token, *Node]) (*Node, error) {
//		nd, ndErr := f("term", r)
//		if ndErr != nil {
//			return nil, ndErr
//		}
//
//		for r.HasTokens() && p(r.GetToken()) {
//			term, err := f("term", r)
//			if err != nil {
//				return nil
//			}
//			nd = NewExpressionNode(nd, term)
//		}
//		return nd, nil
//	}
//
//	func parseTerm(r parser.Reader[*Token], p TokenPredicate[*Token], f ParseRuleFinder[*Token, *Node]) (*Node, error) {
//		nd, ndErr := f("factor", r)
//		if ndErr != nil {
//			return nil, ndErr
//		}
//
//		for r.HasTokens() && p(r.GetToken()) {
//			factor, err := f("factor", r)
//			if err != nil {
//				return nil
//			}
//			nd = NewTermNode(nd, factor)
//		}
//		return nd, nil
//	}
//
//	func parseFactor(r parser.Reader[*Token], p TokenPredicate[*Token], f ParseRuleFinder[*Token, *Node]) (*Node, error) {
//		if !r.HasToken() {
//			return nil, ErrInvalidSyntax
//		}
//		if !p(r.GetToken()) {
//			return nil, ErrUnhandledToken
//		}
//		return NewFactorNode(r.GetToken())
//	}
//
//	prs := NewParser[*Token, *Node](&ParserOptions{
//		ParseRules: map[string]ParseRule{
//			"expression": NewParseRule(IsExpressionToken, parseExpression),
//			"term": NewParseRule(IsTermToken, parseTerm),
//			"factor": NewParseRule(IsFactorToken, parseFactor)
//		}
//	})
func NewParser[Tt any, Tn any](ops *ParserOptions[Tt, Tn]) parser.Parser[Tt, Tn] {
	return &defaultParser[Tt, Tn]{
		ops: ops,
	}
}
