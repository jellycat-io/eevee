package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/jellycat-io/eevee/config"
	"github.com/jellycat-io/eevee/lexer"
	"github.com/jellycat-io/eevee/logger"
	"github.com/jellycat-io/eevee/parser"
	"github.com/spf13/cobra"
)

var log = logger.New()

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Executes file at given path",
	Long:  `This command takes a filepath as argument`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config := config.GetConfig()

		filepath := args[0]
		if _, err := os.Stat(filepath); err != nil {
			log.Error(fmt.Sprintf(color.InRed("Invalid filepath. got=%q"), filepath))
			os.Exit(1)
		}

		buf, err := os.ReadFile(filepath)
		if err != nil {
			log.Error(fmt.Sprintf(color.InRed("Cannot read file: %q"), filepath))
		}
		source := strings.TrimSpace(string(buf))

		fmt.Println(source)

		l := lexer.New(source, config.TabSize)
		for _, t := range l.Tokens {
			fmt.Println(t)
		}

		p := parser.New(l.Tokens, false)
		ast := p.Parse()
		json, err := json.MarshalIndent(ast, "", "    ")
		if err != nil {
			log.Error(err.Error())
		}

		fmt.Printf("%s\n", json)

		if len(p.Errors()) != 0 {
			log.PrintParserErrors(p.Errors())
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
