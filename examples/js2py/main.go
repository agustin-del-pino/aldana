package main

import (
	"fmt"
	"github.com/agustin-del-pino/aldana/pkg/aldana"
)

func main() {
	lexPy := NewPythonLexer()
	cur := aldana.NewCursor([]byte("let bar = 1234"))
	tks, tErr := lexPy.Tokenize(cur)

	if tErr != nil {
		fmt.Println("Lexer Error:")
		fmt.Println(aldana.GetLexerError(tErr, cur))
		return
	}

	prsPy := NewPythonParser()

	rdr := aldana.NewReader(tks)
	nd, nErr := prsPy.Parse(rdr)

	if nErr != nil {
		fmt.Println("Parse Error:")
		fmt.Println(nErr)
		tk := rdr.GetToken()
		fmt.Printf("line: %d column: %d token: %s type: %v\n", tk.Line, tk.Column, string(tk.Value), tk.Type)
		return
	}

	trnJs := NewJSTranspiler()
	b, bErr := trnJs.Transpile(nd)

	if bErr != nil {
		fmt.Println("Transpile Error:")
		fmt.Println(bErr)
		return
	}

	fmt.Println(string(b))

}
