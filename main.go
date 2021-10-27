package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type (
	Evaluator map[string]float64
)

func NewEvaluator() Evaluator {
	return make(Evaluator)
}

func (e Evaluator) Eval(input string) string {
	l := NewLexer(input, 10).Run()

	root, err := NewParser(l.out, e).Parse()
	if err != nil {
		return err.Error()
	}

	val, err := root.Eval()
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("=%f", val)
}

func (e Evaluator) String() string {
	b := strings.Builder{}
	b.WriteString("Values:")
	for k, v := range e {
		b.WriteString(fmt.Sprintf(" %s=%f", k, v))
	}
	return b.String()
}

func (e Evaluator) PrettyString() string {
	b := strings.Builder{}
	b.WriteString("Values:")
	for k, v := range e {
		b.WriteString(fmt.Sprintf(" %s=%f\n", k, v))
	}
	return b.String()
}

func main() {
	// check for one shot mode
	input := flag.String("e", "", "put a raw expression here, if an argument is supplied an interactive session will not be started")
	flag.Parse()

	e := NewEvaluator()

	if *input != "" {
		fmt.Println(e.Eval(*input))
		return
	}

	// do interactive mode
	fmt.Println("Type 'exit' followed by pressing return/enter to exit the program.")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		switch {
		case input == "exit":
			fmt.Println("exiting...")
			return
		case input == "clear":
			e = NewEvaluator()
		case input == "list":
			fmt.Println(e)
		case input == "list-pretty":
			fmt.Println(e.PrettyString())
		// case: TODO input starts with precision, followed by number, update precision for ouputs
		default:
			fmt.Println(e.Eval(input))
		}
	}
}
