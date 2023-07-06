package main

import (
	"fmt"

	"github.com/jellycat-io/eevee/config"
	"github.com/jellycat-io/eevee/lexer"
)

func main() {
	config := config.GetConfig()

	source := `
if true then
    return 10
`

	l := lexer.NewLexer(source, config.TabSize)
	tokens := l.Tokens

	for _, token := range tokens {
		fmt.Printf("%v\n", token)
	}
}
