package repl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/user"

	"github.com/TwiN/go-color"
	"github.com/jellycat-io/eevee/lexer"
	"github.com/jellycat-io/eevee/logger"
	"github.com/jellycat-io/eevee/parser"
)

const PROMPT = "> "

var log = logger.New()

func Start(in io.Reader, out io.Writer) {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(in)

	fmt.Printf(color.InBlue("Eevee REPL 0.1.0 - Welcome %s\n"), user.Username)

	for {
		fmt.Fprint(out, color.InBold(PROMPT))
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line, 4)
		for _, t := range l.Tokens {
			fmt.Println(t)
		}
		p := parser.New(l.Tokens, true)
		ast := p.Parse()

		if len(p.Errors()) != 0 {
			log.PrintParserErrors(p.Errors())
		}

		json, err := json.MarshalIndent(ast, "", "    ")
		if err != nil {
			log.Error(err.Error())
		}

		fmt.Printf("%s\n", json)
	}
}
