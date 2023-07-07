package main

import (
	"encoding/json"
	"fmt"

	"github.com/jellycat-io/eevee/config"
	"github.com/jellycat-io/eevee/lexer"
	"github.com/jellycat-io/eevee/parser"
	"github.com/jellycat-io/eevee/test"
)

func main() {
	config := config.GetConfig()

	source := test.MakeInput(
		"42",
		`"eevee"`,
	)

	l := lexer.NewLexer(source, config.TabSize)
	tokens := l.Tokens

	for _, token := range tokens {
		fmt.Printf("%v\n", token)
	}

	p := parser.NewParser(tokens)
	ast := p.Parse()

	jsonData, err := json.MarshalIndent(ast, "", "	")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonData))
	fmt.Println(ast.String())
}
