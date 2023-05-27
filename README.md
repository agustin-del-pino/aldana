# aldana
Lexer/Parser/Transpiler tools

# Install

````shell
go get -u github.com/agustin-del-pino/aldana@latest
````

# The interfaces
There are three main interfaces: `Lexer`, `Parser`, `Transpiler`. Those can be used by implementing them or by using the default implementation of this lib.

# Default Implementation
Those implementations were thought as most generic possible. So, can be useful for not-big-deal use cases.

The approach of the "most generic possible" is bringing by Generics Types and Options-Configurations.

The `tokens` and `nodes` are generic types to the implementation, so theses delegate all behavior to callbacks and only collects what the callbacks returns or forwards the errors.

The Options-Configurations are the callbacks to be called at the implementation, usually they're named as "...Rule". *(`LexicalRule`, `NodeRule`, etc)*.

# How to use default implementations

For use the default implementations it's necessary create type for Tokens and Nodes. Both types can be anything, for this demonstration, both will be a structure.

````go
type Token {
    Type string
    Value []byte
}

type Node {
    Type string
    Value *Token
    Children []*Node
}
````

Now, declare a new token rule. In this case, the rule will create Numeric Tokens.
````go
func LexNum(c lexer.Cursor, r ranges.ByteRange) *Token {
    t := &Token { Type: "NUM" }
    for c.HasChar() && r(c.GetChar()) {
        t.Value = append(t.Value, c.GetChar())
        c.Next()
    }
    return t
}
````

After that, declare de node rule. In this case, will parse the Numeric Tokens.

````go
func ParseNum(r parser.Reader, f aldana.ParseRuleFinder[*Token, *Node]) (*Node, error) {
    nd := &Node{
        Type: "ROOT",
    }

    for r.HasTokens() {
        if r.GetToken().Type != "NUM" {
            return nil, parser.ErrInvalidSyntax
        }
        nd.Children = append(nd.Children, &Node{Type:"NUM_LIT", Token: r.GetToken()})
        r.Next()
    }

    return nd, nil
}
````

Finally, declare the byte transpiler. In this cases, will re-write the numbers but without any space.

````go
func TranspileNum(n *Node, t transpiler.Transpiler[*¨Node]) ([]byte, error) {
    var b []byte

    if n.Type == "ROOT" {
        for _, c := range n.Children {
            cb, cbErr := t.Transpile(c)

            if cbErr != nil {
                return nil, cbErr
            }

            b = append(b, cb...)
        }
    } else {
        b = n.Token.Value
    }
    return b, nil
}
````

So far, so good. Now, only left the main.

````go
func main() {
    lex := aldana.NewLexer(&aldana.LexerOptions[*Token]{
		Ignore: aldana.IgnoreWhiteSpaces(),
		LexRules: []aldana.LexicalRule[*Token]{
			aldana.NewLexicalRule(ranges.BytesBounded(0x30, 0x39), LexNum),
		},
	})

    prs := aldana.NewParser(&aldana.ParserOptions[*Token, *Node]{
        ParseRules: map[string]aldana.NodeRule[*Token, *Node]{
			"num": ParseNum,
		},
        Root: "num",
    })

    trn := aldana.NewTranspiler(&aldana.TranspilerOptions[*Node]{
		TranspileRules: []aldana.TranspileRule[*Node]{
			aldana.NewTranspileRule(func(_ *Node){return true}, TranspileNum),
        }
    })

    tks, tErr := lex.Tokenize(aldana.NewCursor([]byte("123 456 789")))

    if tErr != nil {
        panic(tErr)
    }

    nd, nErr := prs.Parse(aldana.NewReader(tks))

    if nErr != nil {
        panic(nErr)
    }

    b, bErr := trn.Transpile(nd)

    if bErr != nil {
        panic(bErr)
    }

    fmt.Println(string(b))
}
````
